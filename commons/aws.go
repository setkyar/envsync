package commons

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/viper"
)

func UploadToS3(envName, filePath string) error {
	awsRegion := viper.GetString("aws.region")
	awsAccessKeyID := viper.GetString("aws.access_key_id")
	awsSecretAccessKey := viper.GetString("aws.secret_access_key")
	s3Bucket := viper.GetString("aws.s3_bucket")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %w", err)
	}

	encryptedFilePath, err := EncryptFileWithPrivateKey(filePath)
	if err != nil {
		return fmt.Errorf("failed to encrypt file: %w", err)
	}
	defer os.Remove(encryptedFilePath)

	encryptedFile, err := os.Open(encryptedFilePath)
	if err != nil {
		return fmt.Errorf("failed to open encrypted file %s: %w", encryptedFilePath, err)
	}
	defer encryptedFile.Close()

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(filepath.Join(envName, ".env")),
		Body:   encryptedFile,
	})
	if err != nil {
		return fmt.Errorf("failed to upload encrypted file to S3: %w", err)
	}

	fmt.Println("Encrypted file successfully uploaded to S3.")
	return nil
}

func DownloadFromS3(envName, destinationPath string) error {
	awsRegion := viper.GetString("aws.region")
	awsAccessKeyID := viper.GetString("aws.access_key_id")
	awsSecretAccessKey := viper.GetString("aws.secret_access_key")
	s3Bucket := viper.GetString("aws.s3_bucket")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %w", err)
	}

	downloader := s3manager.NewDownloader(sess)

	// Create a temporary file to store the encrypted data
	tempFile, err := os.CreateTemp("", "encrypted-env")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name()) // Clean up the temp file after use

	// Download the encrypted file from S3
	_, err = downloader.Download(tempFile, &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(filepath.Join(envName, ".env")),
	})
	if err != nil {
		return fmt.Errorf("failed to download file from S3: %w", err)
	}

	err = DecryptFileWithPrivateKey(tempFile.Name(), destinationPath)
	if err != nil {
		return fmt.Errorf("failed to decrypt file: %w", err)
	}

	fmt.Println("File successfully downloaded and decrypted.")
	return nil
}
