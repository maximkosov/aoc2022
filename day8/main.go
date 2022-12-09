package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var inputFile = flag.String("input", "", "Input file path")

func main() {
	flag.Parse()

	input, err := os.Open(*inputFile)
	if err != nil {
		panic(err)
	}
	defer input.Close()

	var forest forest

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		var row []*tree
		for _, b := range scanner.Text() {
			size, _ := strconv.Atoi(string(b))
			row = append(row, &tree{size: size})
		}
		forest = append(forest, row)
	}

	updateVisibility(forest)

	var (
		visible  int
		maxScore int
	)
	for row := range forest {
		for col := range forest[row] {
			tree := forest[row][col]
			if tree.score() > maxScore {
				maxScore = tree.score()
			}
			if isTreeVisible(forest, tree, row, col) {
				visible++
			}
		}
	}

	fmt.Println(visible, maxScore)
}

type tree struct {
	size                                 int
	visLeft, visRight, visTop, visBottom int
}

func (t *tree) score() int {
	return t.visLeft * t.visRight * t.visTop * t.visBottom
}

func isTreeVisible(forest forest, tree *tree, row, col int) bool {
	// Tree tree is visibile if it's at the border or the last visibile is smaller than tree.
	return col == 0 ||
		tree.size > forest[row][col-tree.visLeft].size ||
		col == len(forest)-1 ||
		tree.size > forest[row][col+tree.visRight].size ||
		row == 0 ||
		tree.size > forest[row-tree.visTop][col].size ||
		row == len(forest[col])-1 ||
		tree.size > forest[row+tree.visBottom][col].size

}

type forest [][]*tree

func updateVisibility(forest forest) {
	for row := range forest {
		for col := range forest[row] {
			tree := forest[row][col]
			for left := col - 1; left >= 0; left-- {
				tree.visLeft++
				if forest[row][left].size >= tree.size {
					break
				}
			}
			for top := row - 1; top >= 0; top-- {
				tree.visTop++
				if forest[top][col].size >= tree.size {
					break
				}
			}
		}
	}
	for row := len(forest) - 1; row >= 0; row-- {
		for col := len(forest[row]) - 1; col >= 0; col-- {
			tree := forest[row][col]
			for right := col + 1; right < len(forest[row]); right++ {
				tree.visRight++
				if forest[row][right].size >= tree.size {
					break
				}
			}
			for bottom := row + 1; bottom < len(forest); bottom++ {
				tree.visBottom++
				if forest[bottom][col].size >= tree.size {
					break
				}
			}
		}
	}
}
