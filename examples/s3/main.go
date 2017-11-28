package main

import (
	"fmt"
	"io"
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

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "USAGE: s3 path/in/bucket <file >out\n")
		os.Exit(1)
	}
	path := os.Args[1]

	c, err := s3.NewClient(&s3.Client{
		AccessKeyID:     aki,
		SecretAccessKey: key,
		Region:          reg,
		Bucket:          bkt,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "!!! unable to configure s3 client: %s\n", err)
		os.Exit(1)
	}

	u, err := c.NewUpload(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "!!! unable to start multipart upload: %s\n", err)
		os.Exit(1)
	}

	n, err := u.Stream(os.Stdin, 5*1024*1024*1024)
	if err != nil {
		fmt.Fprintf(os.Stderr, "!!! unable to stream <stdin> in 5m parts: %s\n", err)
		os.Exit(1)
	}

	err = u.Done()
	if err != nil {
		fmt.Fprintf(os.Stderr, "!!! unable to complete multipart upload: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "wrote %d bytes\n", n)

	out, err := c.Get(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "!!! unable to retrieve our uploaded file: %s\n", err)
		os.Exit(1)
	}
	io.Copy(os.Stdout, out)

	err = c.Delete(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "!!! unable to remove our uploaded file: %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
