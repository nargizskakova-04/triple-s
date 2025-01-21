package utils

import (
	"encoding/csv"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var bucketNameRegex = regexp.MustCompile(`^[a-z0-9]([a-z0-9-\\.]{1,61}[a-z0-9])?$`)

func ValidateBucketName(name string) error {
	if len(name) < 3 || len(name) > 63 {
		return errors.New("bucket name must be between 3 and 63 characters")
	}

	if !bucketNameRegex.MatchString(name) {
		return errors.New("bucket name can only contain lowercase letters, numbers, hyphens, and dots, and must not start or end with a hyphen")
	}

	if strings.Contains(name, "..") || strings.Contains(name, "--") || name[0] == '-' || name[len(name)-1] == '-' {
		return errors.New("bucket name must not contain consecutive periods or dashes, and must not start or end with a hyphen")
	}

	if isIPAddress(name) {
		return errors.New("bucket name must not be formatted as an IP address")
	}

	return nil
}

func isIPAddress(name string) bool {
	parts := strings.Split(name, ".")

	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return false
		}

		if len(part) > 1 && part[0] == '0' {
			return false
		}
	}

	return true
}

func CheckBucketUniqueness(name, csvPath string) (bool, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return false, err
	}

	for _, record := range records {
		if len(record) > 0 && record[0] == name {
			return false, nil
		}
	}

	return true, nil
}
