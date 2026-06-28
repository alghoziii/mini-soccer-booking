package storage

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

const uploadTimeout = 60 * time.Second

type S3Config struct {
	Region          string
	BucketName      string
	BaseURL         string
	AccessKeyID     string
	SecretAccessKey string
}

type Client struct {
	s3Client   *s3.Client
	bucketName string
	baseURL    string
}

type IClient interface {
	UploadFile(context.Context, string, []byte) (string, error)
}

func NewS3Client(ctx context.Context, s3Config S3Config) (IClient, error) {
	cfg, err := awsConfig.LoadDefaultConfig(
		ctx,
		awsConfig.WithRegion(s3Config.Region),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			s3Config.AccessKeyID,
			s3Config.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	return &Client{
		s3Client:   s3.NewFromConfig(cfg),
		bucketName: s3Config.BucketName,
		baseURL:    strings.TrimRight(s3Config.BaseURL, "/"),
	}, nil
}

func (c *Client) UploadFile(ctx context.Context, filename string, data []byte) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, uploadTimeout)
	defer cancel()

	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/octet-stream"),
	})
	if err != nil {
		logrus.Errorf("failed to upload file to s3: %v", err)
		return "", fmt.Errorf("upload file to s3: %w", err)
	}

	return fmt.Sprintf("%s/%s", c.baseURL, filename), nil
}
