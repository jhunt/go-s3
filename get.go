package s3

import (
	"io"
)

func (c *Client) Get(key string) (io.Reader, error) {
	res, err := c.get(key)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, ResponseError(res)
	}

	return res.Body, nil
}
