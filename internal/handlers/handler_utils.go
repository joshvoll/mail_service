package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type intrfc interface{}

// httpResponse standar json response
type httpResponse struct {
	Error string `json:"error,omitempty"`
	Data  intrfc `json:"data,omitempty"`
}

// jsonConverter converts json data into model
func jsonConverter(data io.ReadCloser, model interface{}) error {
	decoder := json.NewDecoder(data)
	err := decoder.Decode(model)
	return err
}

// jsonResponse sends standar json response
func jsonResponse(w http.ResponseWriter, data interface{}) {
	response := &httpResponse{
		Data: data,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// jsonError sends json response with specified error
func jsonErrorResponse(w http.ResponseWriter, err error, data interface{}) {
	response := &httpResponse{
		Error: err.Error(),
		Data:  data,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// httpRequest creates request struct
type httpRequest struct {
	Response *http.Response
	URL      string
	Data     []byte
	Method   string
	request  *http.Request
	Headers  map[string]string
}

// Send Sends http request
func (r *httpRequest) Send() error {
	httpClient := http.Client{}

	var err error

	r.request, err = http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.Data))
	if err != nil {
		return err
	}

	// Add headers to the request
	for index, header := range r.Headers {
		r.request.Header.Add(index, header)
	}

	// make requests and store response
	r.Response, err = httpClient.Do(r.request)

	return err
}
