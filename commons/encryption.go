package commons

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func GenerateRSAKeys() error {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	keyDir := filepath.Join(homeDir, ".envsync")

	if err := os.MkdirAll(keyDir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	privateKeyPath := filepath.Join(keyDir, "private_key.pem")
	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %w", err)
	}
	defer privateKeyFile.Close()

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	// Export the public key
	publicKeyPath := filepath.Join(keyDir, "public_key.pem")
	publicKeyFile, err := os.Create(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to create public key file: %w", err)
	}
	defer publicKeyFile.Close()

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	return nil
}
func EncryptFileWithPublicKey(filePath string) (string, error) {
	privateKeyPath := viper.GetString("envsync.private_key")
	fmt.Println(privateKeyPath)

	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read private key file: %w", err)
	}

	block, _ := pem.Decode([]byte(privateKeyBytes))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, &privateKey.PublicKey, data)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %w", err)
	}

	tempFile, err := os.CreateTemp("", "encrypted")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	if _, err := tempFile.Write(encryptedData); err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}

	return tempFile.Name(), nil
}

func DecryptFileWithPrivateKey(encryptedFilePath, decryptedFilePath string) error {
	privateKeyPath := viper.GetString("envsync.private_key")

	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read private key file: %w", err)
	}

	block, _ := pem.Decode([]byte(privateKeyBytes))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	encryptedData, err := os.ReadFile(encryptedFilePath)
	if err != nil {
		return fmt.Errorf("failed to read encrypted file: %w", err)
	}

	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}

	err = os.WriteFile(decryptedFilePath, decryptedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write decrypted data to file: %w", err)
	}

	return nil
}
