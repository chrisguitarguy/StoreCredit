// https://code.google.com/codejam/contest/351101/dashboard
package main

import (
	"io"
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
	"sort"
	"errors"
)

// A single item in the store. Contains the value and the original position
// Makes it easy to sort and not have to worry about dealing with keeping
// the original position around.
type Item struct {
	Value, Position int
}

type ItemSlice []*Item

func (items ItemSlice) Len() int {
	return len(items)
}

func (items ItemSlice) Less(i, j int) bool {
	return items[i].Value < items[j].Value
}

func (items ItemSlice) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func (items ItemSlice) Search(value int) int {
	return sort.Search(items.Len(), func (i int) bool { return items[i].Value >= value })
}

// A "test case" for the store. Contains the credit, the number of items
// in the store, and the list of items.
type Store struct {
	Credit int
	NumItems int
	Items ItemSlice
}

func toInt(s string) (i int, err error ) {
	i, err = strconv.Atoi(strings.Trim(s, " \n\r\t"))
	return
}

func parseStore(reader *bufio.Reader) (store Store, err error) {
	var cost int

	cred, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	store.Credit, err = toInt(cred)
	if err != nil {
		return
	}

	numItems, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	store.NumItems, err = toInt(numItems)
	if err != nil {
		return
	}

	store.Items = make(ItemSlice, store.NumItems)
	itemLine, err := reader.ReadString('\n')
	itemsStrs := strings.Split(itemLine, " ")
	for i, c := range itemsStrs {
		cost, err = toInt(c)
		if err != nil {
			break
		}
		store.Items[i] = &Item{Value: cost, Position: i}
	}

	return
}

func parseInput(input io.Reader) []Store {
	reader := bufio.NewReader(input)

	// the number of test cases
	caseStr, err := reader.ReadString('\n')
	if err != nil {
		panic("Could not read number of cases")
	}

	numCases, err := toInt(caseStr)
	if err != nil {
		panic("Could convert number of cases to an int")
	}

	stores := make([]Store, numCases)
	for i := 0; i < numCases; i++ {
		s, err := parseStore(reader)
		if err != nil {
			panic(err)
		}
		stores[i] = s
	}

	return stores
}

func findSolution(store Store) (pos1, pos2 int, err error) {
	var diff, item1, item2 int
	var item *Item

	sort.Sort(store.Items)

	item2 = -1
	for item1, item = range store.Items {
		diff = store.Credit - item.Value
		item2 = store.Items.Search(diff)
		if item2 < store.Items.Len() && item2 != item1 && store.Items[item2].Value == diff {
			break // got it
		} else {
			item2 = -1
		}
	}

	if item2 < 0 {
		err = errors.New("Could not find solution")
		return
	}

	pos1, pos2 = store.Items[item1].Position, store.Items[item2].Position
	if pos1 > pos2 {
		pos1, pos2 = pos2, pos1
	}

	return
}

func main() {
	var pos1, pos2 int
	var err error

	stores := parseInput(os.Stdin)

	for i, store := range stores {
		pos1, pos2, err = findSolution(store)
		if err != nil {
			fmt.Printf("Case #%d: no result", i+1)
		} else {
			fmt.Printf("Case #%d: %d %d\n", i+1, pos1+1, pos2+1)
		}
	}
}
