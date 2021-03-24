package platerecognizer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/sno6/gate-god/recognition"
)

const (
	baseURL = "https://api.platerecognizer.com/v1/plate-reader"
	token   = "<TOKEN>"
)

type Recognizer struct{}

func New() *Recognizer {
	return &Recognizer{}
}

func (rr *Recognizer) Recognize(r io.Reader) (*recognition.Result, error) {
	body, boundary, err := rr.buildPostBody(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", boundary)
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	var result *Response
	err = json.NewDecoder(rsp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if len(result.Results) < 1 {
		return nil, errors.New("recogizer: no results")
	}

	return &recognition.Result{
		Plate: result.Results[0].Plate,
		Score: result.Results[0].Score,
	}, nil
}

func (rr *Recognizer) buildPostBody(r io.Reader) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, "", err
	}

	part, err := writer.CreateFormFile("upload", "upload.jpg")
	if err != nil {
		return nil, "", err
	}
	part.Write(contents)

	if err := writer.WriteField("regions", "nz"); err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}
