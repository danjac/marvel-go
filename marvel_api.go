package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL        = "http://gateway.marvel.com:80/v1/public/"
	defaultDateFmt = "2006-01-02T15:04:05-0700"
)

type JsonTime time.Time

func (jt *JsonTime) UnmarshalJSON(data []byte) error {

	b := bytes.NewBuffer(data)
	dec := json.NewDecoder(b)
	var s string
	if err := dec.Decode(&s); err != nil {
		return err
	}
	t, err := time.Parse(defaultDateFmt, s)
	if err != nil {
		return err
	}
	*jt = (JsonTime)(t)
	return nil
}

type TextObject struct {
	Type     string `json "type"`
	Language string `json "language"`
	Text     string `json "text"`
}

type SeriesSummary struct {
	ResourceURI string `json "resourceURI"`
	Name        string `json "name"`
}

type ComicSummary struct {
	ResourceURI string `json "resourceURI"`
	Name        string `json "name"`
}

type ComicDate struct {
	Type string   `json "type"`
	Date JsonTime `json "date"`
}

type ComicPrice struct {
	Type  string  `json "type"`
	Price float64 `json "price"`
}

type CreatorList struct {
	Available     int              `json "available"`
	Returned      int              `json "returned"`
	CollectionURI string           `json "collectionURI"`
	Items         []CreatorSummary `json "items"`
}

type CreatorSummary struct {
	ResourceURI string `json "resourceURI"`
	Name        string `json "name"`
	Role        string `json "role"`
}

type CharacterList struct {
	Available     int                `json "available"`
	Returned      int                `json "returned"`
	CollectionURI string             `json "collectionURI"`
	Items         []CharacterSummary `json "items"`
}

type CharacterSummary struct {
	ResourceURI string `json "resourceURI"`
	Name        string `json "name"`
	Role        string `json "role"`
}

type StoryList struct {
	Available     int            `json "available"`
	Returned      int            `json "returned"`
	CollectionURI string         `json "collectionURI"`
	Items         []StorySummary `json "items"`
}

type StorySummary struct {
	ResourceURI string `json "resourceURI"`
	Name        string `json "name"`
	Type        string `json "type"`
}

type EventList struct {
	Available     int            `json "available"`
	Returned      int            `json "returned"`
	CollectionURI string         `json "collectionURI"`
	Items         []EventSummary `json "items"`
}

type EventSummary struct {
	ResourceURI string `json "resourceURI"`
	Name        string `json "name"`
}

type Image struct {
	Path      string `json "path"`
	Extension string `json "extension"`
}

type URL struct {
	Type string `json "type"`
	URL  string `json "url"`
}

type Comic struct {
	Id                 int            `json "id"`
	DigitalId          int            `json "digitalId"`
	Title              string         `json "title"`
	Format             string         `json "format"`
	IssueNumber        float64        `json "issueNumber"`
	VariantDescription string         `json "variantDescription"`
	Description        string         `json "description"`
	Modified           JsonTime       `json "modified"`
	Isbn               string         `json "isbn"`
	UPC                string         `json "upc"`
	DiamondCode        string         `json "diamondCode"`
	EAN                string         `json "ean"`
	ISSN               string         `json "issn"`
	PageCount          int            `json "pageCount"`
	TextObjects        []TextObject   `json "textObjects"`
	ResourceURI        string         `json "resourceURI"`
	URLs               []URL          `json "urls"`
	Series             SeriesSummary  `json "series"`
	Variants           []ComicSummary `json "variants"`
	Collections        []ComicSummary `json "collections"`
	CollectedIssues    []ComicSummary `json "collectedIssues"`
	Dates              []ComicDate    `json "dates"`
	Prices             []ComicPrice   `json "prices"`
	Thumbnail          Image          `json "thumbnail"`
	Images             []Image        `json "images"`
	Creators           CreatorList    `json "creators"`
	Characters         CharacterList  `json "characters"`
	Stories            StoryList      `json "stories"`
	Events             EventList      `json "events"`
}

type ComicDataContainer struct {
	Offset  int     `json: "offset"`
	Limit   int     `json: "limit"`
	Total   int     `json: "total"`
	Count   int     `json: "count"`
	Results []Comic `json "results"`
}

type ComicDataWrapper struct {
	Code            int                `json: "code"`
	Status          string             `json: "status"`
	Copyright       string             `json: "copyright"`
	AttributionText string             `json: "attributionText"`
	AttributionHTML string             `json: "attributionHTML"`
	Data            ComicDataContainer `json  "data"`
	Etag            string             `json: "string"`
}

func makeIntString(a []int) string {
	l := make([]string, len(a))
	for _, v := range a {
		l = append(l, strconv.Itoa((int)(v)))
	}
	return strings.Join(l, ",")
}

type ComicQueryParams struct {
	Format            string
	FormatType        string
	NoVariants        bool
	DateDescriptor    string
	DateRange         int
	DiamondCode       string
	DigitalId         int
	UPC               string
	ISBN              string
	EAN               string
	ISSN              string
	HasDigitalIssue   bool
	ModifiedSince     time.Time
	Creators          []int
	Characters        []int
	Series            []int
	Events            []int
	Stories           []int
	SharedAppearances []int
	Collaborators     []int
	OrderBy           string
	Limit             int
	Offset            int
}

func (p *ComicQueryParams) ToQueryString() map[string]string {
	args := make(map[string]string)
	if p.Format != "" {
		args["format"] = p.Format
	}
	if p.FormatType != "" {
		args["formatType"] = p.FormatType
	}
	if p.DateDescriptor != "" {
		args["dateDescriptor"] = p.DateDescriptor
	}
	if p.DiamondCode != "" {
		args["diamondCode"] = p.DiamondCode
	}
	if p.ISBN != "" {
		args["isbn"] = p.ISBN
	}
	if p.UPC != "" {
		args["upc"] = p.UPC
	}
	if p.NoVariants {
		args["noVariants"] = "true"
	}
	if p.HasDigitalIssue {
		args["HasDigitalIssue"] = "true"
	}
	if p.DigitalId != 0 {
		args["digitalId"] = strconv.Itoa((int)(p.DigitalId))
	}
	if p.Limit != 0 {
		args["limit"] = strconv.Itoa((int)(p.Limit))
	}
	if p.Offset != 0 {
		args["offset"] = strconv.Itoa((int)(p.Offset))
	}
	if p.DateRange != 0 {
		args["limit"] = strconv.Itoa((int)(p.DateRange))
	}
	if len(p.Creators) > 0 {
		args["creators"] = makeIntString(p.Creators)
	}
	if !p.ModifiedSince.IsZero() {
		args["modifiedSince"] = p.ModifiedSince.Format(defaultDateFmt)
	}

	return args
}

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
	resp, err := client.GetResponse("comics", p.ToQueryString())
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
