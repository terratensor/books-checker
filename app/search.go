package main

type EMatchMode uint32

const (
	MatchAll    EMatchMode = iota // match all query words
	MatchAny                      // match any query word
	MatchPhrase                   // match this exact phrase
)

type Search struct {
	Index     string
	MatchMode EMatchMode
	Query     string
}

// NewSearch construct default search which then may be customized. You may just customize 'Query' and m.b. 'Indexes'
// from default one, and it will work like a simple 'Query()' call.
func NewSearch(query string, index string) Search {
	return Search{
		Index:     index,
		Query:     query,
		MatchMode: MatchAll,
	}
}
