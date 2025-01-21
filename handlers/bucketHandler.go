package handlers

import (
	"encoding/xml"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
	"triple-s/utils"
)

func (b *BucketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		b.ListBuckets(w, r)
	case http.MethodPut:
		b.CreateBucket(w, r)
	case http.MethodDelete:
		b.DeleteBucket(w, r)
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

func (b *BucketHandler) CreateBucket(w http.ResponseWriter, r *http.Request) {
	if err := utils.EnsureDirExists(b.BaseDir); err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Failed to create base directory")
		return
	}

	csvPath := path.Join(b.BaseDir, "buckets.csv")
	if err := utils.EnsureFileExists(csvPath); err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Failed to create buckets.csv")
		return
	}

	bucketName := r.URL.Path[1:]

	if err := utils.ValidateBucketName(bucketName); err != nil {
		WriteXMLError(w, http.StatusBadRequest, err.Error())
		return
	}

	unique, err := utils.CheckBucketUniqueness(bucketName, csvPath)
	if err != nil {
		WriteXMLError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !unique {
		WriteXMLError(w, http.StatusConflict, "Bucket name already exists")
		return
	}

	bucketPath := filepath.Join(b.BaseDir, bucketName)
	if _, err := os.Stat(bucketPath); !os.IsNotExist(err) {
		WriteXMLError(w, http.StatusConflict, "Bucket already exists")
		return
	}

	if err := os.Mkdir(bucketPath, os.ModePerm); err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Failed to create bucket")
		return
	}

	now := time.Now()
	bucket := Bucket{
		Name:             bucketName,
		CreationTime:     now,
		LastModifiedTime: now,
		Status:           "marked for deletion",
	}

	if err := b.appendBucketMetadata(bucket); err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Failed to update metadata")
		return
	}
	response := struct {
		XMLName xml.Name `xml:"CreateBucketResponse"`
		Bucket  Bucket   `xml:"Bucket"`
	}{
		Bucket: bucket,
	}

	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(response)
}

func (b *BucketHandler) ListBuckets(w http.ResponseWriter, r *http.Request) {
	csvPath := path.Join(b.BaseDir, "buckets.csv")

	if r.URL.Path != "/" {
		WriteXMLError(w, http.StatusBadRequest, "List buckets only available at root path")
		return
	}

	records, err := utils.ReadCSVFile(csvPath)
	if err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Failed to read metadata")
		return
	}

	var buckets []Bucket
	for _, record := range records {
		createdTime, err := time.Parse(time.RFC3339, record[1])
		if err != nil {
			WriteXMLError(w, http.StatusInternalServerError, "Invalid creation time format in CSV")
			return
		}
		lastModifiedTime, err := time.Parse(time.RFC3339, record[2])
		if err != nil {
			WriteXMLError(w, http.StatusInternalServerError, "Invalid last modified time format in CSV")
			return
		}

		bucket := Bucket{
			Name:             record[0],
			CreationTime:     createdTime,
			LastModifiedTime: lastModifiedTime,
			Status:           record[3],
		}

		buckets = append(buckets, bucket)
	}

	response := struct {
		XMLName xml.Name `xml:"ListBucketsResponse"`
		Buckets []Bucket `xml:"Bucket"`
	}{
		Buckets: buckets,
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(response)
}

func (b *BucketHandler) DeleteBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Path[1:]
	if bucketName == "" || bucketName == "/" {
		WriteXMLError(w, http.StatusBadRequest, "Bucket name not specified")
		return
	}

	csvPath := filepath.Join(b.BaseDir, "buckets.csv")
	records, err := utils.ReadCSVFile(csvPath)
	if err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Failed to read bucket metadata")
		return
	}

	bucketExists := false
	for _, record := range records {
		if record[0] == bucketName {
			bucketExists = true
			break
		}
	}

	if !bucketExists {
		WriteXMLError(w, http.StatusNotFound, "Bucket does not exist in metadata")
		return
	}

	bucketPath := filepath.Join(b.BaseDir, bucketName)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		WriteXMLError(w, http.StatusNotFound, "Bucket directory does not exist")
		return
	}

	objectsCSVPath := filepath.Join(bucketPath, "objects.csv")
	if _, err := os.Stat(objectsCSVPath); err == nil {
		records, err := utils.ReadCSVFile(objectsCSVPath)
		if err != nil {
			WriteXMLError(w, http.StatusInternalServerError, "Failed to read object metadata")
			return
		}
		if len(records) > 0 {
			WriteXMLError(w, http.StatusConflict, "Cannot delete non-empty bucket")
			return
		}
	}

	if err := os.RemoveAll(bucketPath); err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Failed to delete bucket directory")
		return
	}

	if err := b.removeBucketMetadata(bucketName); err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Failed to update metadata")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
