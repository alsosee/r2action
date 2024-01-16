package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2 is a struct describing r2 cloudflare storage bucket.
type R2 struct {
	Bucket string
	client *s3.Client
	ctx    context.Context
}

// NewR2 creates new R2 struct.
func NewR2(
	accountID string,
	accessKeyID string,
	accessKeySecret string,
	bucket string,
) (*R2, error) {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("creating config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	return &R2{
		Bucket: bucket,
		client: client,
		ctx:    context.Background(),
	}, nil
}

// Get downloads object from given key.
func (r2 *R2) Get(key string) ([]byte, error) {
	output, err := r2.client.GetObject(r2.ctx, &s3.GetObjectInput{
		Bucket: aws.String(r2.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("getting object: %w", err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(output.Body)
	if err != nil {
		return nil, fmt.Errorf("reading object body: %w", err)
	}

	return buf.Bytes(), nil
}

// Put uploads given body to given key.
func (r2 *R2) Put(key string, body []byte) error {
	_, err := r2.client.PutObject(r2.ctx, &s3.PutObjectInput{
		Bucket: aws.String(r2.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})
	if err != nil {
		return fmt.Errorf("uploading object: %w", err)
	}

	return nil
}

// Delete deletes object from given key.
func (r2 *R2) Delete(key string) error {
	_, err := r2.client.DeleteObject(r2.ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r2.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("deleting object: %w", err)
	}

	return nil
}
