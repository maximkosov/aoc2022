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
	fsSize    = flag.Int("fs-size", 0, "Total filesystem size")
	reqSize   = flag.Int("req-size", 0, "Required size for update")
)

func main() {
	flag.Parse()

	input, err := os.Open(*inputFile)
	if err != nil {
		panic(err)
	}
	defer input.Close()

	var (
		root    = newDirectory(nil, "/")
		current *node
	)

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		switch token := parseToken(scanner.Text()).(type) {
		case tokenCommadCD:
			switch token.target {
			case "/":
				current = root
			case "..":
				current = current.parent
			default:
				current = current.findChild(token.target)
			}

		case tokenCommandLS:
			continue

		case tokenFile:
			current.addChild(newFile(current, token.name, token.size))

		case tokenDirectory:
			current.addChild(newDirectory(current, token.name))
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	root.updateSize()

	dirs := root.where(func(n *node) bool {
		return n.nodeType == directory
	})
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].size < dirs[j].size
	})

	freeSize := *fsSize - root.size
	for _, n := range dirs {
		if freeSize+n.size >= *reqSize {
			fmt.Printf("Total size is %d. I need to delete %s to free up %d and have total %d\n", root.size, n.name, n.size, freeSize+n.size)
			return
		}
	}
}

type tokenCommadCD struct {
	target string
}

type tokenCommandLS struct {
}

type tokenDirectory struct {
	name string
}

type tokenFile struct {
	name string
	size int
}

var fileRegex = regexp.MustCompile(`(\d+) (.+)`)

func parseToken(line string) any {
	switch {
	case strings.HasPrefix(line, "$ cd "):
		return tokenCommadCD{
			target: strings.TrimPrefix(line, "$ cd "),
		}

	case strings.HasPrefix(line, "$ ls"):
		return tokenCommandLS{}

	case strings.HasPrefix(line, "dir "):
		return tokenDirectory{
			name: strings.TrimPrefix(line, "dir "),
		}

	default:
		groups := fileRegex.FindStringSubmatch(line)
		if len(groups) != 3 {
			return nil
		}

		size, _ := strconv.Atoi(groups[1])
		return tokenFile{
			name: groups[2],
			size: size,
		}
	}
}

type nodeType int

const (
	directory nodeType = iota
	file
)

type node struct {
	name     string
	nodeType nodeType
	size     int
	parent   *node
	children map[string]*node
}

func newFile(parent *node, name string, size int) *node {
	return &node{
		name:     name,
		size:     size,
		nodeType: file,
		parent:   parent,
	}
}

func newDirectory(parent *node, name string) *node {
	return &node{
		name:     name,
		nodeType: directory,
		parent:   parent,
		children: make(map[string]*node),
	}
}

func (n *node) addChild(child *node) {
	n.children[child.name] = child
}

func (n *node) findChild(name string) *node {
	return n.children[name]
}

func (n *node) traverse(cb func(*node)) {
	if n.children != nil {
		for _, child := range n.children {
			child.traverse(cb)
		}
	}
	cb(n)
}

func (n *node) updateSize() {
	n.traverse(func(descendant *node) {
		if descendant.children != nil {
			for _, child := range descendant.children {
				descendant.size += child.size
			}
		}
	})
}

func (n *node) where(predicate func(*node) bool) []*node {
	var res []*node
	n.traverse(func(descendant *node) {
		if predicate(descendant) {
			res = append(res, descendant)
		}
	})
	return res
}
