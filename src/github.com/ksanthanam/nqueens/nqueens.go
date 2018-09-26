// Package nqueens for Matrix related problems
package nqueens

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	p      = fmt.Println
	pf     = fmt.Printf
	s      = fmt.Sprintf
	digits = 2
)

var (
	// DEBUG ON/OF
	DEBUG = false
	// LEVEL for DEBUG
	LEVEL = 1
	d     = func(level int, str string) {
		if DEBUG && level <= LEVEL {
			p(str)
		}
	}
	df = func(level int, str string) {
		if DEBUG && level <= LEVEL {
			pf(strings.Join([]string{str, "\n"}, ""))
		}
	}
	dfr = func(level int, str string) {
		if DEBUG && level <= LEVEL {
			pf(strings.Join([]string{str, "                             \r"}, ""))
		}
	}
)

// StackAction actions on PieceStack
type StackAction int

const (
	// NoAction type
	NoAction StackAction = iota
	// NextRow progress to Next Row
	NextRow
	// PopStack Pop the stack
	PopStack
	// Traversed the row has been traversed
	Traversed
	// RowHasNoPosition to indicate can't place the Piece
	RowHasNoPosition
)

// Placeable type to indicate a Piece is placeable or not
type Placeable bool

// String function
func (p Placeable) String() string {
	if p {
		return "Placeable"
	} else {
		return "Not Placeable"
	}
}
func (sa StackAction) String() string {
	size := 20
	sizeText := func(display string) string {
		format := fmt.Sprintf("%%%ds", size)
		return fmt.Sprintf(format, display)
	}
	switch sa {
	case NoAction:
		return sizeText("NoAction")
	case NextRow:
		return sizeText("NextRow")
	case PopStack:
		return sizeText("PopStack")
	case Traversed:
		return sizeText("Traversed")
	case RowHasNoPosition:
		return sizeText("RowHasNoPosition")
	default:
		return sizeText("")
	}
}

/*
 *   Cell Type
 */

// Cell type for each cell on a Chess board
type Cell struct {
	Row, Col int
}

// String to help with Printing
func (c Cell) String() string {
	format := fmt.Sprintf("(%%%dd,%%%dd)", digits, digits)
	return fmt.Sprintf(format, c.Row, c.Col)
}

// Equal function to compare two cells
func (c Cell) Equal(o Cell) bool {
	return c.Row == o.Row && c.Col == o.Col
}

// AtFirstRow function to compare two cells
func (c Cell) AtFirstRow() bool {
	return c.Row == 0
}

/*
 *   Stack
 */

// PieceStack Interface for ChessBoard
type PieceStack interface {
	processAnchor(Cell, int, chan Solution, *sync.WaitGroup)
	canPieceBePlaced(Cell) piececanbeplaced
	GetStacks() []Piece
}

// QueenStack type that converts CoordStack
type QueenStack struct {
	lock      sync.Mutex
	size      int
	pieces    []Piece
	visited   map[Cell]bool
	traversed map[int]bool
	processed bool
}

func (qs *QueenStack) canPieceBePlaced(anchor Cell) piececanbeplaced {
	return func(queens []Piece, queen Piece) bool {
		position := queen.GetPosition()
		if position.Row >= qs.size || position.Col >= qs.size {
			return false
		}
		blocked := false
		for _, qpos := range queens {
			blocked = blocked || position.InPathOfQueen(qpos.GetPosition())
		}
		return !blocked
	}
}

// GetStacks func
func (qs *QueenStack) GetStacks() []Piece {
	qs.lock.Lock()
	defer qs.lock.Unlock()
	return qs.pieces
}

func (qs *QueenStack) hasBeenProcessed() bool {
	qs.lock.Lock()
	defer qs.lock.Unlock()
	return qs.processed
}

func (qs *QueenStack) setProcessed(flag bool) {
	qs.lock.Lock()
	defer qs.lock.Unlock()
	qs.processed = flag
}
func (qs *QueenStack) getPieces() Positions {
	qs.lock.Lock()
	defer qs.lock.Unlock()
	positions := make([]Cell, 0)
	for _, piece := range qs.pieces {
		positions = append(positions, piece.GetPosition())
	}
	qs.processed = true
	return Positions(positions)
}

