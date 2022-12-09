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
	length    = flag.Int("length", 2, "Rope length")
)

func main() {
	flag.Parse()

	input, err := os.Open(*inputFile)
	if err != nil {
		panic(err)
	}
	defer input.Close()

	var (
		rope    = make([]coord, *length)
		visited = make(map[coord]int)
	)

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		n, _ := strconv.Atoi(line[1])
		for i := 0; i < n; i++ {
			head := rope[0]
			switch line[0] {
			case "L":
				head.x -= 1
			case "R":
				head.x += 1
			case "U":
				head.y += 1
			case "D":
				head.y -= 1
			}
			rope[0] = head
			for i := 1; i < len(rope); i++ {
				rope[i] = updateTail(rope[i-1], rope[i])
			}
			visited[rope[len(rope)-1]] += 1
		}
	}

	fmt.Println(len(visited))
}

type coord struct {
	x, y int
}

func updateTail(head, tail coord) coord {
	diffX, diffY := head.x-tail.x, head.y-tail.y
	if abs(diffX) >= 2 && abs(diffY) >= 2 {
		tail.x += (diffX - sign(diffX))
		tail.y += (diffY - sign(diffY))
		return tail
	}
	if abs(diffX) >= 2 {
		tail.x += (diffX - sign(diffX))
		tail.y = head.y
		return tail
	}
	if abs(diffY) >= 2 {
		tail.x = head.x
		tail.y += (diffY - sign(diffY))
	}
	return tail
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x < 0 {
		return -1
	}
	return 1
}
