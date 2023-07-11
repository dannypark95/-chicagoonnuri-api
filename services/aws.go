package services

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	awsConfig *aws.Config
	bucketName string
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("unable to load AWS configuration")
	}

	awsConfig = &cfg
	bucketName = "chicagoonnuri"
}

// UploadToS3 uploads a file to the specified S3 bucket with the provided key and content type.
func UploadToS3(file io.Reader, key string, contentType string) (string, error) {
	client := s3.NewFromConfig(*awsConfig)

	input := &s3.PutObjectInput{
		Body: file,
		Bucket: aws.String(bucketName),
		Key: aws.String(key),
		ContentType: aws.String(contentType),
	}

	_, err := client.PutObject(context.Background(), input)
	if err != nil {
		return "", err
	}

	return GetObjectURL(bucketName, key), nil
}

// ListPDFs lists all PDF files in the "jubo" folder in the S3 bucket.
func ListPDFs() ([]string, error) {
	client := s3.NewFromConfig(*awsConfig)
	prefix := "jubo/"

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	}

	output, err := client.ListObjectsV2(context.Background(), input)
	if err != nil {
		return nil, err
	}

	fileURLs := make([]string, len(output.Contents))
	for i, object := range output.Contents {
		fileURLs[i] = GetObjectURL(bucketName, *object.Key)
	}

	return fileURLs, nil
}

//ReadLiveJuboFromS3 gets the live jubo from the metadata file in S3.
func ReadLiveJuboFromS3() (string, error) {
	client := s3.NewFromConfig(*awsConfig)

	input := &s3.GetObjectInput {
		Bucket: aws.String(bucketName),
		Key: aws.String("jubo_metadata.json"),
	}

	resp, err := client.GetObject(context.Background(), input)
	if err != nil {
	  log.Println(err)
	  return "", err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.String(), nil
}

// WriteLiveJuboToS3 sets the live jubo in the metadata files in S3.
func WriteLiveJuboToS3(liveJubo string) error {
	client := s3.NewFromConfig(*awsConfig)

	liveJuboData := []byte(`{"liveJubo": "` + liveJubo + `"}`)

	input := &s3.PutObjectInput{
		Body: strings.NewReader(string(liveJuboData)),
		Bucket: aws.String(bucketName),
		Key: aws.String("jubo_metadata.json"),
		ContentType: aws.String("application/json"),
	}

	_, err := client.PutObject(context.Background(), input)
	if err != nil {
		return err
	}

	return nil
}

// getObjectURL constructs the URL for an object in the S3 bucket with the given key.
func GetObjectURL(bucket, key string) string {
	escapedKey := url.QueryEscape(key)
	return "https://" + bucket + ".s3.amazonaws.com/" + escapedKey
}

// DeleteFromS3 deletes an object from the S3 bucket with the given key
func DeleteFromS3(key string) error {
	client := s3.NewFromConfig(*awsConfig)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}

	_, err := client.DeleteObject(context.Background(), input)
	if err != nil {
		return err
	}

	return nil
}