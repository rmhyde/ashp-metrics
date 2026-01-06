package emoncms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func UpdateInput[T any](c *Client, node string, values T) error {

	jsonBytes, err := json.Marshal(values)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/input/post", nil)
	q := req.URL.Query()
	q.Add("node", node)
	q.Add("fulljson", string(jsonBytes))
	req.URL.RawQuery = q.Encode()

	if err != nil {
		return err
	}

	resp, err := c.Send(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("body: %s", string(body))
	return err
}
