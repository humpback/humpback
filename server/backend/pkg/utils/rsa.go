package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"humpback/common/enum"
)

func RSAEncrypt(text string) string {
	pubBlock, _ := pem.Decode([]byte(enum.PublicKey))
	pubKeyValue, _ := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	pub := pubKeyValue.(*rsa.PublicKey)
	ciphertext, _ := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(text))
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func RSADecrypt(base64Text string) string {
	ciphertext, _ := base64.StdEncoding.DecodeString(base64Text)
	priBlock, _ := pem.Decode([]byte(enum.PrivateKey))
	priKey, _ := x509.ParsePKCS1PrivateKey(priBlock.Bytes)
	text, _ := rsa.DecryptPKCS1v15(rand.Reader, priKey, ciphertext)
	return string(text)
}
