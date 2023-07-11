package main

import "strings"

type EMatchMode uint32

const (
	MatchAll    EMatchMode = iota // match all query words
	MatchAny                      // match any query word
	MatchPhrase                   // match this exact phrase
)

type Search struct {
	Index      string
	MatchMode  EMatchMode
	Query      string
	LogMessage string
}

type SearchResult struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Hits     struct {
		Total         int    `json:"total"`
		TotalRelation string `json:"total_relation"`
		Hits          []struct {
			Id     string `json:"_id"`
			Score  int    `json:"_score"`
			Source struct {
				Name string `json:"name"`
				Date string `json:"date"`
			} `json:"_source"`
			Highlight struct {
				Name []string `json:"name"`
				Date []string `json:"date"`
			} `json:"highlight"`
		} `json:"hits"`
	} `json:"hits"`
}

// NewSearch construct default search which then may be customized. You may just customize 'Query' and m.b. 'Indexes'
// from default one, and it will work like a simple 'Query()' call.
func NewSearch(query string, index string) Search {
	return Search{
		Index:     index,
		Query:     strings.TrimSpace(query),
		MatchMode: MatchAll,
	}
}
