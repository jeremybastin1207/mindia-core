package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
	mindiaerr "github.com/jeremybastin1207/mindia-core/internal/error"
)

func toAbsolutePath(relativePath string) string {
	return "/" + relativePath
}

func writeError(w http.ResponseWriter, httpStatus int, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatus)
	if merr, ok := err.(*mindiaerr.Error); ok {
		res := encodeJSON(w, mindiaerr.NewApiError(*merr))
		_, err := w.Write(res.JsonContent)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, jsonContentResponse JsonEncoderResult) error {
	if jsonContentResponse.Err != nil {
		return jsonContentResponse.Err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(jsonContentResponse.JsonContent)
	return err
}

func writeBytes(w http.ResponseWriter, bytes []byte, mime string) error {
	w.Header().Set("Content-Type", mime)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(bytes)
	return err
}

func writeMessage(w http.ResponseWriter, message string) error {
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(message))
	return err
}

type JsonEncoderResult struct {
	JsonContent []byte
	Err         error
}

func encodeJSON(w http.ResponseWriter, obj interface{}) JsonEncoderResult {
	jsonContent, err := json.MarshalIndent(obj, "", "	")
	return JsonEncoderResult{
		JsonContent: jsonContent,
		Err:         err,
	}
}

func bearerToken(r *http.Request, header string) (string, error) {
	rawToken := r.Header.Get(header)
	splitted := strings.SplitN(rawToken, " ", 2)
	if len(splitted) < 2 {
		return "", errors.New("token with incorrect bearer format")
	}
	token := strings.TrimSpace(splitted[1])
	return token, nil
}

func parseBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var b T
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return nil, err
	}
	return &b, nil
}

func parseQuery[T any](url *url.URL) (*T, error) {
	var (
		params  T
		decoder = schema.NewDecoder()
	)
	u, err := url.Parse(url.RawPath)
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(&params, u.Query())
	if err != nil {
		return nil, err
	}
	return &params, nil
}
