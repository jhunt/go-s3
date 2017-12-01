package main

import (
	"fmt"
	"os"

	"github.com/jhunt/go-s3"
)

func main() {
	aki := os.Getenv("S3_AKI")
	if aki == "" {
		fmt.Fprintf(os.Stderr, "!!! no S3_AKI env var set...\n")
		os.Exit(1)
	}

	key := os.Getenv("S3_KEY")
	if key == "" {
		fmt.Fprintf(os.Stderr, "!!! no S3_KEY env var set...\n")
		os.Exit(1)
	}

	bkt := os.Getenv("S3_BUCKET")
	if bkt == "" {
		fmt.Fprintf(os.Stderr, "!!! no S3_BUCKET env var set...\n")
		os.Exit(1)
	}

	reg := os.Getenv("S3_REGION")
	if reg == "" {
		reg = "us-east-1"
	}

	sigv := 4
	if os.Getenv("S3_SIGV2") != "" {
		sigv = 2
	}

	if len(os.Args) != 1 {
		fmt.Fprintf(os.Stderr, "USAGE: list\n")
		os.Exit(1)
	}

	c, err := s3.NewClient(&s3.Client{
		AccessKeyID:      aki,
		SecretAccessKey:  key,
		Bucket:           bkt,
		Region:           reg,
		SignatureVersion: sigv,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "!!! unable to configure s3 client: %s\n", err)
		os.Exit(1)
	}

	files, err := c.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "!!! unable to list files in bucket '%s': %s\n", bkt, err)
		os.Exit(1)
	}

	for _, f := range files {
		fmt.Printf("- %s\n", f.Key)
	}
	os.Exit(0)
}
