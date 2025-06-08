package prom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func QueryValue(promURL, query string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/query?query=%s", promURL, query))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data struct {
		Data struct {
			Result []struct {
				Value [2]interface{}
			}
		}
	}

	// TODO: debug puprpose, remove later
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&data); err != nil {
		fmt.Printf("Error decoding JSON response: %v\nResponse body: %s, Query: %s, URL: %s\n", err, string(b), query, resp.Request.URL.String())
		return 0, err
	}
	if len(data.Data.Result) == 0 {
		return 0, fmt.Errorf("no data")
	}

	valueStr := data.Data.Result[0].Value[1].(string)
	var value float64
	fmt.Sscanf(valueStr, "%f", &value)
	return value, nil
}