func (qs *QueenStack) popStack(anchor Cell) StackAction {
	qs.lock.Lock()
	var action StackAction
	defer func() {
		if action == Traversed {
			qs.traversed[anchor.Row] = true
		}
		//        12345678901234567890123456789012345678901234567890
		dfr(1, s("     %s: popStack - %s with result %d long", anchor, action, len(qs.pieces)))

		qs.lock.Unlock()
	}()
	poppable := func(r int) bool {
		return r <= anchor.Row && len(qs.pieces) > 1
	}
	newQueen := func(row, col int) Piece {
		return NewQueen(Cell{row, col})
	}
	clearDown := func() {
		for row := len(qs.pieces); row <= anchor.Row; row++ { // 3 to 3
			for col := 0; col < qs.size; col++ {
				qs.visited[Cell{row, col}] = false // f(3,0), f(3,1), f(3,2), f(3,3)
			}
		}
	}
	if qs.traversed[anchor.Row] {
		action = Traversed
		return action
	}
	stackSize := len(qs.pieces) // [(0,3) (1,0) (2,2)], [(0,3) (1,0)], [(0,3) (1,1)]
	if stackSize <= 1 {         // 3 <= 1, 2 <= 1, 2 <= 1
		action = Traversed
		return action
	}
	can := qs.canPieceBePlaced(anchor)
	r := stackSize - 1 // 2, 1, 2
	clearDown()
	for poppable(r) { // 2, 1 <= 3 && 2 > 1, 2 <= 3 && 2 > 1
		lastPiece := qs.pieces[stackSize-1]         // (2,2), (1,0), (1,1)
		r = lastPiece.GetRow()                      // 2, 1, 1
		c := lastPiece.GetCol() + 1                 // 3, 1, 2
		qs.pieces = qs.pieces[:stackSize-1]         // [(0,3) (1,0)], [(0,3)], [(0,3)]
		stackSize--                                 // 2,1,1
		nextPosition := Cell{r, c}                  // (1,1), (1,2)
		queen := newQueen(r, c)                     // Q(1,1), Q(1,2)
		for c < qs.size && !can(qs.pieces, queen) { // 3 < 4, 1 < 4 && can(1,1), 2 < 4 && cant(1,2), 3 < 4 && cant(1,3)
			qs.visited[nextPosition] = true // t(2,3), t(1,2), t(1,3)
			c++                             // 3, 2, 3, 4
			nextPosition := Cell{r, c}      // (1,2), (1,3), (1,4)
			queen = NewQueen(nextPosition)  // Q(1,2), Q(1,3), Q(1,4)
		}
		if c < qs.size && can(qs.pieces, newQueen(r, c)) { //  1 < 4 && can(1,1)
			qs.pieces = append(qs.pieces, newQueen(r, c)) // [(0,3) (1,1)]
			qs.visited[Cell{r, c}] = true                 // t(1,1)
			stackSize++                                   // 2
			if (stackSize - 1) >= anchor.Row {            // 1 >= 3
				action = Traversed
				return action
			}
			action = NextRow
			return action
		}

		if c >= qs.size { // 4 >= 4
			action = PopStack
			return action
		}
	}
	action = Traversed
	return action
}

