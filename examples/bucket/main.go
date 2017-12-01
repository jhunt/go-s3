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

	reg := os.Getenv("S3_REGION")
	/* region can be blank - no location constraint */

	sigv := 4
	if os.Getenv("S3_SIGV2") != "" {
		sigv = 2
	}

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "USAGE: bucket NAME ACL\n")
		os.Exit(1)
	}
	bkt := os.Args[1]
	acl := os.Args[2]

	c, err := s3.NewClient(&s3.Client{
		AccessKeyID:      aki,
		SecretAccessKey:  key,
		Region:           "us-east-1", /* AWS requires bucket creation to go to us-east-1 */
		SignatureVersion: sigv,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "!!! unable to configure s3 client: %s\n", err)
		os.Exit(1)
	}

	err = c.CreateBucket(bkt, reg, acl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "!!! unable to create bucket '%s': %s\n", bkt, err)
		os.Exit(1)
	}

	err = c.DeleteBucket(bkt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "!!! unable to delete bucket '%s': %s\n", bkt, err)
		os.Exit(1)
	}

	os.Exit(0)
}
