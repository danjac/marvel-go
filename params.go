package main

import (
	"strconv"
	"time"
)

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
