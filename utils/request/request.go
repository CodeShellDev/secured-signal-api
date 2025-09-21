package req

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/codeshelldev/secured-signal-api/utils/query"
)

const (
	Json BodyType = "Json"
	Form  BodyType = "Form"
	Unknown BodyType = "Unknown"
)

type BodyType string

type Body struct {
	Data	map[string]any
	Raw     []byte
	Empty	bool
}

func (body Body) ToString() string {
	return string(body.Raw)
}

func CreateBody(data map[string]any) (Body, error) {
	if len(data) <= 0 {
		err := errors.New("empty data map")

		return Body{Empty: true}, err
	}

	bytes, err := json.Marshal(data)

	if err != nil {

		return Body{Empty: true}, err
	}

	isEmpty := len(data) <= 0

	return Body{
		Data: data,
		Raw: bytes,
		Empty: isEmpty,
	}, nil
}

func GetJsonData(body []byte) (map[string]any, error) {
	var data map[string]any

	err := json.Unmarshal(body, &data)

	if err != nil {

		return nil, err
	}

	return data, nil
}

func GetFormData(body []byte) (map[string]any, error) {
	data := map[string]any{}

	queryData := query.ParseRawQuery(string(body))

	if len(queryData) <= 0 {
		err := errors.New("invalid form data")

		return nil, err
	}

	for key, value := range queryData {	
		data[key] = query.ParseTypedQueryValues(value)
	}

	return data, nil
}

func GetBody(req *http.Request) ([]byte, error) {
	bodyBytes, err := io.ReadAll(req.Body)
	
	if err != nil {
		req.Body.Close()

		return nil, err
	}
	defer req.Body.Close()

	return bodyBytes, nil
}

func GetReqHeaders(req *http.Request) (map[string]any) {
	data := map[string]any{}

	for key, value := range req.Header {
		data[key] = value
	}

	return data
}

func GetReqBody(w http.ResponseWriter, req *http.Request) (Body, error) {
	bytes, err := GetBody(req)

	var isEmpty bool
	
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return Body{Empty: true}, err
	}

	if len(bytes) <= 0 {
		return Body{Empty: true}, nil
	}

	var data map[string]any

	switch GetBodyType(req) {
		case Json:
			data, err = GetJsonData(bytes)

			if err != nil {
				http.Error(w, "Bad Request: invalid JSON", http.StatusBadRequest)

				return Body{Empty: true}, err
			}
		case Form:
			data, err = GetFormData(bytes)

			if err != nil {
				http.Error(w, "Bad Request: invalid Form", http.StatusBadRequest)

				return Body{Empty: true}, err
			}
	}

	isEmpty = len(data) <= 0

	return Body{
		Raw: bytes,
		Data: data,
		Empty: isEmpty,
	}, nil
}

func GetBodyType(req *http.Request) BodyType {
	contentType := req.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(contentType, "application/json"):
		return Json

	case strings.HasPrefix(contentType, "multipart/form-data"):
		return Form

	case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
		return Form
	default:
		return Unknown
	}
}