// Пакет отвечает за шифрование/дешифрование данных.
package mycrypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/KirillKhitev/metrics/internal/flags"
)

// Encrypting шифрует тело запроса используя публичный ключ
func Encrypting(data []byte, keyFilePath string) ([]byte, error) {
	result := []byte{}

	if keyFilePath == "" {
		return data, nil
	}

	publicKeyData, err := os.ReadFile(keyFilePath)
	if err != nil {
		return result, fmt.Errorf("error by read public key file, error: %w", err)
	}

	block, _ := pem.Decode(publicKeyData)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return result, errors.New("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return result, err
	}

	result, err = rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return result, fmt.Errorf("error from encryption: %w", err)
	}

	return result, nil
}

// Decripting расшифровывает тело запроса используя приватный ключ.
func Decrypting(data []byte, keyFilePath string) ([]byte, error) {
	result := []byte{}

	if keyFilePath == "" {
		return data, nil
	}

	privateKeyData, err := os.ReadFile(keyFilePath)
	if err != nil {
		return result, fmt.Errorf("error by read private key file, error: %w", err)
	}

	block, _ := pem.Decode(privateKeyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return result, errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return result, err
	}

	result, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)

	if err != nil {
		return result, fmt.Errorf("error from decripting: %w", err)
	}

	return result, nil
}

// Middleware, расшифровываем данные из запроса.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)

		bodyDecrypted, err := Decrypting(body, flags.Args.CryptoKey)
		if err != nil {
			log.Fatalf(err.Error())
		}

		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(bodyDecrypted))

		next.ServeHTTP(w, r)
	})
}
