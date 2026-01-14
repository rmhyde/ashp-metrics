package emoncms

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func (c *Client) InsertOrUpdateDatapoint(id int, time time.Time, value string) error {

	req, err := http.NewRequest("POST", c.BaseURL+"/feed/post.json", nil)
	q := req.URL.Query()
	q.Add("id", fmt.Sprint(id))
	q.Add("time", fmt.Sprint(time.Unix()))
	q.Add("value", value)
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