func (qs *QueenStack) nextRow(anchor Cell) StackAction { // (3,3)

	qs.lock.Lock()
	var action StackAction
	defer func() {
		if action == Traversed {
			qs.traversed[anchor.Row] = true
			if anchor.RowIs(qs.size - 1) {
				qs.processed = true
			}
		}
		dfr(1, s("     %s: nextRow  - %s with result %d long", anchor, action, len(qs.pieces)))
		qs.lock.Unlock()
	}()
	newQueen := func(row, col int) Piece {
		return NewQueen(Cell{row, col})
	}
	if qs.processed || qs.traversed[anchor.Row] {
		action = Traversed
		return Traversed
	}

	stackSize := len(qs.pieces) // [(0,3)], [(0,3) (1,0)], [(0,3) (1,0) (2,2)], [(0,3) (1,1)]
	if stackSize == 0 {
		qs.pieces = append(qs.pieces, NewQueen(Cell{0, anchor.Col})) //[(0,3)]
		stackSize++                                                  // 1
		qs.visited[Cell{0, anchor.Col}] = true
	}
	if (stackSize - 1) >= anchor.Row { // 0 >= 3, 1 >= 3, 2 >= 3, 1 >= 3
		action = Traversed
		return action
	}
	r := stackSize // 1, 2, 3, 2
	c := 0
	nextPosition := Cell{r, c} // (1,0), (2,0), (3,0), (2,0)
	for c < qs.size {          // 0 < 4, 0 < 4
		nextPosition = Cell{r, c} // (1,0), (2,0), (3,0), (2,0)
		// d(1, s("     %s: O:Inspecting position %t%s with %s", anchor, qs.visited[nextPosition], nextPosition, qs.pieces))
		for c < qs.size && qs.visited[nextPosition] { // 0 < 4 && f(1,0), 0 < 4 && f(2,0), 0 < 4 && f(3,0),
			// 0 < 4 && f(2,0)
			c++
			nextPosition = Cell{r, c}
			// d(1, s("     %s: I:Inspecting position %t%s with %s", anchor, qs.visited[nextPosition], nextPosition, qs.pieces))
		}
		nextPosition = Cell{r, c}       // (1,0), (2,0), (3,0), (2,0)
		queen := NewQueen(nextPosition) // Q(1,0), Q(2,0), Q(3,0), Q(2,0)
		can := qs.canPieceBePlaced(anchor)
		// d(1, s("     %s: O:Can position %t%s with %s", anchor, !can(qs.pieces, queen), nextPosition, qs.pieces))
		for c < qs.size && !can(qs.pieces, queen) { // 0 < 4 && can(1,0), 0 < 4 && cant(2,0), 1 < 4 && cant(2,1),
			// 2 < 4 && can(2,2), 0 < 4 && cant(3,0), 1 < 4 && cant(3,1), 2 < 4 && cant(3,2), 3 < 4 && cant(3,3),
			// 0 < 4 && cant(2,0), 1 < 0 && cant(2,1), 2 < 4 && cant(2,2), 3 < 4 && cant(2,3)
			qs.visited[nextPosition] = true // t(2,0), t(2,1), t(3,0), t(3,1), t(3,2), t(3,3), t(2,0), t(2,1), t(2,2), t(2,3)
			c++                             // 1,2, 1, 2, 3, 4, 1, 2,3, 4
			nextPosition = Cell{r, c}       // (2,1), (2,2), (3,1), (3,2), (3,3), (3,4), (2,1), (2,2), (2,3), (2,4)
			queen = NewQueen(nextPosition)  // Q(2,1), Q(2,2), Q(3,1), Q(3,2), Q(3,3), Q(3,4), Q(2,1), Q(2,2), Q(2,3), Q(2,4)
			// d(1, s("     %s: O:Can position %t%s with %s", anchor, can(qs.pieces, queen), nextPosition, qs.pieces))
		}
		if c < qs.size && can(qs.pieces, newQueen(r, c)) { // 2 < 4 && can(2,2), 4 < 4, 4 < 4
			qs.pieces = append(qs.pieces, newQueen(r, c)) // [(0,3) (1,0)], [(0,3) (1,0) (2,2)]
			qs.visited[Cell{r, c}] = true                 // t(1,0), t(2,2)
			stackSize++                                   // 2, 3
			if (stackSize - 1) >= anchor.Row {            // 1 >= 3, 2 >= 3
				action = Traversed
				return action
			}
			action = NextRow
			return action
		}
	}

	action = PopStack
	return action
}

func (qs *QueenStack) traverseToAnchor(anchor Cell, solutions chan Solution, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() { // (3,0)
		st := time.Now()
		actions := make(chan StackAction, qs.size)
		// done := make(chan bool, 1)
		action := NextRow
		defer func() {
			dfr(1, s("     %s: - Stop with %s after %s", anchor, action, time.Since(st)))
			wg.Done()
		}()
		for {
			select {
			case action = <-actions:
				switch action {
				case NextRow:
					actions <- qs.nextRow(anchor)
				case PopStack:
					actions <- qs.popStack(anchor)
				case Traversed:
					if anchor.RowIs(qs.size - 1) {
						solutions <- Solution{anchor, qs.getPieces()}
					}
					return
				}
			default:
				switch action {
				case NextRow:
					actions <- qs.nextRow(anchor)
				case PopStack:
					actions <- qs.popStack(anchor)
				case Traversed:
					if anchor.RowIs(qs.size - 1) {
						solutions <- Solution{anchor, qs.getPieces()}
					}
					return
				}
			}
		}
	}()
}
func (qs *QueenStack) processAnchor(anchor Cell, size int, solutions chan Solution, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		qs.traverseToAnchor(anchor, solutions, wg)
	}()
	return
}

