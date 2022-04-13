package stoic

import "time"

// Thought represent user's comment regarding quote with given QuoteId
// and registered at given Time
type Thought struct {
	Time time.Time
	Text string
	QuoteId int
}
