package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
)

var (
	inputFile = flag.String("input", "", "Input file path")
	nElves    = flag.Int("n", 1, "Number of elves to sum")
)

func main() {
	flag.Parse()

	input, err := os.Open(*inputFile)
	if err != nil {
		panic(err)
	}
	defer input.Close()

	elves := []*elf{{}}
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		if scanner.Text() == "" {
			elves = append(elves, &elf{number: len(elves), calories: 0})
			continue
		}
		cal, _ := strconv.Atoi(scanner.Text())
		elves[len(elves)-1].calories += cal
	}

	sort.Slice(elves, func(i, j int) bool {
		return elves[i].calories > elves[j].calories
	})

	var sum int
	for i := 0; i < *nElves; i++ {
		sum += elves[i].calories
	}
	fmt.Println(sum)
}

type elf struct {
	number   int
	calories int
}
