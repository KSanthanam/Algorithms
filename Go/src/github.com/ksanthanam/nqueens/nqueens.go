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
	nth    = func(n, size int) string {
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
		format := fmt.Sprintf("%%%dd%%s", size)
		return fmt.Sprintf(format, n, sfx)
	}
)

var (
	wg      = &sync.WaitGroup{}
	solnset = SolutionDisplay{sync.Mutex{}, 4, 2, make(map[int]bool)}
	logs    = make(chan string, 1000)
	done    = make(chan bool, 1)
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
	logn = func(str string) string {
		return strings.Join([]string{str, "\n"}, "")
	}
	logr = func(str string) string {
		return strings.Join([]string{str, "                             \r"}, "")
	}
	logger = func(done chan bool) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case str := <-logs:
					pf(str)
				case <-done:
					return

				}

			}
		}()
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
)

// Placeable type to indicate a Piece is placeable or not
type Placeable bool

// String function
func (p Placeable) String() string {
	if p {
		return "Placeable"
	}
	return "Not Placeable"
}
func (sa StackAction) String() string {
	size := 9
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
	processAnchor(Cell, int, chan Solution)
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
	semaphore chan struct{}
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
	pieces := qs.pieces
	qs.lock.Unlock()
	return pieces
}

// String function
func (qs *QueenStack) String() string {
	length := len(qs.pieces)
	blank := qs.size - length
	bars := make([]string, 0)
	for c := 0; c < length; c++ {
		bars = append(bars, "⚄")
	}
	for c := 0; c < blank; c++ {
		bars = append(bars, "⚀")
	}
	bar := strings.Join(bars, "")
	return bar
}
func (qs *QueenStack) hasBeenProcessed() bool {
	qs.lock.Lock()
	processed := qs.processed
	qs.lock.Unlock()
	return processed
}

func (qs *QueenStack) setProcessed(flag bool) {
	qs.lock.Lock()
	qs.processed = flag
	qs.lock.Unlock()
}
func (qs *QueenStack) getPieces() Positions {
	qs.lock.Lock()
	positions := make([]Cell, 0)
	for _, piece := range qs.pieces {
		positions = append(positions, piece.GetPosition())
	}
	qs.processed = true
	result := Positions(positions)
	qs.lock.Unlock()
	return result
}

