package utils

import (
    "io/ioutil"
    "io"
	"net/http"
)

func Fetch(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil { return nil, err }

	return resp.Body, nil
}

func FetchString(url string) (string, error) {
    body, err := Fetch(url)
    if err != nil { return "", err }

    bytes, err := ioutil.ReadAll(body)
    if err != nil { return "", err }

	return string(bytes), nil
}