// NewQueenStack function for CellStack
func NewQueenStack(size int) PieceStack {
	return &QueenStack{sync.Mutex{}, size, make([]Piece, 0), make(map[Cell]bool, 0), make(map[int]bool, 0), false}
}

/*
 *   Board
 */

// Board interface
type Board interface {
	Size() int
}

/*
 *   ChessBoard
 */

// ChessBoard with size
type ChessBoard struct {
	size uint
}

// Size function
func (cb *ChessBoard) Size() int {
	return int(cb.size)
}

// Piece interface
type Piece interface {
	CanPlace(Cell) bool
	GetPosition() Cell
	withInRows(upper, lower int) bool
	GetRow() int
	GetCol() int
	String() string
}

type piececanbeplaced func([]Piece, Piece) bool

// Queen struct
type Queen struct {
	Name     string
	Position Cell
}

// String for print function
func (q *Queen) String() string {
	if q == nil {
		return ""
	}
	return s("%s%s", q.Name, q.Position)
}

// NewQueen function
func NewQueen(position Cell) Piece {
	return &Queen{"Queen", Cell{position.Row, position.Col}}
}

// CanPlace function
func (q Queen) CanPlace(cell Cell) bool {
	return !q.Position.InPathOfQueen(cell)
}

// GetPosition function
func (q *Queen) GetPosition() Cell {
	return q.Position
}

// withInRows function to check the row boundary
func (q Queen) withInRows(upper, lower int) bool {
	return q.Position.Row <= upper && q.Position.Row >= lower
}

// GetRow function to get Row of the piece
func (q Queen) GetRow() int {
	return q.Position.Row
}

// GetCol function to get Col of the piece
func (q Queen) GetCol() int {
	return q.Position.Col
}

/*
 *   QueenBoard
 */

// QueenBoard with only Queens
type QueenBoard struct {
	ChessBoard
	queens    []PieceStack
	processed map[int]bool
	solutions map[int]bool
}

// InPathOfQueen function to check if coordinate is in the path of a Queen coordinate
func (c Cell) InPathOfQueen(q Cell) bool {
	inpath := c.Row == q.Row || // Same Row
		c.Col == q.Col || // Same Column
		c.Row-c.Col == q.Row-q.Col || // Same Diff
		c.Row+c.Col == q.Row+q.Col // Same total
	return inpath
}

// RowIs function to find out the row
func (c Cell) RowIs(row int) bool {
	return c.Row == row
}

// ColIs function to find out the col
func (c Cell) ColIs(col int) bool {
	return c.Col == col
}

// NewQueenBoard function
func NewQueenBoard(size uint) *QueenBoard {
	stacks := make([]PieceStack, 0)
	for col := 0; col < int(size); col++ {
		stacks = append(stacks, NewQueenStack(int(size)))
	}
	return &QueenBoard{ChessBoard{size}, stacks, make(map[int]bool, 0), make(map[int]bool, 0)}
}

// PlaceQueen function to add to Queens placed
func (q *QueenBoard) PlaceQueen(row, col int) bool {
	// cell := Cell{row: row, col: col}
	return false
}

// String function to display QueenBoards
func (q *QueenBoard) String() string {
	return fmt.Sprintf("Board %dx%d Placed(%v) Solutions(%v)", q.size, q.size, q.processed, q.solutions)
}

func (q *QueenBoard) startAnchorQueue(anchors chan<- Cell, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for row := 0; row < q.Size(); row++ {
			for col := 0; col < q.Size(); col++ {
				anchors <- Cell{row, col}
			}
		}
	}()
}

