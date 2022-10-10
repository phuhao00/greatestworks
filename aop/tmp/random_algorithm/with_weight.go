package random_algorithm

import (
	"errors"
	"math/rand"
	"sort"
)

type Choice struct {
	Item   interface{}
	Weight uint
}

func NewChoice(item interface{}, weight uint) Choice {
	return Choice{Item: item, Weight: weight}
}

type Chooser struct {
	data   []Choice
	totals []int
	max    int
}

func NewChooser(choices ...Choice) (*Chooser, error) {
	sort.Slice(choices, func(i, j int) bool {
		return choices[i].Weight < choices[j].Weight
	})

	totals := make([]int, len(choices))
	runningTotal := 0
	for i, c := range choices {
		weight := int(c.Weight)
		if (maxInt - runningTotal) <= weight {
			return nil, errWeightOverflow
		}
		runningTotal += weight
		totals[i] = runningTotal
	}

	if runningTotal < 1 {
		return nil, errNoValidChoices
	}

	return &Chooser{data: choices, totals: totals, max: runningTotal}, nil
}

const (
	intSize = 32 << (^uint(0) >> 63) // cf. strconv.IntSize
	maxInt  = 1<<(intSize-1) - 1
)

var (
	errWeightOverflow = errors.New("sum of Choice Weights exceeds max int")
	errNoValidChoices = errors.New("zero Choices with Weight >= 1")
)

func (c Chooser) Pick() interface{} {
	r := rand.Intn(c.max) + 1
	i := searchInts(c.totals, r)
	return c.data[i].Item
}

func (c Chooser) PickSource(rs *rand.Rand) interface{} {
	r := rs.Intn(c.max) + 1
	i := searchInts(c.totals, r)
	return c.data[i].Item
}

func searchInts(a []int, x int) int {
	i, j := 0, len(a)
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		if a[h] < x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
