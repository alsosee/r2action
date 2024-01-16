// GitHub Action that does something
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v9"
)

// type vars map[string]interface{}

type conf struct {
	AccountID       string `env:"INPUT_ACCOUNT_ID" long:"account-id" description:"R2 account id"`
	AccessKeyID     string `env:"INPUT_ACCESS_KEY_ID" long:"access-key-id" description:"R2 access key id"`
	AccessKeySecret string `env:"INPUT_ACCESS_KEY_SECRET" long:"access-key-secret" description:"R2 access key secret"`
	Bucket          string `env:"INPUT_BUCKET" long:"bucket" description:"R2 bucket"`
	Operation       string `env:"INPUT_OPERATION" long:"operation" description:"Operation to perform"`
	Key             string `env:"INPUT_KEY" long:"key" description:"Object key"`
	File            string `env:"INPUT_FILE" long:"file" description:"Local file path"`
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("::error::%v", err)
		os.Exit(1)
	}
}

func run() error {
	var c conf

	if err := env.Parse(&c); err != nil {
		return err
	}

	r2, err := NewR2(
		c.AccountID,
		c.AccessKeyID,
		c.AccessKeySecret,
		c.Bucket,
	)
	if err != nil {
		return fmt.Errorf("failed to create R2 client: %v", err)
	}

	err = performOperation(r2, c.Operation, c.Key, c.File)
	if err != nil {
		return fmt.Errorf("failed to perform operation: %v", err)
	}

	return nil
}

func performOperation(r2 *R2, operation, key, file string) error {
	switch operation {
	case "get":
		return get(r2, key, file)
	case "put":
		return put(r2, key, file)
	case "delete":
		return del(r2, key)
	default:
		return fmt.Errorf("unknown operation %q", operation)
	}
}

func get(r2 *R2, key, file string) error {
	if key == "" {
		return fmt.Errorf("key is required")
	}

	if file == "" {
		return fmt.Errorf("file is required")
	}

	body, err := r2.Get(key)
	if err != nil {
		return fmt.Errorf("downloading object: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	if err := os.WriteFile(file, body, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

func put(r2 *R2, key, file string) error {
	if key == "" {
		return fmt.Errorf("key is required")
	}

	if file == "" {
		return fmt.Errorf("file is required")
	}

	body, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	if err := r2.Put(key, body); err != nil {
		return fmt.Errorf("uploading object: %w", err)
	}

	return nil
}

func del(r2 *R2, key string) error {
	if key == "" {
		return fmt.Errorf("key is required")
	}

	if err := r2.Delete(key); err != nil {
		return fmt.Errorf("deleting object: %w", err)
	}

	return nil
}
