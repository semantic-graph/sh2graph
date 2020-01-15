package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"mvdan.cc/sh/syntax"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Pair struct {
	a, b interface{}
}

type Graph struct {
	Nodes [][]string
	Edges [][]int
}

func main() {
	inputPath := os.Args[1]
	dat, err := ioutil.ReadFile(inputPath)
	check(err)
	r := strings.NewReader(string(dat))
	f, err := syntax.NewParser().Parse(r, "")
	if err != nil {
		return
	}
	var g Graph
	syntax.Walk(f, func(node syntax.Node) bool {
		switch x := node.(type) {
		case *syntax.CallExpr:
			if len(x.Args) == 0 || len(x.Args[0].Parts) == 0 {
				fmt.Println("\nSkip CallExpr: ")
				syntax.DebugPrint(os.Stdout, x)
				return true
			}
			prog := x.Args[0].Parts[0]
			switch progNode := prog.(type) {
			case *syntax.Lit:
				{
					progName := progNode.Value
					operations := []string{progName}
					for _, arg := range x.Args[1:] {
						switch x := arg.Parts[0].(type) {
						case *syntax.Lit:
							operations = append(operations, x.Value)
						case *syntax.SglQuoted:
							operations = append(operations, x.Value)
						}
					}
					g.Nodes = append(g.Nodes, operations)
				}
			}
		}
		return true
	})
	// syntax.DebugPrint(os.Stdout, f)

	file, err := json.MarshalIndent(g, "", " ")
	check(err)
	err = ioutil.WriteFile(inputPath+".sh2graph.json", file, 0644)
	check(err)
}
