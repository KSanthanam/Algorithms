// Package nqueens for Matrix related problems
package nqueens

import (
	"testing"
)

// TestGetQueenPositions
func TestGetQueenPositions(t *testing.T) {
	cases := []struct {
		size      int
		positions Positions
		solutions []Positions
	}{
		{4,
			Positions([]Cell{Cell{0, 1}, Cell{1, 3}, Cell{2, 0}, Cell{3, 2}}),
			[]Positions{
				Positions([]Cell{Cell{0, 1}, Cell{1, 3}, Cell{2, 0}, Cell{3, 2}}),
				Positions([]Cell{Cell{0, 1}, Cell{1, 3}, Cell{2, 0}, Cell{3, 2}}),
				Positions([]Cell{Cell{0, 1}, Cell{1, 3}, Cell{2, 0}, Cell{3, 2}}),
				Positions([]Cell{Cell{0, 2}, Cell{1, 0}, Cell{2, 3}, Cell{3, 1}})},
		},

		//[(0,0) (1,2) (2,5) (3,7) (4,9) (5,4) (6,8) (7,1) (8,3) (9,6)]
		{
			10,
			Positions([]Cell{Cell{0, 0}, Cell{1, 2}, Cell{2, 5}, Cell{3, 7}, Cell{4, 9}, Cell{5, 4}, Cell{6, 8}, Cell{7, 1}, Cell{8, 3}, Cell{9, 6}}),
			[]Positions{
				Positions([]Cell{Cell{0, 0}, Cell{1, 2}, Cell{2, 5}, Cell{3, 7}, Cell{4, 9}, Cell{5, 4}, Cell{6, 8}, Cell{7, 1}, Cell{8, 3}, Cell{9, 6}}),
				Positions([]Cell{Cell{0, 1}, Cell{1, 3}, Cell{2, 5}, Cell{3, 7}, Cell{4, 9}, Cell{5, 0}, Cell{6, 2}, Cell{7, 4}, Cell{8, 6}, Cell{9, 8}}),
				Positions([]Cell{Cell{0, 2}, Cell{1, 0}, Cell{2, 5}, Cell{3, 8}, Cell{4, 4}, Cell{5, 9}, Cell{6, 7}, Cell{7, 3}, Cell{8, 1}, Cell{9, 6}}),
				Positions([]Cell{Cell{0, 3}, Cell{1, 0}, Cell{2, 4}, Cell{3, 7}, Cell{4, 9}, Cell{5, 2}, Cell{6, 6}, Cell{7, 8}, Cell{8, 1}, Cell{9, 5}}),
				Positions([]Cell{Cell{0, 4}, Cell{1, 0}, Cell{2, 3}, Cell{3, 8}, Cell{4, 6}, Cell{5, 1}, Cell{6, 9}, Cell{7, 2}, Cell{8, 5}, Cell{9, 7}}),
				Positions([]Cell{Cell{0, 5}, Cell{1, 0}, Cell{2, 2}, Cell{3, 9}, Cell{4, 7}, Cell{5, 1}, Cell{6, 3}, Cell{7, 8}, Cell{8, 6}, Cell{9, 4}}),
				Positions([]Cell{Cell{0, 6}, Cell{1, 0}, Cell{2, 2}, Cell{3, 5}, Cell{4, 7}, Cell{5, 9}, Cell{6, 3}, Cell{7, 8}, Cell{8, 4}, Cell{9, 1}}),
				Positions([]Cell{Cell{0, 8}, Cell{1, 0}, Cell{2, 2}, Cell{3, 7}, Cell{4, 5}, Cell{5, 1}, Cell{6, 9}, Cell{7, 4}, Cell{8, 6}, Cell{9, 3}}),
				Positions([]Cell{Cell{0, 7}, Cell{1, 0}, Cell{2, 2}, Cell{3, 5}, Cell{4, 8}, Cell{5, 6}, Cell{6, 9}, Cell{7, 3}, Cell{8, 1}, Cell{9, 4}}),
				Positions([]Cell{Cell{0, 9}, Cell{1, 0}, Cell{2, 3}, Cell{3, 5}, Cell{4, 2}, Cell{5, 8}, Cell{6, 1}, Cell{7, 7}, Cell{8, 4}, Cell{9, 6}}),
			},
		},
	}
	for _, c := range cases {

		sols := NQueenSolutions(uint(c.size))
		// if err != nil {
		// 	t.Errorf("GetGoQueenPositions with size %d got error %v", c.size, err)
		// }
		if !AssertEqual(sols, c.solutions) {
			t.Errorf("GetGoQueenPositions with same size %d (got) %v == %v (expected)", c.size, sols, c.solutions)
		}
	}
}

func AssertEqual(left, right []Positions) bool {
	var isInRight func(Positions) bool
	isInRight = func(l Positions) bool {
		lcell, lok := l.At(0)
		if lok {
			for _, r := range right {
				cell, ok := r.At(0)
				if ok && lcell.Row == cell.Row && l.Equal(r) {
					return true
				}
			}
		}
		return false
	}
	for _, l := range left {
		if !isInRight(l) {
			return false
		}
	}
	return true
}
