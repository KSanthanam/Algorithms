// Package nqueens for Matrix related problems
package nqueens

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	p = fmt.Println
	s = fmt.Sprintf
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
)

/*
 *   Cell Type
 */

// Cell type for each cell on a Chess board
type Cell struct {
	Row, Col int
}

// String to help with Printing
func (c Cell) String() string {
	return fmt.Sprintf("(%d,%d)", c.Row, c.Col)
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

// Stack Interface for ChessBoard
type Stack interface {
	Push(Piece)
	Pop() (Piece, error)
	Place(Piece, Cell) bool
	processAnchor(Cell, int, chan Positions, piececanbeplaced, *sync.WaitGroup)
	canPlaceQueenAmonstQueens(Cell) piececanbeplaced
	GetStacks() []Piece
}

// PieceStack type that converts CoordStack
type PieceStack struct {
	lock      sync.Mutex
	pieces    []Piece
	visited   map[Cell]bool
	processed bool
}

// Place function
func (ps *PieceStack) Place(piece Piece, cell Cell) bool {
	return false
}

func (ps *PieceStack) isEmpty() bool {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	return len(ps.pieces) == 0
}

func (ps *PieceStack) canPlaceQueenAmonstQueens(anchor Cell) piececanbeplaced {
	return func(queens []Piece, queen Piece) bool {
		blocked := false
		position := queen.GetPosition()
		for _, qpos := range queens {
			blocked = blocked || position.InPathOfQueen(qpos.GetPosition())
		}
		return !blocked
	}
}

func (ps *PieceStack) processAnchor(anchor Cell, size int, solutions chan Positions, canPieceBePlaced piececanbeplaced, wg *sync.WaitGroup) {

	wg.Add(1)
	go func() {
		defer func() {
			ps.lock.Unlock()
			wg.Done()
		}()
		ps.lock.Lock()
		stackSize := len(ps.pieces)
		if ps.processed {
			return
		}
		if anchor.AtFirstRow() && stackSize == 0 {
			if stackSize == 0 {
				ps.pieces = append(ps.pieces, NewQueen(anchor))
				ps.visited[anchor] = true
			}
			return
		}
		r := 0
		cStart := 0

		if stackSize == 0 {
			firstQueen := NewQueen(Cell{0, anchor.Col})
			ps.pieces = append(ps.pieces, firstQueen)
			ps.visited[Cell{0, anchor.Col}] = true
			r = 1
			stackSize = len(ps.pieces)
			d(1, s("%s: Inserted first %s", anchor, firstQueen))
		}
		lastPiece := ps.pieces[stackSize-1]
		lastQueen := lastPiece.GetPosition()
		r = lastQueen.Row + 1
		rStop := anchor.Row
		st := time.Now()
		for r <= rStop {
			c := cStart
			placed := false
			for c < size {
				newPosition := Cell{r, c}
				if visited, ok := ps.visited[newPosition]; ok || !visited {
					newQueen := NewQueen(newPosition)
					placed = canPieceBePlaced(ps.pieces, newQueen)
					ps.visited[newPosition] = true
					if placed {
						d(1, s("%s: Placed %s", anchor, newQueen))
						ps.pieces = append(ps.pieces, newQueen)
						break
					}
				}
				c++
			}
			if placed {
				r++
				cStart = 0
			} else {
				if len(ps.pieces) > 0 && r <= rStop {
					stackSize := len(ps.pieces)
					if stackSize == 1 {
						d(1, s("%s: Reached top of the column(%d)", anchor, anchor.Col))
						break
					} else {
						popped := ps.pieces[stackSize-1]
						d(1, s("%s: Popping %s", anchor, popped))
						ps.pieces = ps.pieces[:stackSize-1]
						position := popped.GetPosition()
						r = position.Row
						cStart = position.Col + 1
					}
				} else {
					d(1, s("%s: Nothing to pop", anchor))
					break
				}
			}
		}
		d(0, s("%s: took %s", anchor, time.Since(st)))

		if r < rStop {
			d(0, s("%s: Blanking out rows from %d  to %d", anchor, r, rStop))
			c := anchor.Col
			for vr := r; vr <= rStop; vr++ {
				cell := Cell{vr, c}
				ps.visited[cell] = true
			}
		}
		if anchor.RowIs(size - 1) {
			// d(0, s("%s: Solution for column %d and took %s resulting %s", anchor, anchor.Col, time.Since(st), ps.pieces))
			cells := make([]Cell, 0)
			for _, piece := range ps.pieces {
				cells = append(cells, piece.GetPosition())
			}
			solutions <- Positions(cells)
			ps.processed = true
		}
	}()

	return
}

// NewPieceStack function for CellStack
func NewPieceStack() Stack {
	return &PieceStack{sync.Mutex{}, make([]Piece, 0), make(map[Cell]bool, 0), false}
}

// Push is a function to push to Stack
func (ps *PieceStack) Push(p Piece) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	ps.pieces = append(ps.pieces, p)
}

