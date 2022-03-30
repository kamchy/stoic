package main

import (
	"fmt"
	"os"

	"github.com/kamchy/stoic"
)

func main() {
	if quotes, err := stoic.ReadQuotes(os.Stdin); err == nil {
		for i, quote := range quotes {
			fmt.Printf("Quote %d\n%s\n%s\n\n", i, quote.Text, quote.Author)
		}
	}
}
