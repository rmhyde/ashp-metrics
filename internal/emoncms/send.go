package emoncms

import "net/http"

func (c *Client) Send(req *http.Request) (*http.Response, error) {
	// Add API key to the request headers
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Use the http client from the Client struct
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
