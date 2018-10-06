package knapsack

import "testing"

func TestGetKnapSack(t *testing.T) {
	cases := []struct {
		objects  Objects
		max      int
		selected Objects
	}{
		{[]Object{Object{2, 1}, Object{3, 2}, Object{4, 5}, Object{5, 6}},
			8,
			[]Object{Object{5, 6}, Object{3, 2}}},
	}
	for _, c := range cases {
		got := GetKnapSack(c.objects, c.max)
		if !got.Equal(c.selected) {
			t.Errorf("GetKnapSack (got) %v == %v (expected)", got, c.selected)
		}
	}
}
