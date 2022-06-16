package main

import "time"

type Response struct {
	Items          []Item `json:"items"`
	HasMore        bool   `json:"has_more"`
	QuotaMax       int    `json:"quota_max"`
	QuotaRemaining int    `json:"quota_remaining"`
}

type Item map[string]interface{}

func (i Item) Type() string {
	if i.IsRecognizedMember() {
		return "recognized"
	}
	if i.IsRecommendation() {
		return "recommendation"
	}
	return ""
}

func (i Item) IsRecognizedMember() bool {
	_, ok := i["posted_by_collectives"]
	return ok
}

func (i Item) IsRecommendation() bool {
	_, ok := i["recommendations"]
	return ok
}

func (i Item) RecommendationDate() int64 {
	recs := i["recommendations"].([]interface{})
	n := recs[0].(map[string]interface{})["creation_date"].(float64)
	return int64(n)
}

func (i Item) RecommendationDateFmt() string {
	if !i.IsRecommendation() {
		return ""
	}
	d := i.RecommendationDate()
	return time.Unix(d, 0).Format("2006-01-02")
}

func (i Item) RecommendationTime() time.Time {
	return time.Unix(i.RecommendationDate(), 0)
}

func (i Item) Owner() *Owner {
	owner := i["owner"].(map[string]interface{})
	return &Owner{
		UserID:      maybeInt64(owner["user_id"]),
		Reputation:  maybeInt64(owner["reputation"]),
		UserType:    owner["user_type"].(string),
		DisplayName: maybeString(owner["display_name"]),
		Link:        maybeString(owner["link"]),
	}
}

func (i Item) Answer() *Answer {
	return &Answer{
		IsAccepted:   i["is_accepted"].(bool),
		Score:        int64(i["score"].(float64)),
		CreationDate: int64(i["creation_date"].(float64)),
		AnswerID:     int64(i["answer_id"].(float64)),
		QuestionID:   int64(i["question_id"].(float64)),
	}
}

type Answer struct {
	IsAccepted   bool  `json:"is_accepted"`
	Score        int64 `json:"score"`
	CreationDate int64 `json:"creation_date"`
	AnswerID     int64 `json:"answer_id"`
	QuestionID   int64 `json:"question_id"`
}

func (a Answer) DateFmt() string {
	return time.Unix(a.CreationDate, 0).Format("2006-01-02")
}

type Collective struct {
	Tags          []string `json:"tags"`
	ExternalLinks []struct {
		Type string `json:"type"`
		Link string `json:"link"`
	} `json:"external_links"`
	Slug string `json:"slug"`
}

type Recommendation struct {
	Collective   *Collective `json:"collective"`
	CreationDate int64       `json:"creation_date"`
}

type Owner struct {
	UserID      int64  `json:"user_id"`
	Reputation  int64  `json:"reputation"`
	UserType    string `json:"user_type"`
	DisplayName string `json:"display_name"`
	Link        string `json:"link"`
}

func maybeInt64(v interface{}) int64 {
	n, ok := v.(float64)
	if !ok {
		return -1
	}
	return int64(n)
}

func maybeString(v interface{}) string {
	s, _ := v.(string)
	return s
}
