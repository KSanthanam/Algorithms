# Algorithms 
## NQueens
Given a Chessboard of size N find possible Solutions where N queens are placed on the board

### Time complexity
The **Time Complexity** by Brute force logic is N<sup>N</sup> as one needs to place the queen in (a,b) and check N*N-1 places to validate.

N Queen Algorithm

Example 4x4
All possible solutions are 

1) [(0,1) (1,3) (2,0) (3,2)] 

 | Q | 0 | 1 | 2 | 3 |
 |:--|:--|:--|:--|:--|
 | 0 |   | Q |   |   | 
 | 1 |   |   |   | Q |
 | 2 | Q |   |   |   |
 | 3 |   |   | Q |   |

2) [(0,2) (1,0) (2,3) (3,1)] 

 | Q | 0 | 1 | 2 | 3 |
 |:--|:--|:--|:--|:--|
 | 0 |   |   | Q |   | 
 | 1 | Q |   |   |   |
 | 2 |   |   |   | Q |
 | 3 |   | Q |   |   |

### Solutions
There are few solutions with inefficient timecomplexity. The worst being Q(N<sup>N</sup>)

The solution implemented here is the most efficient possible in Golang using back propagation, channel and go routines.

#### Synopsis of Algorithm

The Algorithm generates N*N anchor points and submits to the anchor channel.
R<sub>i</sub>,C<sub>j</sub>
where <sub>i</sub> is 1..N, <sub>j</sub> is 1..N

Each anchor is read from the channel and submited to go a routine that  traverses the board using the PieceStack dedicated for the column of the given anchor.

The Traverse Logic for a given anchor.<br>
when an anchor (R<sub>i</sub>,C<sub>j</sub>) is receied, a go routine is kick started with the aim to place Queens up to the row i. If the row i can't be reached with placements it back tracks by forwarding the column position. If it exhausts all columns for the rows before i, the go routine stops trying. 
so when the the last row anchor is received and the row can be reached, it is considered to be a solution as there will be N queens in place.

So the logic is as follows:
anchor go routine for (R<sub>i</sub>,C<sub>j</sub>)

Each anchor is submitted to a go routine. So there are approximately N<sup>4</sup> number of go routines. Most of these go routines will exit very quick if the predecessor anchor for the same column reached there before and done the work. 
so if an anchor for R<sub>i+m</sub> has reached before R<sub>i</sub> then the corresponding go routine exits straight away. if anchor for R<sub>n</sub> reaches before any of the other anchors for the same column then all subsequent anchors for the said column will exit immediately.

By implementing the go routines, an efficiency of 5000X was achieved.

<pre>
Anchor is Anchor(R<sub>i</sub>,C<sub>j</sub>) 

if PieceStack for Anchor.Col not processed then

  lastPiece = Last Piece in PieceStack
  nextPiece = NextPiece(lastPiece.Row + 1, 0)

  for nextPiece.Row <= Anchor.Row && 
      Anchor.Row not already traversed do
      c is nextPiece.Col
      placed = false
      for c < size do
        if Position nextPiece.Row,c is not in path then 
          place Queen at nextPiece.Row,c
          break
        end
        c++
      done
      if placed then
        nextPiece = nextPiece.Row + 1, 0
      else 
        pop a piece from PieceStack (up to 1 piece at Row 0)
        if piece not popped then
          break
        end
        nextPiece = poppedPiece.Row, poppedPiece.Col + 1
      end
  done
  Anchor.Row reached if nextPiece.Row >= Anchor.Row else false

  if Anchor.Row == N OR Anchor.Row not reached then
    send Pieces in PieceStack as Solution if length of PieceStack == size
    Set Anchor.Col Processed to true
  end
end
</pre>

### Summary
The Algorithm implemented here runs N*N go routines and traverses the stack for its column. To keep the integrity of the PieceStack the stack for the particular column is processed by a single anchor at a time. The go routine for other anchors for this column waits for the given anchor to traverse and backtrack (if needed) to the row of the given anchor.


[Markdown Syntax](https://stackedit.io/app#)