package dataconverter

import (
	"errors"
	"fmt"

	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/converter"
)
const (
	metadataEncryptionKey     = "encryption"
	metadataCompressionKey    = "compression"
	metadataEncryptedAESV1    = "AESV1"
	metadataCompressionGZV1   = "GZV1"
	minimumPreCompressionSize = 1 << 12 // ~4KB
)

type Options struct {
	// EncryptionKey is the encryption key used to encrypt the payloads
	// this key must be 16, 24, 32 characters in length
	EncryptionKey []byte

	// CompressionEnabled if you wish to compress payloads
	// in addition to encrypting them this might be beneficial if
	// working with larger payloads.
	CompressionEnabled bool
}

type encryptDataConverterV1 struct {
	encryptionService  *AESEncryptionServiceV1
	payloadConverters  map[string]converter.PayloadConverter
	orderedEncodings   []string
	compressionEnabled bool
}

type nightfallTemporalEncodings struct {
	encoding string
	isAESV1 bool
	isGZV1  bool
}

var (

	// ErrMetadataIsNotSet is returned when metadata is not set.
	ErrMetadataIsNotSet = errors.New("metadata is not set")
	// ErrEncodingIsNotSet is returned when payload encoding metadata is not set.
	ErrEncodingIsNotSet = errors.New("payload encoding metadata is not set")
	// ErrEncodingIsNotSupported is returned when payload encoding is not supported.
	ErrEncodingIsNotSupported = errors.New("payload encoding is not supported")
	//ErrUnableToFindConverter is return when payload converter is not found
	ErrUnableToFindConverter = errors.New("payload converter is not found")
)

// NewEncryptDataConverterV1 - Temporal provides a default unencrypted DataConverter however
// for some of our needs we need a DataConverter to encrypt maybe sensitive information
// into workflows. EncryptDataConverterV1 allows the ability to encrypt maybe sensitive
// workflows without compromising sensitive info we send to our temporal service.
func NewEncryptDataConverterV1(opts Options) (converter.DataConverter, error) {
	defaultTemporalPayloadConverters := []converter.PayloadConverter{
		converter.NewNilPayloadConverter(),
		converter.NewByteSlicePayloadConverter(),
		// Only one proto converter should be used.
		// Although they check for different interfaces (proto.Message and proto.Marshaler) all proto messages implements both interfaces.
		converter.NewProtoJSONPayloadConverter(),
		// NewProtoPayloadConverter(),
		converter.NewJSONPayloadConverter(),
		// nightfallCustomConverter
	}
	encryptionService, err := newAESEncryptionServiceV1(opts)
	if err != nil {
		return nil, err
	}
	dc := &encryptDataConverterV1{
		payloadConverters:  make(map[string]converter.PayloadConverter, len(defaultTemporalPayloadConverters)),
		orderedEncodings:   make([]string, len(defaultTemporalPayloadConverters)),
		encryptionService:  encryptionService,
		compressionEnabled: opts.CompressionEnabled,
	}
	for i, payloadConverter := range defaultTemporalPayloadConverters {
		dc.payloadConverters[payloadConverter.Encoding()] = payloadConverter
		dc.orderedEncodings[i] = payloadConverter.Encoding()
	}
	return dc, nil
}

func (dc *encryptDataConverterV1) ToPayloads(values ...interface{}) (*commonpb.Payloads, error) {
	if len(values) == 0 {
		return nil, nil
	}

	result := &commonpb.Payloads{
		Payloads: make([]*commonpb.Payload, len(values)),
	}
	for i := range values {
		payload, err := dc.ToPayload(values[i])
		if err != nil {
			return nil, fmt.Errorf("values[%d]: %w", i, err)
		}
		result.Payloads[i] = payload
	}

	return result, nil
}

// FromPayloads converts to a list of values of different types.
func (dc *encryptDataConverterV1) FromPayloads(payloads *commonpb.Payloads, valuePtrs ...interface{}) error {
	if payloads == nil {
		return nil
	}

	for i, payload := range payloads.GetPayloads() {
		if i >= len(valuePtrs) {
			break
		}

		err := dc.FromPayload(payload, valuePtrs[i])
		if err != nil {
			return fmt.Errorf("payload item %d: %w", i, err)
		}
	}

	return nil
}

// ToPayload converts single value to payload.
func (dc *encryptDataConverterV1) ToPayload(value interface{}) (*commonpb.Payload, error) {
	for _, enc := range dc.orderedEncodings {
		unencryptedPayload, err := dc.payloadConverters[enc].ToPayload(value)
		if err != nil {
			return nil, err
		}
		if unencryptedPayload != nil {
			return dc.compressAndEncryptPayload(unencryptedPayload)
		}
	}

	return nil, fmt.Errorf("value: %v of type: %T: %w", value, value, ErrUnableToFindConverter)
}

