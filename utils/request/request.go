package req

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	"github.com/codeshelldev/secured-signal-api/utils/query"
)

const (
	Json BodyType = "Json"
	Form  BodyType = "Form"
	Unknown BodyType = "Unknown"
)

type BodyType string

type Body struct {
	Data	map[string]interface{}
	Raw     []byte
	Empty	bool
}

func (body Body) ToString() string {
	return string(body.Raw)
}

func CreateBody(data map[string]interface{}) (Body, error) {
	if len(data) <= 0 {
		err := errors.New("empty data map")
		log.Error("Could not encode Body: ", err.Error())
		return Body{Empty: true}, err
	}

	bytes, err := json.Marshal(data)

	if err != nil {
		log.Error("Could not encode Body: ", err.Error())
		return Body{Empty: true}, err
	}

	isEmpty := len(data) <= 0

	return Body{
		Data: data,
		Raw: bytes,
		Empty: isEmpty,
	}, nil
}

func GetJsonData(body []byte) (map[string]interface{}, error) {
	var data map[string]interface{}

	err := json.Unmarshal(body, &data)

	if err != nil {
		log.Error("Could not decode Body: ", err.Error())
		return nil, err
	}

	return data, nil
}

func GetFormData(body []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	queryData := query.ParseRawQuery(string(body))

	if len(queryData) <= 0 {
		err := errors.New("invalid form data")
		log.Error("Could not decode Body: ", err.Error())
		return nil, err
	}

	for key, value := range queryData {	
		data[key] = query.ParseTypedQuery(value)
	}

	return data, nil
}

func GetBody(req *http.Request) ([]byte, error) {
	bodyBytes, err := io.ReadAll(req.Body)
	
	if err != nil {
		log.Error("Could not read Body: ", err.Error())

		req.Body.Close()

		return nil, err
	}
	defer req.Body.Close()

	return bodyBytes, nil
}

func GetReqBody(w http.ResponseWriter, req *http.Request) Body {
	bytes, err := GetBody(req)

	var isEmpty bool
	
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		isEmpty = true
	}

	if len(bytes) <= 0 {
		isEmpty = true
	}

	if isEmpty {
		return Body{Empty: true}
	}

	var data map[string]interface{}

	switch GetBodyType(req) {
	case Json:
		data, err = GetJsonData(bytes)

		if err != nil {
			http.Error(w, "Bad Request: invalid JSON", http.StatusBadRequest)
		}
	case Form:
		data, err = GetFormData(bytes)

		if err != nil {
			http.Error(w, "Bad Request: invalid Form", http.StatusBadRequest)
		}
	}

	isEmpty = len(data) <= 0

	return Body{
		Raw: bytes,
		Data: data,
		Empty: isEmpty,
	}
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