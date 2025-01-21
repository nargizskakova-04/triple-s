package handlers

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"triple-s/utils"
)

func (b *BucketHandler) appendBucketMetadata(bucket Bucket) error {
	metadataFile := filepath.Join(b.BaseDir, "buckets.csv") 
	file, err := os.OpenFile(metadataFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		bucket.Name,
		bucket.CreationTime.Format(time.RFC3339),
		bucket.LastModifiedTime.Format(time.RFC3339),
		bucket.Status,
	}
	return writer.Write(record)
}

func (b *BucketHandler) removeBucketMetadata(bucketName string) error {
	metadataFile := filepath.Join(b.BaseDir, "buckets.csv")

	records, err := utils.ReadCSVFile(metadataFile)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %v", err)
	}

	tempFile, err := os.CreateTemp(b.BaseDir, "buckets_temp_*.csv")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	tempFilePath := tempFile.Name()
	defer func() {
		tempFile.Close()
		if err != nil {
			os.Remove(tempFilePath)
		}
	}()

	writer := csv.NewWriter(tempFile)
	defer writer.Flush()

	for _, record := range records {
		if record[0] != bucketName {
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("failed to write record: %v", err)
			}
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("failed to flush writer: %v", err)
	}

	tempFile.Close()

	if err := os.Rename(tempFilePath, metadataFile); err != nil {
		return fmt.Errorf("failed to rename temporary file: %v", err)
	}

	return nil
}
