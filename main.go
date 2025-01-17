package main

import (
	"fmt"

	"github.com/sbchaos/consume/comb"
	"github.com/sbchaos/consume/par"
	"github.com/sbchaos/consume/par/char"
	sp "github.com/sbchaos/consume/par/strings"
	ss "github.com/sbchaos/consume/stream/strings"
)

func main() {
	s1 := ss.NewStringStream("item0 `item1` \"item2\" 'item3' item4")
	p := sp.StringWithOptionalQuotes()

	lst, err := par.Parse(s1,
		comb.Count(5, comb.SkipAfter(p, char.WhiteSpaces())),
	)
	if err != nil {
		fmt.Printf("Error in parsing: %s", err)
	}
	fmt.Println(lst)
}
