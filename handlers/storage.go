package handlers

import (
	"time"
)

type BucketHandler struct {
	BaseDir string
}

type Bucket struct {
	Name             string    `xml:"Name"`
	CreationTime     time.Time `xml:"LastModified"`
	LastModifiedTime time.Time `xml:"CreationDate"`
	Status           string    `xml:"Status"`
}

type ObjectHandler struct {
	BaseDir string
}

type Object struct {
	ObjectKey    string    `xml:"ObjectKey"`
	Size         int       `xml:"Size"`
	ContentType  string    `xml:"ContentType"`
	LastModified time.Time `xml:"LastModified"`
}

type Storage struct {
	Bucket []Bucket
	Object []Object
}

type RootHandler struct{}
