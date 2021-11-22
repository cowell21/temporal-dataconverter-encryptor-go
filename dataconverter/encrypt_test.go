package dataconverter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const(
	developmentEncryptionStringKey = "00000000000000000000000000000000" // 32 length
	emptyString = ""
	testString = "hello"
	weirdCharsString = `
		♠ ♣ ♥ ♦ ← ↑ ‍ → ↓ ↔ « » ‹ › ◊
		¡ ¿ € £ ¤ ¥ ¢ ‰ ¶ “ ” „ ‌ ¦ ‡ † § © ™ ® ¹ ² ³ ¼ ½ ¾
		· • ª º ¨ × ÷ − √ ∞ ∩ ∫ ± ¬ ~ ≈ ≠ ≡ ◊ ø Ø ≤ ≥
		Δ Ω α β π µ ð ∂ ∏ ∑ ƒ
		Æ Ç Ð Ñ Þ ßßβ æ ç ð ñ ø þ
		À Á Â Ã Ä Å È É Ê Ë Ì Í Î Ï Ò Ó Ô Õ Ö Ø Ù Ú Û Ü Ý
		à á â ã ä å è é ê ë ì í î ï ò ó ô õ ö ù ú û ý ÿ
		Ξ ξ
		Φφ Phi = Golden Ratio
		τ tau = Golden Ratio
		А Б В Г Д Е Ё Ж З И Й К Л М Н О
		Á á Â â ´ Æ æ À à ℵ Α α & ∧
		汉字
		😊 🍔
	`
)

func getNewAESEncryptionService() *AESEncryptionServiceV1 {
	aesService, _ := newAESEncryptionServiceV1(getOpts())
	return aesService
}

func TestBasicEncryptAndDecryptByteArray(t *testing.T) {
	aesEncryptionService := getNewAESEncryptionService()
	// encrypt string
	encryptedBytes, _ := aesEncryptionService.Encrypt([]byte(testString))
	// decrypt string back to what it was
	decryptedBytes, _ := aesEncryptionService.Decrypt(encryptedBytes)
	assert.Equal(t, []byte(testString), decryptedBytes)
}

func TestHandleEmptyStringEncryptAndDecryptString(t *testing.T) {
	aesEncryptionService := getNewAESEncryptionService()
	// encrypt string
	encryptedBytes, err := aesEncryptionService.Encrypt([]byte(emptyString))
	assert.Equal(t, err, nil)

	// decrypt string back to what it was
	decryptedString, err := aesEncryptionService.Decrypt(encryptedBytes)
	assert.Equal(t, err, nil)
	assert.Equal(t, []byte(emptyString), decryptedString)
}

func TestCanDecryptOnNewAesInstance(t *testing.T) {
	encryptedBytes := []byte("csLdzmyY8H7-iGXVWCWKEQAETs1gUFhqBdxNH8BG_raM") // "hello" string encrypted
	aesEncryptionService := getNewAESEncryptionService()
	decryptedBytes, _ := aesEncryptionService.Decrypt(encryptedBytes)
	assert.Equal(t, testString, string(decryptedBytes))
}

func TestCanDecryptOnNewAesInstanceByteArray(t *testing.T) {
	encryptedString := "csLdzmyY8H7-iGXVWCWKEQAETs1gUFhqBdxNH8BG_raM" // "hello" string encrypted
	aesEncryptionService := getNewAESEncryptionService()
	decryptedString, _ := aesEncryptionService.Decrypt([]byte(encryptedString))
	assert.Equal(t, testString, string(decryptedString))
}

func TestWeirdCharsEncryptAndDecryptByteArray(t *testing.T) {
	aesEncryptionService := getNewAESEncryptionService()
	// encrypt string
	encryptedBytes, _ := aesEncryptionService.Encrypt([]byte(weirdCharsString))
	// decrypt string back to what it was
	decryptedBytes, _ := aesEncryptionService.Decrypt(encryptedBytes)
	assert.Equal(t, []byte(weirdCharsString), decryptedBytes)
}

func TestShouldHandleEmptyBytesDecrypt(t *testing.T) {
	aesEncryptionService := getNewAESEncryptionService()
	// decrypt string back to what it was
	decryptedString, err := aesEncryptionService.Decrypt([]byte(emptyString))
	assert.Equal(t, err, nil)
	assert.Equal(t, decryptedString, []byte(emptyString))
}
