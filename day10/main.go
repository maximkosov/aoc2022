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
		cpu         = cpu{x: 1}
		sumStrength int
	)

	f, _ := os.Open(*inputFile)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		var cmd command
		switch {
		case line == "noop":
			cmd = noop()
		case strings.HasPrefix(line, "addx "):
			v, _ := strconv.Atoi(strings.TrimPrefix(line, "addx "))
			cmd = addx(v)
		}

		for i := 0; i < cmd.cycles; i++ {
			cpu.cycle += 1

			if (cpu.cycle-20)%40 == 0 {
				sumStrength += cpu.signalStrength()
			}

			pos := cpu.cycle % 40
			if pos >= cpu.x && pos < cpu.x+3 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
			if pos == 0 {
				fmt.Println()
			}
		}
		cpu.x = cmd.action(cpu.x)
	}

	fmt.Println()
	fmt.Println(sumStrength)
}

type command struct {
	cycles int
	action func(x int) int
}

func noop() command {
	return command{
		cycles: 1,
		action: func(x int) int { return x },
	}
}

func addx(v int) command {
	return command{
		cycles: 2,
		action: func(x int) int { return x + v },
	}
}

type cpu struct {
	cycle int
	x     int
}

func (c cpu) signalStrength() int {
	return c.cycle * c.x
}
