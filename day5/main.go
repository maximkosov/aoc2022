package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var (
	inputFile    = flag.String("input", "", "Input file")
	useCrane9001 = flag.Bool("crane9001", false, "Use fancy crane 9001")
)

func exit(msg string, err error) {
	fmt.Printf("%s: %s\n", msg, err)
	os.Exit(1)
}

func main() {
	flag.Parse()

	input, err := os.Open(*inputFile)
	if err != nil {
		exit("Err open file", err)
	}
	defer input.Close()

	scanner := bufio.NewScanner(input)

	stacks := readStacks(scanner)

	mover := move9000
	if *useCrane9001 {
		mover = move9001
	}

	err = processMoves(scanner, func(from, to, count int) {
		mover(stacks, from, to, count)
	})
	if err != nil {
		exit("Err processing moves", err)
	}

	if err := scanner.Err(); err != nil {
		exit("Scan error", err)
	}

	fmt.Println(topCrates(stacks))
}

type crate string

type stack []crate

func readStacks(scanner *bufio.Scanner) []stack {
	var (
		stacks []stack
	)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		var (
			stack  int = -1
			spaces int
		)
		for _, b := range line {
			switch {
			case b == ' ':
				spaces++
				continue
			case b == '[':
				stack += (1 + spaces/4)
				spaces = 0
				continue
			case b == ']' || b < 'A' || b > 'Z':
				continue
			default:
				for stack > len(stacks)-1 {
					stacks = append(stacks, nil)
				}
				stacks[stack] = append(stacks[stack], crate(b))
			}
		}
	}

	for i := range stacks {
		reverse(stacks[i])
	}

	return stacks
}

func (s stack) pop(n int) ([]crate, stack) {
	return s[len(s)-n:], s[:len(s)-n]
}

func reverse[T any](arr []T) {
	for i := 0; i < len(arr)/2; i++ {
		arr[i], arr[len(arr)-i-1] = arr[len(arr)-i-1], arr[i]
	}
}

var moveRe = regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)

func processMoves(scanner *bufio.Scanner, handler func(from, to, count int)) error {
	for scanner.Scan() {
		groups := moveRe.FindStringSubmatch(scanner.Text())
		if len(groups) != 4 {
			return fmt.Errorf("unexpected line %q", scanner.Text())
		}

		count, err := strconv.Atoi(groups[1])
		if err != nil {
			return fmt.Errorf("error parsing count: %w", err)
		}
		from, err := strconv.Atoi(groups[2])
		if err != nil {
			return fmt.Errorf("error parsing from: %w", err)
		}
		to, err := strconv.Atoi(groups[3])
		if err != nil {
			return fmt.Errorf("error parsing to: %w", err)
		}

		handler(from, to, count)
	}
	return nil
}

func move9000(stacks []stack, from, to, count int) {
	var toMove []crate
	toMove, stacks[from-1] = stacks[from-1].pop(count)
	reverse(toMove)
	stacks[to-1] = append(stacks[to-1], toMove...)
}

func move9001(stacks []stack, from, to, count int) {
	var toMove []crate
	toMove, stacks[from-1] = stacks[from-1].pop(count)
	stacks[to-1] = append(stacks[to-1], toMove...)
}

func topCrates(stacks []stack) []crate {
	res := make([]crate, len(stacks))
	for i, s := range stacks {
		if len(s) > 0 {
			res[i] = s[len(s)-1]
		}
	}
	return res
}
