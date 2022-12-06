package main

import (
	"flag"
	"fmt"
)

var (
	lookupChars = flag.Int("n", 4, "Number of unmatching characters")
	input       = flag.String("input", "", "Input string")
)

func main() {
	flag.Parse()

	proc := &processor{lookupChars: *lookupChars}
	for i, ch := range *input {
		if proc.process(ch) {
			fmt.Println(i + 1)
			break
		}
	}
}

type processor struct {
	lookupChars int
	seq         []rune
}

func (p *processor) process(ch rune) bool {
	for i := len(p.seq) - 1; i >= 0; i-- {
		if p.seq[i] == ch {
			p.seq = p.seq[i+1:]
			break
		}
	}
	p.seq = append(p.seq, ch)
	return len(p.seq) == p.lookupChars
}
