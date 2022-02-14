package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func main() {
	var (
		bucket  string
		keyFile string
		file    string
	)

	flag.StringVar(&bucket, "b", "", "Bucket name")
	flag.StringVar(&keyFile, "k", "", "GCP key.json file")
	flag.StringVar(&file, "f", "", "File to upload, will be named the same in the bucket")

	flag.Parse()

	if bucket == "" {
		fmt.Println("Bucket Required")
		os.Exit(1)
	}

	if file == "" {
		fmt.Println("File to upload required")
		os.Exit(1)
	}

	err := uploadFile(os.Stdout, bucket, file, keyFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func uploadFile(w io.Writer, bucket, object, jsonPath string) error {
	ctx := context.Background()
	var client *storage.Client
	var err error
	if jsonPath == "" {
		client, err = storage.NewClient(ctx)
	} else {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(jsonPath))
	}
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open(object)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	wc.ProgressFunc = progressfunc

	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Fprintf(w, "Blob %v uploaded.\n", object)
	return nil
}

func progressfunc(copiedBytes int64) {
	mb := copiedBytes / 1024 / 1024
	fmt.Printf("%dMB Uploaded\n", mb)
}
