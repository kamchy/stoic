package model

import "time"

// Thought represent user's comment regarding quote with given QuoteId
// and registered at given Time
type Thought struct {
	Time    time.Time
	Text    string
	QuoteId int64
	Id      int64
}

// ThoughtWithQuote is a thought with a quote on which it is based
type ThoughtWithQuote struct {
	Thought
	Quote
}

// Quote represents single quote and contains text and author
type Quote struct {
	Text   string
	Author string
	Id     int64
}

//QuoteWithCount is a quote with the number of related thoughts
type QuoteWithCount struct {
	Quote        Quote
	ThoughtCount int64
}