func (qs *QueenStack) popStack(anchor Cell) StackAction {
	qs.lock.Lock()
	var action StackAction
	debug := func() {
		wg.Add(1)
		go func() {
			level := 1
			if DEBUG && level <= LEVEL {
				logs <- logr(s(" %s: %10s count %s %s", anchor, "popStack", qs.String(), solnset.String()))
			}
			wg.Done()
		}()
	}
	defer func() {
		if action == Traversed {
			qs.traversed[anchor.Row] = true
		}
		debug()
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
	debug := func() {
		wg.Add(1)
		go func() {
			level := 1
			if DEBUG && level <= LEVEL {
				logs <- logr(s(" %s: %10s count %s %s", anchor, "nextRow", qs.String(), solnset.String()))
			}
			wg.Done()
		}()
	}

	defer func() {
		if action == Traversed {
			qs.traversed[anchor.Row] = true
			if anchor.RowIs(qs.size - 1) {
				qs.processed = true
			}
		}
		debug()
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
		nextPosition = Cell{r, c}                     // (1,0), (2,0), (3,0), (2,0)
		for c < qs.size && qs.visited[nextPosition] { // 0 < 4 && f(1,0), 0 < 4 && f(2,0), 0 < 4 && f(3,0),
			// 0 < 4 && f(2,0)
			c++
			nextPosition = Cell{r, c}
		}
		nextPosition = Cell{r, c}       // (1,0), (2,0), (3,0), (2,0)
		queen := NewQueen(nextPosition) // Q(1,0), Q(2,0), Q(3,0), Q(2,0)
		can := qs.canPieceBePlaced(anchor)
		for c < qs.size && !can(qs.pieces, queen) { // 0 < 4 && can(1,0), 0 < 4 && cant(2,0), 1 < 4 && cant(2,1),
			// 2 < 4 && can(2,2), 0 < 4 && cant(3,0), 1 < 4 && cant(3,1), 2 < 4 && cant(3,2), 3 < 4 && cant(3,3),
			// 0 < 4 && cant(2,0), 1 < 0 && cant(2,1), 2 < 4 && cant(2,2), 3 < 4 && cant(2,3)
			qs.visited[nextPosition] = true // t(2,0), t(2,1), t(3,0), t(3,1), t(3,2), t(3,3), t(2,0), t(2,1), t(2,2), t(2,3)
			c++                             // 1,2, 1, 2, 3, 4, 1, 2,3, 4
			nextPosition = Cell{r, c}       // (2,1), (2,2), (3,1), (3,2), (3,3), (3,4), (2,1), (2,2), (2,3), (2,4)
			queen = NewQueen(nextPosition)  // Q(2,1), Q(2,2), Q(3,1), Q(3,2), Q(3,3), Q(3,4), Q(2,1), Q(2,2), Q(2,3), Q(2,4)
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
func (qs *QueenStack) pipe(anchor Cell, inAction StackAction, outAction chan StackAction) {
	wg.Add(1)
	go func(anchor Cell, action StackAction) {
		switch action {
		case NextRow:
			outAction <- qs.nextRow(anchor)
		case PopStack:
			outAction <- qs.popStack(anchor)
		}
		wg.Done()
	}(anchor, inAction)
	return
}
func (qs *QueenStack) traverseToAnchor(anchor Cell, solutions chan Solution) {
	wg.Add(1)
	qs.semaphore <- struct{}{}
	go func() { // (3,0)
		defer func() {
			<-qs.semaphore
			wg.Done()
		}()
		debug := func() {
			wg.Add(1)
			go func() {

				level := 1
				if DEBUG && level <= LEVEL {
					logs <- logr(s(" %s: %10s count %s %s", anchor, "nextAction", qs.String(), solnset.String()))
				}
				wg.Done()
			}()
		}

		nextAction := NextRow
		next := make(chan StackAction, qs.size)
		for {
			qs.pipe(anchor, nextAction, next)
			nextAction = <-next
			debug()
			if nextAction == Traversed {
				if anchor.RowIs(qs.size - 1) {
					solutions <- Solution{anchor, qs.getPieces()}
				}
				return
			}
		}
	}()
}
func (qs *QueenStack) processAnchor(anchor Cell, size int, solutions chan Solution) {
	wg.Add(1)
	go func() {
		qs.traverseToAnchor(anchor, solutions)
		wg.Done()
	}()
	return
}

// NewQueenStack function for QueenStack
func NewQueenStack(size int) PieceStack {
	return &QueenStack{sync.Mutex{}, size, make([]Piece, 0), make(map[Cell]bool, 0), make(map[int]bool, 0), make(chan struct{}, 1), false}
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

func (q *QueenBoard) startAnchorQueue(anchors chan<- Cell) {
	wg.Add(1)
	go func() {
		for row := 0; row < q.Size(); row++ {
			for col := 0; col < q.Size(); col++ {
				anchors <- Cell{row, col}
			}
		}
		wg.Done()
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
	solnset.SetSize(q.Size())
	if DEBUG {
		logger(done)
	}

	noOfAnchors := q.Size() * q.Size()
	anchors := make(chan Cell, noOfAnchors)
	solutions := make(chan Solution, q.Size())
	solns := make([]Positions, 0)

	q.startAnchorQueue(anchors)
	for row := 0; row < q.Size(); row++ {
		for col := 0; col < q.Size(); col++ {
			anchor := <-anchors
			q.queens[anchor.Col].processAnchor(anchor, q.Size(), solutions)
		}
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		debug := func(solution Solution, noOfSolutions int) {
			if !DEBUG || 0 < LEVEL {
				return
			}
			wg.Add(1)
			go func() {
				filler := s(s("%%%ds", q.Size()*2), " ")
				length := s(s("%%%dd", digits), solution.Size())
				if solution.Size() == q.Size() {
					logs <- logn(s(" %s Solution for %s Col took %18s and is %s long %s  \n", nth(noOfSolutions, digits), nth(solution.forAnchor.Col, digits), time.Since(st), length, filler))
				} else {
					logs <- logn(s(" %s Partial  for %s Col took %18s and is %s long %s  \n", nth(noOfSolutions, digits), nth(solution.forAnchor.Col, digits), time.Since(st), length, filler))
				}
				wg.Done()
			}()
		}
		noOfSolutions := 0
		for {
			select {
			case solution := <-solutions:
				noOfSolutions++
				debug(solution, noOfSolutions)
				if solution.Size() == q.Size() {
					solns = append(solns, solution.positions)
				}
				solnset.SetSolution(solution.forAnchor.Col)
				if noOfSolutions == q.Size() {
					done <- true
					return
				}
			}
		}
	}()
	wg.Wait()
	p("Found", len(solns), fmt.Sprintf("solutions for Board %dx%d  in %s \n", q.Size(), q.Size(), time.Since(st)), solns)
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

// SolutionDisplay type
type SolutionDisplay struct {
	lock   sync.Mutex
	size   int
	digits int
	solns  map[int]bool
}

// SetSize function
func (sd *SolutionDisplay) SetSize(size int) {
	sd.lock.Lock()
	digits = 0
	for d := size; d >= 1; d = d / 10 {
		digits++
	}
	sd.size = size
	sd.digits = digits
	sd.lock.Unlock()
}

// SetSolution function
func (sd *SolutionDisplay) SetSolution(col int) {
	sd.lock.Lock()
	sd.solns[col] = true
	sd.lock.Unlock()
}

// String function
func (sd *SolutionDisplay) String() string {
	sd.lock.Lock()
	solnsdisplay := make([]string, 0)
	for pos := 0; pos < sd.size; pos++ {
		if sd.solns[pos] {
			solnsdisplay = append(solnsdisplay, s(s("%%%ds", sd.digits), "x"))
		} else {
			solnsdisplay = append(solnsdisplay, s(s("%%%dd", sd.digits), pos))
		}
	}
	sd.lock.Unlock()
	return strings.Join(solnsdisplay, "")
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
