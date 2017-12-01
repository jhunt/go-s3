package s3

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type Object struct {
	Key string
}

func (c *Client) List() ([]Object, error) {
	objects := make([]Object, 0)
	ctok := ""
	for {
		res, err := c.get(fmt.Sprintf("/?list-type=2%s", ctok), nil)
		if err != nil {
			return nil, err
		}

		var r struct {
			XMLName  xml.Name `xml:"ListBucketResult"`
			Next     string   `xml:"NextContinuationToken"`
			Contents []struct {
				Key string `xml:"Key"`
			} `xml:"Contents"`
		}
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if res.StatusCode != 200 {
			return nil, ResponseErrorFrom(b)
		}

		err = xml.Unmarshal(b, &r)
		if err != nil {
			return nil, err
		}

		for _, f := range r.Contents {
			objects = append(objects, Object{
				Key: f.Key,
			})
		}

		if r.Next == "" {
			return objects, nil
		}

		ctok = fmt.Sprintf("&continuation-token=%s", r.Next)
	}
}
