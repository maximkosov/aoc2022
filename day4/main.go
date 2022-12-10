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
		r1, r2, err := parseRanges(scanner.Text())
		if err != nil {
			panic(err)
		}

		if r1.contains(r2) || r2.contains(r1) {
			nContains++
		}
		if r1.overlaps(r2) {
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
	return r.start >= other.start && r.start <= other.end ||
		other.start >= r.start && other.start <= r.end
}

func parseRanges(text string) (assignedRange, assignedRange, error) {
	ranges := strings.Split(text, ",")
	if len(ranges) != 2 {
		return assignedRange{}, assignedRange{}, fmt.Errorf("cannot parse ranges from %s", text)
	}

	r1, err := parseRange(ranges[0])
	if err != nil {
		return assignedRange{}, assignedRange{}, err
	}

	r2, err := parseRange(ranges[1])
	if err != nil {
		return assignedRange{}, assignedRange{}, err
	}

	return r1, r2, nil
}

func parseRange(text string) (assignedRange, error) {
	startEnd := strings.Split(text, "-")
	if len(startEnd) != 2 {
		return assignedRange{}, fmt.Errorf("cannot parse range from %s", startEnd)
	}
	start, err := strconv.Atoi(startEnd[0])
	if err != nil {
		return assignedRange{}, fmt.Errorf("cannot parse start: %w", err)
	}
	end, err := strconv.Atoi(startEnd[1])
	if err != nil {
		return assignedRange{}, fmt.Errorf("cannot parse end: %w", err)
	}
	return assignedRange{start: start, end: end}, nil
}
