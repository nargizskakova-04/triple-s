package handlers

import (
	"os"
	"path/filepath"
	"strings"
	"time"
	"triple-s/utils"
)

func parseBucketAndObject(path string) (string, string) {
	parts := strings.SplitN(strings.Trim(path, "/"), "/", 2)
	if len(parts) < 2 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

func updateObjectMetadata(bucketPath, objectKey string, metadata []string) error {
	metadataFile := filepath.Join(bucketPath, "objects.csv")
	records, err := utils.ReadCSVFile(metadataFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var updatedRecords [][]string
	found := false

	for _, record := range records {
		if record[0] == objectKey {
			updatedRecords = append(updatedRecords, metadata) 
			found = true
		} else {
			updatedRecords = append(updatedRecords, record)
		}
	}

	if !found {
		updatedRecords = append(updatedRecords, metadata) 
	}

	return utils.WriteCSVFile(metadataFile, updatedRecords)
}

func updateBucketLastModified(baseDir, bucketName string, lastModified time.Time) error {
	csvPath := filepath.Join(baseDir, "buckets.csv")
	records, err := utils.ReadCSVFile(csvPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var updatedRecords [][]string
	for _, record := range records {
		if record[0] == bucketName {
			record[2] = lastModified.Format(time.RFC3339) 
		}
		updatedRecords = append(updatedRecords, record)
	}

	return utils.WriteCSVFile(csvPath, updatedRecords)
}

func updateBucketMetadata(baseDir, bucketName string, lastModified time.Time, status string) error {
	csvPath := filepath.Join(baseDir, "buckets.csv")
	records, err := utils.ReadCSVFile(csvPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var updatedRecords [][]string
	for _, record := range records {
		if record[0] == bucketName {
			record[2] = lastModified.Format(time.RFC3339) 
			record[3] = status
		}
		updatedRecords = append(updatedRecords, record)
	}

	return utils.WriteCSVFile(csvPath, updatedRecords)
}

func getContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".txt":
		return "text/plain"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

func removeObjectMetadata(bucketPath, objectKey string) error {
	metadataFile := filepath.Join(bucketPath, "objects.csv")
	records, err := utils.ReadCSVFile(metadataFile)
	if err != nil {
		return err
	}

	var updatedRecords [][]string
	for _, record := range records {
		if record[0] != objectKey { 
			updatedRecords = append(updatedRecords, record)
		}
	}

	return utils.WriteCSVFile(metadataFile, updatedRecords)
}

func isBucketEmpty(bucketPath string) (bool, error) {
	entries, err := os.ReadDir(bucketPath)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.Name() != "objects.csv" && !entry.IsDir() {
			return false, nil 
		}
	}
	return true, nil
}
