package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	inputFile = flag.String("input", "", "Input file path")
)

func main() {
	flag.Parse()

	var (
		nContains int
		nOverlaps int
	)

	f, _ := os.Open(*inputFile)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ranges := strings.Split(scanner.Text(), ",")
		r1, r2 := parseRange(ranges[0]), parseRange(ranges[1])
		if r1.contains(r2) || r2.contains(r1) {
			nContains++
		}
		if r1.overlaps(r2) || r2.overlaps(r1) {
			nOverlaps++
		}
	}

	fmt.Println(nContains, nOverlaps)
}

type assignedRange struct {
	start, end int
}

func (r assignedRange) contains(other assignedRange) bool {
	return r.start >= other.start && r.end <= other.end
}

func (r assignedRange) overlaps(other assignedRange) bool {
	return r.start >= other.start && r.start <= other.end
}

func parseRange(text string) assignedRange {
	startEnd := strings.Split(text, "-")
	start, _ := strconv.Atoi(startEnd[0])
	end, _ := strconv.Atoi(startEnd[1])
	return assignedRange{start: start, end: end}
}
