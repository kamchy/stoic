package model

import "time"

// Thought represent user's comment regarding quote with given QuoteId
// and registered at given Time
type Thought struct {
	Time time.Time
	Text string
	QuoteId int
}

// Quote represents single quote and contains text and author
type Quote struct {
	Text   string
	Author string
	Id int
}