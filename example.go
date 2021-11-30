package temporal_dataconverter_encryptor_go

import (
	"github.com/cowell21/temporal-dataconverter-encryptor-go/dataconverter"
	"go.temporal.io/sdk/client"
)

// NewTemporalClient new temporal client example with encryption payload converter
func NewTemporalClient() (client.Client, error) {
	encryptedDataConverter, err := dataconverter.NewEncryptDataConverterV1(dataconverter.Options{
		EncryptionKey: []byte("00000000~secretGoesHere~00000000"),
	})
	if err != nil {
		return nil, err
	}
	opts := client.Options{
		Namespace:    "default",
		DataConverter: encryptedDataConverter,
	}
	return client.NewClient(opts)
}