// PlaceNQueens function to place N queens
func (q *QueenBoard) PlaceNQueens() []Positions {
	st := time.Now()
	// processed := 0
	size := q.Size()
	digits = 0
	for d := size; d >= 1; d = d / 10 {
		digits++
	}
	noOfAnchors := q.Size() * q.Size()
	wg := &sync.WaitGroup{}
	done := make(chan bool, 1)
	anchors := make(chan Cell, noOfAnchors)
	solutions := make(chan Solution, q.Size())
	solns := make([]Positions, 0)
	nth := func(n int) string {
		m := n % 10
		sfx := ""
		if n >= 11 && n <= 13 {
			sfx = "th"
		} else {
			switch m {
			case 1:
				sfx = "st"
			case 2:
				sfx = "nd"
			case 3:
				sfx = "rd"
			default:
				sfx = "th"
			}
		}
		format := fmt.Sprintf("%%%dd%%s", digits)
		return fmt.Sprintf(format, n, sfx)
	}

	q.startAnchorQueue(anchors, wg)
	for row := 0; row < q.Size(); row++ {
		for col := 0; col < q.Size(); col++ {
			anchor := <-anchors
			dfr(1, s("     %s: PROCESSING Anchor %s", anchor, anchor))
			q.queens[anchor.Col].processAnchor(anchor, q.Size(), solutions, wg)
		}
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		noOfSolutions := 0
		for {
			select {
			case solution := <-solutions:
				noOfSolutions++
				d(0, s("     %s Solution for %s Col took %20s and is %d long", nth(noOfSolutions), nth(solution.forAnchor.Col), time.Since(st), solution.Size()))
				if solution.Size() == q.Size() {
					solns = append(solns, solution.positions)
				} else {
					d(0, s("     %s Partial  for %s Col took %20s and is %d long", nth(noOfSolutions), nth(solution.forAnchor.Col), time.Since(st), solution.Size()))
				}
				if noOfSolutions == q.Size() {
					done <- true
				}
			case <-done:
				return
			}
		}
	}()
	wg.Wait()
	p("Found", len(solns), fmt.Sprintf("solutions in %s ", time.Since(st)), solns)
	return solns
}

// Positions type is a collection of Cells
type Positions []Cell

// Equal function to compare two Positions
func (p Positions) Equal(c Positions) bool {
	left := []Cell(p)
	right := []Cell(c)
	if len(left) != len(right) {
		return false
	}
	for i := 0; i < len(left); i++ {
		if !left[i].Equal(right[i]) {
			return false
		}
	}
	return true
}

// Size function to get size
func (p Positions) Size() int {
	positions := []Cell(p)
	return len(positions)
}

// At function to get size
func (p *Positions) At(i int) (Cell, bool) {
	positions := []Cell(*p)
	size := p.Size()
	if i < size {
		return positions[i], true
	}
	return Cell{}, false
}

// First function to get size
func (p *Positions) First() (Cell, bool) {
	return p.At(0)
}

// Last function to get size
func (p *Positions) Last() (Cell, bool) {
	size := p.Size()
	if size > 0 {
		return p.At(size - 1)
	}
	return Cell{}, false
}

// CanPlaceQueen function
func (p Positions) CanPlaceQueen(cell Cell) bool {
	positions := []Cell(p)
	inpath := false
	for _, qpos := range positions {
		inpath = inpath || cell.InPathOfQueen(qpos)
	}
	return !inpath
}

// Solution type
type Solution struct {
	forAnchor Cell
	positions Positions
}

// String function
func (s Solution) String() string {
	return fmt.Sprintf("%s: %s", s.forAnchor, s.positions)
}

// Size function
func (s Solution) Size() int {
	return s.positions.Size()
}

/*
	0   1   2   3
   _______________
0 |   | Q |   |   |
  |---------------|
1 |   |   |   | Q |
  |---------------|
2 | Q |   |   |   |
  |---------------|
3 |   |   | Q |   |
   ---------------
         x-y  y+x
   (1,2) (-1) (3)
         (0,1),(0,2),(0,3)
   (1,0),(1,1),(1,2),(1,3)
         (2,1),(2,2),(2,3)
    (3,0)      (3,2)
*/

// NQueenSolutions function
func NQueenSolutions(size uint) []Positions {

	qb := NewQueenBoard(size)
	return qb.PlaceNQueens()
}
