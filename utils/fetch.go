package utils

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

func Fetch(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil { return nil, err }
    if resp.StatusCode != 200 {
        return nil, errors.New("Got a bad status code " + strconv.Itoa(resp.StatusCode) + " for url: " + url)
    }
	return resp.Body, nil
}

func FetchString(url string) (string, error) {
    body, err := Fetch(url)
    if err != nil { return "", err }

    bytes, err := ioutil.ReadAll(body)
    if err != nil { return "", err }

	return string(bytes), nil
}
