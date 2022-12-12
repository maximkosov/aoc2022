package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	inputFile = flag.String("input", "", "Input file path")
	noReleaf  = flag.Bool("no-relief", false, "No relief worry level when item passed")
	rounds    = flag.Int("rounds", 20, "Rounds to pass")
)

func main() {
	flag.Parse()

	var monkeys []*monkey

	f, _ := os.Open(*inputFile)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "Monkey") {
			m, err := readMonkey(scanner)
			if err != nil {
				panic(err)
			}

			monkeys = append(monkeys, m)
		}
	}

	divProduct := 1
	for _, m := range monkeys {
		divProduct *= m.divisor
	}

	for i := 0; i < *rounds; i++ {
		for _, m := range monkeys {
			for _, it := range m.items {
				m.inspect(it)

				if *noReleaf {
					it.worryLevel %= divProduct
				} else {
					it.worryLevel /= 3
				}

				m.throwToNext(it, monkeys)
			}

			m.items = nil
		}
	}

	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].inspectedTimes > monkeys[j].inspectedTimes
	})

	monkeyBusiness := 1
	for i := 0; i < 2; i++ {
		fmt.Println(monkeys[i].inspectedTimes)
		monkeyBusiness *= monkeys[i].inspectedTimes
	}

	fmt.Println(monkeyBusiness)
}

type item struct {
	worryLevel int
}

type monkey struct {
	items          []*item
	inspectedTimes int
	divisor        int
	inspect        func(it *item)
	throwToNext    func(it *item, monkeys []*monkey)
}

func readMonkey(scanner *bufio.Scanner) (*monkey, error) {
	var m = &monkey{}

	groups, err := expect(scanner, `  Starting items: (.+)`)
	if err != nil {
		return nil, fmt.Errorf("parse starting items: %w", err)
	}
	for _, strLevel := range strings.Split(groups[1], ", ") {
		it := &item{}
		it.worryLevel, _ = strconv.Atoi(strLevel)
		m.items = append(m.items, it)
	}

	groups, err = expect(scanner, `  Operation: new = old (\+|\*) (old|\d+)`)
	if err != nil {
		return nil, fmt.Errorf("parse operation: %w", err)
	}

	getY := func(it *item) int { return it.worryLevel }
	if groups[2] != "old" {
		y, _ := strconv.Atoi(groups[2])
		getY = func(it *item) int { return y }
	}

	var op func(x, y int) int
	switch groups[1] {
	case "+":
		op = func(x, y int) int { return x + y }
	case "*":
		op = func(x, y int) int { return x * y }
	}

	m.inspect = func(it *item) {
		it.worryLevel = op(it.worryLevel, getY(it))
		m.inspectedTimes++
	}

	groups, err = expect(scanner, `  Test: divisible by (\d+)`)
	if err != nil {
		return nil, fmt.Errorf("parse test: %w", err)
	}
	m.divisor, _ = strconv.Atoi(groups[1])

	groups, err = expect(scanner, `    If true: throw to monkey (\d+)`)
	if err != nil {
		return nil, fmt.Errorf("parse test: %w", err)
	}
	trueMonkey, _ := strconv.Atoi(groups[1])

	groups, err = expect(scanner, `    If false: throw to monkey (\d+)`)
	if err != nil {
		return nil, fmt.Errorf("parse test: %w", err)
	}
	falseMonkey, _ := strconv.Atoi(groups[1])

	m.throwToNext = func(it *item, monkeys []*monkey) {
		if it.worryLevel%m.divisor == 0 {
			monkeys[trueMonkey].items = append(monkeys[trueMonkey].items, it)
		} else {
			monkeys[falseMonkey].items = append(monkeys[falseMonkey].items, it)
		}
	}

	return m, nil
}

func expect(scanner *bufio.Scanner, pattern string) ([]string, error) {
	if !scanner.Scan() {
		return nil, fmt.Errorf("can't scan: %w", scanner.Err())
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("compile regexp: %w", err)
	}
	return re.FindStringSubmatch(scanner.Text()), nil
}
