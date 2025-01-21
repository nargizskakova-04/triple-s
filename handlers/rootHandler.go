package handlers

import (
	"fmt"
	"net/http"
)

func (rt *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "Object created")
	case http.MethodPut:
		fmt.Fprintln(w, "What do do. You can only use methid get with that")
	case http.MethodDelete:
		fmt.Fprintln(w, "What do do. You can only use methid get with that")
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}
