package main

import (
	"fmt"

	h "github.com/sbchaos/consume/base/higher"
	"github.com/sbchaos/consume/char"
	"github.com/sbchaos/consume/run"
	ss "github.com/sbchaos/consume/stream/strings"
	sp "github.com/sbchaos/consume/strings"
)

func main() {
	s1 := ss.NewStringStream("item0 `item1` \"item2\" 'item3' item4")
	p := sp.StringWithOptionalQuotes()

	lst, err := run.Parse(s1,
		h.Count(5, h.SkipAfter(p, char.WhiteSpaces())),
	)
	if err != nil {
		fmt.Printf("Error in parsing: %s", err)
	}
	fmt.Println(lst)
}
