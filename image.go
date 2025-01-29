package malak

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ayinke-llc/hermes"
)

var imageClient = &http.Client{
	Timeout: time.Second * 3,
}

func IsImageFromURL(s string) (bool, error) {
	if hermes.IsStringEmpty(s) {
		return false, errors.New("please provide the url")
	}

	u, err := url.Parse(s)
	if err != nil {
		return false, err
	}

	resp, err := imageClient.Get(u.String())
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	// Read the first 512 bytes of the response body
	buffer := make([]byte, 512)
	_, err = resp.Body.Read(buffer)
	if err != nil {
		return false, err
	}

	// Detect the MIME type using the first 512 bytes
	mimeType := http.DetectContentType(buffer)

	if strings.HasPrefix(mimeType, "image/") {
		return true, nil
	}

	return false, errors.New("url does not contain an image")
}
