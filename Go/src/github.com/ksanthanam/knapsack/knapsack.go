package knapsack

import (
	"fmt"
	"sort"
)

var (
	p  = fmt.Println
	pf = fmt.Printf
	s  = fmt.Sprintf
)

// MaxPosition struct
type MaxPosition struct {
	row, col, profit int
}

// Object with Weight
type Object struct {
	Weight, Profit int
}

// Equal for Object
func (o Object) Equal(c Object) bool {
	return o.Weight == c.Weight && o.Profit == c.Profit
}

// Objects a Collection of Objects
type Objects []Object

// Len function for Sorting
func (os Objects) Len() int { return len(os) }

// Swap function for Sorting
func (os Objects) Swap(i, j int) { os[i], os[j] = os[j], os[i] }

// Less function for Sorting
func (os Objects) Less(i, j int) bool { return os[i].Weight < os[j].Weight }

// Equal function
func (os Objects) Equal(cs Objects) bool {
	if os.Len() != cs.Len() {
		p("Different Length")
		return false
	}
	left := []Object(os)
	right := []Object(cs)
	for i := range left {
		if !left[i].Equal(right[i]) {
			return false
		}
	}
	return true
}

/*
KnapSack 0/1 Problem Example
m = 8
n = 4
P = {1,2,5,6}
W = {2,3,4,5}

Pi Wi         0   1   2   3   4   5   6   7   8
            |^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
          0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
 1   2    1 | 0 | 0 | 1 | 1 | 1 | 1 | 1 | 1 | 1 |
 2   3    2 | 0 | 0 | 1 | 2 | 2 | 3 | 3 | 3 | 3 |
 5   4    3 | 0 | 0 | 1 | 2 | 5 | 5 | 6 | 7 | 7 |
 6   5    4 | 0 | 0 | 1 | 2 | 5 | 6 | 6 | 7 | 8 |

*/
// GetKnapSack solution
func GetKnapSack(objects Objects, m int) Objects {
	os := objects
	sort.Sort(os)
	ordered := []Object{Object{0, 0}}
	ordered = append(ordered, []Object(os)...)
	n := len(ordered)
	V := make([][]int, n)
	max := MaxPosition{0, 0, 0}
	placeValue := func(r, w int) int {
		o := ordered[r]
		preProfit := func() int {
			if w-o.Weight >= 0 {
				return V[r-1][w-o.Weight]
			}
			return 0
		}
		pP := V[r-1][w]
		if w >= o.Weight {
			cP := o.Profit + preProfit()
			if pP > cP {
				return pP
			}
			return cP
		}
		return pP
	}

	for i := 0; i < n; i++ {
		V[i] = make([]int, m+1)
		for w := 0; w <= m; w++ {
			if i == 0 || w == 0 {
				V[i][w] = 0
			} else {
				V[i][w] = placeValue(i, w)
			}
			if V[i][w] > max.profit {
				max = MaxPosition{i, w, V[i][w]}
			}
		}
	}

	selected := make([]Object, 0)
	profitLeft := max.profit
	maxi := max.row
	maxw := max.col
	for profitLeft > 0 && maxi > 0 && maxw > 0 {
		if V[maxi][maxw] > V[maxi-1][maxw] {
			selected = append(selected, ordered[maxi])
			profitLeft -= ordered[maxi].Profit
			maxw -= ordered[maxi].Weight
		}
		maxi--
	}
	return selected
}