// Pop is a function to pop a stack
func (ps *PieceStack) Pop() (Piece, error) {
	ps.lock.Lock()
	ps.lock.Lock()
	defer ps.lock.Unlock()
	l := len(ps.pieces)
	if l == 0 {
		return nil, errors.New("Stack is empty")
	}
	popped := ps.pieces[l-1]
	ps.pieces = ps.pieces[:l-1]
	return popped, nil
}

// GetStacks function returns the stack
func (ps *PieceStack) GetStacks() []Piece {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	return ps.pieces
}

/*
 *   Function Types
 */

// Placeable function type
type Placeable func(*PieceStack, Piece, Cell) bool

// PiecePlaceable function type
type PiecePlaceable func(Piece, Cell) func(*PieceStack, Piece, Cell) bool

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
	// Place(Cell, func(piece Piece) bool) bool
	CanPlace(Cell) bool
	GetPosition() Cell
	withInRows(upper, lower int) bool
	GetRow() int
}

type piececanbeplaced func([]Piece, Piece) bool

// Queen struct
type Queen struct {
	Name     string
	Position Cell
}

// String for print function
func (q Queen) String() string {
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

/*
 *   QueenBoard
 */

// QueenBoard with only Queens
type QueenBoard struct {
	ChessBoard
	queens    []Stack
	processed map[int]bool
	solutions map[int]bool
}

// CanPlace function to get the Can function
func (q *QueenBoard) CanPlace(sid int, cell Cell) func(Cell) bool {
	return func(cell Cell) bool {
		piece := NewQueen(cell)
		q.queens[sid].Place(piece, cell)
		return false
	}
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
	stacks := make([]Stack, 0)
	for col := 0; col < int(size); col++ {
		stacks = append(stacks, NewPieceStack())
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

// traverseWith
func (q *QueenBoard) traverseWith(anchors <-chan Cell, solutions chan<- Positions) {
	// for anchor := range anchors {
	// 	solutions <- Positions([]Cell{Cell{anchor.row, anchor.col}})
	// }
	anchor := <-anchors
	solutions <- Positions([]Cell{Cell{anchor.Row, anchor.Col}})
}

func (q *QueenBoard) startAnchorQueue(anchors chan<- Cell, done <-chan bool, wg *sync.WaitGroup) {
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

func (q *QueenBoard) processStack(anchor Cell, solutions chan Positions, wg *sync.WaitGroup) {
	defer wg.Done()
	col := anchor.Col
	if processed, ok := q.processed[col]; !ok || !processed {
		q.queens[col].processAnchor(anchor, q.Size(), solutions, q.queens[col].canPlaceQueenAmonstQueens(anchor), wg)
	}
}
func (q *QueenBoard) collectSolutions(anchors <-chan Cell, solutions chan Positions, done chan bool, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case cell := <-anchors:
				wg.Add(1)
				go q.processStack(cell, solutions, wg)
			case <-done:
				return
			}
		}
	}()
	return
}

// PlaceNQueens function to place N queens
func (q *QueenBoard) PlaceNQueens() []Positions {
	st := time.Now()
	processed := 0
	noOfAnchors := q.Size() * q.Size()
	wg := &sync.WaitGroup{}
	done := make(chan bool, 1)
	anchors := make(chan Cell, noOfAnchors)
	solutions := make(chan Positions, q.Size())
	solns := make([]Positions, 0)

	defer func() {
		d(0, s("Processed %d solutions with %d anchors and they took %s", processed, noOfAnchors, time.Since(st)))
		p(s("%dx%d ChessBoard has %d solutions (generated in %s) and they are:\n%v", q.Size(), q.Size(), len(solns), time.Since(st), solns))
	}()

	q.startAnchorQueue(anchors, done, wg)
	q.collectSolutions(anchors, solutions, done, wg)
	go func() {
		for {
			select {
			case positions := <-solutions:
				processed++
				cells := []Cell(positions)
				col := cells[0].Col
				q.processed[col] = true
				if positions.Size() == q.Size() {
					q.solutions[col] = true
					solns = append(solns, positions)
					d(0, s("%d: processed for col(%d) \nprocessed(%d/%d) %v \nsolution(%d/%d)  %v",
						processed, col, len(q.processed), q.Size(), q.processed, len(q.solutions), q.Size(), q.solutions))
				}
				if processed == q.Size() {
					done <- true
					return
				}
			}
		}
	}()
	wg.Wait()
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

// isComplete function to get full size
func (p *Positions) isComplete(size int) bool {
	positions := []Cell(*p)
	return len(positions) == size
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
