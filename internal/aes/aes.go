package aes

import (
	"github.com/golang-module/dongle"
)

func GetCypher(key string) *dongle.Cipher {
	b := dongle.Encrypt.FromString(key).BySha256().ToRawBytes()
	cipher := dongle.NewCipher()
	cipher.SetMode(dongle.CBC)
	cipher.SetPadding(dongle.PKCS7)
	cipher.SetKey(b)
	cipher.SetIV("9859102938491658")
	return cipher
}

func Encrypt(cipher *dongle.Cipher, message string) string {
	return dongle.Encrypt.FromString(message).ByAes(cipher).ToBase64String()
}

func Decrypt(cipher *dongle.Cipher, message string) string {
	return dongle.Decrypt.FromBase64String(message).ByAes(cipher).ToString()
}