// FromPayload converts single value from payload.
func (dc *encryptDataConverterV1) FromPayload(payload *commonpb.Payload, valuePtr interface{}) error {
	if payload == nil {
		return nil
	}
	nightfallEncodings, err := dc.decryptAndDecompress(payload)
	if err != nil {
		return err
	}
	payloadConverter, ok := dc.payloadConverters[nightfallEncodings.encoding]
	if !ok {
		return fmt.Errorf("encoding %s: %w", nightfallEncodings.encoding, ErrEncodingIsNotSupported)
	}

	return payloadConverter.FromPayload(payload, valuePtr)
}

// ToString converts payload object into human readable string.
func (dc *encryptDataConverterV1) ToString(payload *commonpb.Payload) string {
	if payload == nil {
		return ""
	}
	nightfallEncodings, err := dc.decryptAndDecompress(payload)
	if err != nil {
		return err.Error()
	}
	payloadConverter, ok := dc.payloadConverters[nightfallEncodings.encoding]
	if !ok {
		return fmt.Errorf("encoding %s: %w", nightfallEncodings.encoding, ErrEncodingIsNotSupported).Error()
	}

	return payloadConverter.ToString(payload)
}

// ToStrings converts payloads object into human-readable strings.
func (dc *encryptDataConverterV1) ToStrings(payloads *commonpb.Payloads) []string {
	if payloads == nil {
		return nil
	}

	result := make([]string, len(payloads.GetPayloads()) )
	for idx := range payloads.GetPayloads(){
		result[idx] = dc.ToString(payloads.GetPayloads()[idx])
	}

	return result
}

// the unencrypted payload isn't mutable and need to reconstruct the payload.
func (dc *encryptDataConverterV1) compressAndEncryptPayload(unencryptedPayload *commonpb.Payload) (*commonpb.Payload, error){
	newMetadata := unencryptedPayload.GetMetadata()
	if newMetadata == nil {
		newMetadata = make(map[string][]byte)
	}
	newMetadata[metadataEncryptionKey] = []byte(metadataEncryptedAESV1)

	// compression can be resource intensive on cpu/mem
	// we only should be compressing things that are bigger in payload size
	// as these payloads we should benefit from a size reduction
	if dc.compressionEnabled && len(unencryptedPayload.GetData()) > minimumPreCompressionSize {
		compressedData, err := compressGZV1(unencryptedPayload.GetData())
		if err != nil {
			return nil, err
		}
		newMetadata[metadataCompressionKey] = []byte(metadataCompressionGZV1)
		unencryptedPayload.Data = compressedData
	}
	encryptedBytes, err := dc.encryptionService.Encrypt(unencryptedPayload.GetData())
	if err != nil {
		return &commonpb.Payload{}, err
	}

	return &commonpb.Payload{
		Metadata: newMetadata,
		Data: encryptedBytes,
	}, nil
}

// decryptAndDecompress figures out from metadata whether the payload needs to be decrypted and/or decompressed
func (dc *encryptDataConverterV1) decryptAndDecompress(payload *commonpb.Payload) (nightfallTemporalEncodings, error) {
	nightfallEncodings, err := encoding(payload)
	if err != nil {
		return nightfallTemporalEncodings{}, err
	}
	if nightfallEncodings.isAESV1 {
		if payload.Data, err = dc.encryptionService.Decrypt(payload.GetData()); err != nil {
			return nightfallTemporalEncodings{}, err
		}
	}
	if nightfallEncodings.isGZV1 {
		if payload.Data, err = decompressGZV1(payload.GetData()); err != nil {
			return nightfallTemporalEncodings{}, err
		}
	}
	return nightfallEncodings, nil
}

func encoding(payload *commonpb.Payload) (nightfallTemporalEncodings, error) {
	metadata := payload.GetMetadata()
	if metadata == nil {
		return nightfallTemporalEncodings{}, ErrMetadataIsNotSet
	}
	encryptionType, hasEncryption := metadata[metadataEncryptionKey]
	compressionType, hasCompression := metadata[metadataCompressionKey]
	if encoding, ok := metadata[converter.MetadataEncoding]; ok {
		return nightfallTemporalEncodings{
			encoding: string(encoding),
			isAESV1: hasEncryption && (string(encryptionType) == metadataEncryptedAESV1),
			isGZV1: hasCompression && (string(compressionType) == metadataCompressionGZV1),
		}, nil
	}
	return nightfallTemporalEncodings{}, ErrEncodingIsNotSet
}
