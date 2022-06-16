package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const host = "https://api.stackexchange.com"
const path = "/2.3/collectives/go/answers"
const pagesize = "50"

func main() {
	page := 1
	obj := fetch(page)

	var data []*Data
	data = aggregateData(obj.Items)

	for obj.HasMore {
		page++
		obj = fetch(page)
		data = append(data, aggregateData(obj.Items)...)
	}

	write(data)
}

func fetch(page int) *Response {
	q := url.Values{}
	q.Set("order", "desc")
	q.Set("sort", "activity")
	q.Set("site", "stackoverflow")
	q.Set("pagesize", pagesize)
	q.Set("page", strconv.Itoa(page))

	resp, err := http.DefaultClient.Get(host + path + "?" + q.Encode())
	if err != nil {
		panic(err)
	}
	obj := &Response{}
	err = json.NewDecoder(resp.Body).Decode(obj)
	if err != nil {
		panic(err)
	}
	return obj
}

type Data struct {
	UserID      int64
	UserName    string
	Type        string // recognized,recommendation
	Posted      string // YYYY-MM-DD
	Recommended string // YYYY-MM-DD
	QLink       string
}

func (d Data) AsCSV() string {
	return strings.Join([]string{
		strconv.FormatInt(d.UserID, 10),
		d.UserName,
		d.Type,
		d.Posted,
		d.Recommended,
		d.QLink,
	}, ",")
}

func aggregateData(items []Item) (data []*Data) {
	for _, item := range items {
		owner := item.Owner()
		d := &Data{
			UserID:      owner.UserID,
			UserName:    owner.DisplayName,
			Type:        item.Type(),
			Posted:      item.Answer().DateFmt(),
			Recommended: item.RecommendationDateFmt(),
			QLink:       "https://stackoverflow.com/questions/" + strconv.FormatInt(item.Answer().QuestionID, 10),
		}
		data = append(data, d)
	}
	return
}

func write(data []*Data) {
	f, err := os.OpenFile("file.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// dump
		fmt.Println("id,name,type,posted,recommended,qlink")
		for _, d := range data {
			fmt.Println(d.AsCSV())
		}
		return
	}
	_, err = fmt.Fprintln(f, "id,name,type,posted,recommended,qlink")
	if err != nil {
		panic(err)
	}
	for _, d := range data {
		_, err = fmt.Fprintln(f, d.AsCSV())
		if err != nil {
			panic(err)
		}
	}
}
