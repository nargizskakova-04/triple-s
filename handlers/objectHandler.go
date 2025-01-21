package handlers

import (
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func (o *ObjectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		o.GetObject(w, r)
	case http.MethodPut:
		o.UploadObject(w, r)
	case http.MethodDelete:
		o.DeleteObject(w, r)
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

func (o *ObjectHandler) UploadObject(w http.ResponseWriter, r *http.Request) {
	bucketName, objectKey := parseBucketAndObject(r.URL.Path)
	bucketPath := filepath.Join(o.BaseDir, bucketName)

	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		WriteXMLError(w, http.StatusNotFound, "Bucket does not exist")
		return
	}

	objectPath := filepath.Join(bucketPath, objectKey)
	file, err := os.Create(objectPath)
	if err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Could not save object")
		return
	}
	defer file.Close()

	size, err := io.Copy(file, r.Body)
	if err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Failed to save object")
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	lastModified := time.Now()
	objectMetadata := []string{
		objectKey,
		strconv.FormatInt(size, 10),
		contentType,
		lastModified.Format(time.RFC3339),
	}

	if err := updateObjectMetadata(bucketPath, objectKey, objectMetadata); err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Failed to update object metadata")
		return
	}

	if err := updateBucketMetadata(o.BaseDir, bucketName, lastModified, "active"); err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Failed to update bucket metadata")
		return
	}

	response := struct {
		XMLName xml.Name `xml:"UploadObjectResponse"`
		Object  Object   `xml:"Object"`
	}{
		Object: Object{
			ObjectKey:    objectKey,
			Size:         int(size),
			ContentType:  contentType,
			LastModified: lastModified,
		},
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(response)
}

func (o *ObjectHandler) GetObject(w http.ResponseWriter, r *http.Request) {
	bucketName, objectKey := parseBucketAndObject(r.URL.Path)
	bucketPath := filepath.Join(o.BaseDir, bucketName)

	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		WriteXMLError(w, http.StatusNotFound, "Bucket does not exist")
		return
	}

	objectPath := filepath.Join(bucketPath, objectKey)
	file, err := os.Open(objectPath)
	if err != nil {
		WriteXMLError(w, http.StatusNotFound, "Object not found")
		return
	}
	defer file.Close()

	contentType := getContentType(objectPath)
	w.Header().Set("Content-Type", contentType)

	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, file); err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Failed to retrieve object")
	}
}

func (o *ObjectHandler) DeleteObject(w http.ResponseWriter, r *http.Request) {
	bucketName, objectKey := parseBucketAndObject(r.URL.Path)
	bucketPath := filepath.Join(o.BaseDir, bucketName)

	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		WriteXMLError(w, http.StatusNotFound, "Bucket does not exist")
		return
	}

	objectPath := filepath.Join(bucketPath, objectKey)
	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		WriteXMLError(w, http.StatusNotFound, "Object not found")
		return
	}

	if err := os.Remove(objectPath); err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Failed to delete object")
		return
	}

	if err := removeObjectMetadata(bucketPath, objectKey); err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Failed to update object metadata")
		return
	}

	empty, err := isBucketEmpty(bucketPath)
	if err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Failed to verify bucket status")
		return
	}

	newStatus := "active"
	if empty {
		newStatus = "marked for deletion"
	}
	if err := updateBucketMetadata(o.BaseDir, bucketName, time.Now(), newStatus); err != nil {
		WriteXMLError(w, http.StatusInternalServerError, "Failed to update bucket metadata")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
