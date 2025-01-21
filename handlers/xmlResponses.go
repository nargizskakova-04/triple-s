package handlers

import (
	"encoding/xml"
	"net/http"
)

type XMLErrorResponse struct {
	XMLName xml.Name `xml:"Error"`
	Code    int      `xml:"Code"`
	Message string   `xml:"Message"`
}

func WriteXMLError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(statusCode)

	errResponse := XMLErrorResponse{
		Code:    statusCode,
		Message: message,
	}

	if err := xml.NewEncoder(w).Encode(errResponse); err != nil {
		http.Error(w, "Failed to encode XML error response", http.StatusInternalServerError)
	}
}
