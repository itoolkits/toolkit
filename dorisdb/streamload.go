// doris stream load

package dorisdb

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/itoolkits/toolkit/retry"
)

type StreamLoad struct {
	URL      string
	Port     string
	DB       string
	Table    string
	Username string
	Password string

	Timeout time.Duration
}

// auth - doris auth
func (s *StreamLoad) auth() string {
	return base64.StdEncoding.EncodeToString([]byte(s.Username + ":" + s.Password))
}

// LoadData - doris stream load
func (s *StreamLoad) LoadData(data []byte) error {
	client := &http.Client{
		Timeout: s.Timeout,
	}
	url := "http://%s:%s/api/%s/%s/_stream_load"
	url = fmt.Sprintf(url, s.URL, s.Port, s.DB, s.Table)

	reader := bytes.NewReader(data)

	request, err := http.NewRequest(http.MethodPut, url, reader)
	if err != nil {
		return err
	}

	uid, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("doris load data label mark error, %w", err)
	}

	request.Header.Add("Authorization", "Basic "+s.auth())
	request.Header.Add("EXPECT", "100-continue")
	request.Header.Add("label", uid.URN())
	// request.Header.Add("column_separator", ",")
	request.Header.Add("timeout", "10")
	request.Header.Add("format", "json")
	request.Header.Add("strip_outer_array", "true")

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if response == nil {
		return fmt.Errorf("doris response nil error")
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("doris response status error, %d %v", response.StatusCode, response)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if len(body) < 1 {
		return fmt.Errorf("doris response body nil error")
	}

	defer response.Body.Close()

	responseBody := &DorisResponseBody{}
	err = json.Unmarshal(body, responseBody)
	if err != nil {
		return fmt.Errorf("doris response unmarshal error, %s %s", string(body), err.Error())
	}
	if responseBody.Status != "Success" && responseBody.Status != "Publish Timeout" {
		return fmt.Errorf("doris response status error, %s", string(body))
	}
	if responseBody.NumberFilteredRows > 0 {
		return fmt.Errorf("doris response number filter rows error, %s", string(body))
	}
	return nil
}

// LoadDataWithRetry - load data with retry
func (s *StreamLoad) LoadDataWithRetry(data []byte, times int) error {
	return retry.Do(func() error {
		return s.LoadData(data)
	}, times)
}

type DorisResponseBody struct {
	TxnID                  int    `json:"TxnId"`
	Label                  string `json:"Label"`
	Status                 string `json:"Status"`
	Message                string `json:"Message"`
	NumberTotalRows        int    `json:"NumberTotalRows"`
	NumberLoadedRows       int    `json:"NumberLoadedRows"`
	NumberFilteredRows     int    `json:"NumberFilteredRows"`
	NumberUnselectedRows   int    `json:"NumberUnselectedRows"`
	LoadBytes              int    `json:"LoadBytes"`
	LoadTimeMs             int    `json:"LoadTimeMs"`
	BeginTxnTimeMs         int    `json:"BeginTxnTimeMs"`
	StreamLoadPutTimeMs    int    `json:"StreamLoadPutTimeMs"`
	ReadDataTimeMs         int    `json:"ReadDataTimeMs"`
	WriteDataTimeMs        int    `json:"WriteDataTimeMs"`
	CommitAndPublishTimeMs int    `json:"CommitAndPublishTimeMs"`
	ErrorURL               string `json:"ErrorURL"`
}
