package Response

import "net/http"

func SendError(w http.ResponseWriter, message string, code int) {
	http.Error(w, message, code)
}

type ErrorResponse struct {
	Errors map[string]string `json:"errors"`
}
