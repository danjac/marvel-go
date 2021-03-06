package marvel

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "http://gateway.marvel.com:80/v1/public/"

type Client struct {
	PublicKey  string
	PrivateKey string
}

func NewClient(publicKey string, privateKey string) Client {
	c := Client{}
	c.PublicKey = publicKey
	c.PrivateKey = privateKey
	return c
}

func (client *Client) makeHash(timestamp string) string {
	// makes a hash key for the url with timestamp etc
	h := md5.New()
	io.WriteString(h, timestamp)
	io.WriteString(h, client.PrivateKey)
	io.WriteString(h, client.PublicKey)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (client *Client) buildUrl(endpoint string, arguments map[string]string) string {

	ts := time.Now().Local().Format("20060102150405")
	hash := client.makeHash(ts)

	q := url.Values{}

	q.Set("ts", ts)
	q.Set("hash", hash)
	q.Set("apikey", client.PublicKey)

	for k, v := range arguments {
		q.Set(k, v)
	}

	return fmt.Sprintf("%s%s?%s", baseURL, endpoint, q.Encode())
}

func (client *Client) GetResponse(endpoint string, arguments map[string]string) ([]byte, error) {
	res, err := http.Get(client.buildUrl(endpoint, arguments))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func (client *Client) GetComic(comicId string) (*Comic, error) {
	resp, err := client.GetResponse("comics/"+comicId, nil)
	if err != nil {
		return nil, err
	}
	wrapper := &ComicDataWrapper{}
	err = json.Unmarshal(resp, &wrapper)
	if err != nil {
		return nil, err
	}

	if wrapper.Code == 404 {
		return nil, nil
	}

	return &wrapper.Data.Results[0], nil
}

func (client *Client) GetComics(p *ComicQueryParams) ([]Comic, error) {
	resp, err := client.GetResponse("comics", p.QueryArgs())
	if err != nil {
		return nil, err
	}
	wrapper := &ComicDataWrapper{}
	err = json.Unmarshal(resp, &wrapper)
	if err != nil {
		return nil, err
	}
	return wrapper.Data.Results, nil
}
