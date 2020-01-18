package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type symbol string
type mconfigName string

type operationType int

const (
	P operationType = iota
	E
	R
	L
)

type operation struct {
	t   operationType
	arg symbol
}

func (o *operation) String() string {
	switch o.t {
	case P:
		return fmt.Sprintf("P%s", o.arg)
	case E:
		return "E"
	case R:
		return "R"
	case L:
		return "L"
	}
	return "X"
}

type behavior struct {
	ops          []*operation
	finalMconfig mconfigName
}

func (b *behavior) String() string {
	return fmt.Sprintf("%v -> %s", b.ops, b.finalMconfig)
}

type mconfig struct {
	rules map[symbol]*behavior
}

func (m *mconfig) String() string {
	return fmt.Sprintf("%s", m.rules)
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	file, err := os.Open("foo.txt")
	if err != nil {
		// HACK
		file, err = os.Open("../../foo.txt")
		if err != nil {
			log.Fatal(err)
		}
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	text := string(b)

	lines := strings.Split(text, "\n")

	configuration := make(map[mconfigName]*mconfig)

	for lineNum, line := range lines {
		fmt.Printf("%d: %s\n", lineNum, line)

		items := strings.Split(line, " ")
		if len(items) != 4 {
			log.Fatal("Line should only have four items: " + line)
		}

		rawMconfig := items[0]
		rawSymbol := items[1]
		rawOps := items[2]
		rawFinalMconfig := items[3]

		ops := parseOps(rawOps)
		b := &behavior{
			ops:          ops,
			finalMconfig: mconfigName(rawFinalMconfig),
		}

		m, ok := configuration[mconfigName(rawMconfig)]
		if !ok {
			m = &mconfig{rules: make(map[symbol]*behavior)}
			configuration[mconfigName(rawMconfig)] = m
		}

		_, ok = m.rules[symbol(rawSymbol)]
		if ok {
			log.Fatal("duplicate rule")
		}
		m.rules[symbol(rawSymbol)] = b
	}

	fmt.Printf("Result: %s", configuration)
}

func parseOps(input string) []*operation {
	groups := strings.Split(input, ",")
	out := make([]*operation, len(groups))

	for i, group := range groups {
		switch group[0:1] {
		case "P":
			arg := group[1:]
			out[i] = &operation{t: P, arg: symbol(arg)}
		case "E":
			out[i] = &operation{t: E}
		case "R":
			out[i] = &operation{t: R}
		case "L":
			out[i] = &operation{t: L}
		}
	}

	return out
}
