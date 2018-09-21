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

The solution implemented is the most efficient possible in Golang using back propagation, channel and go routines.

#### Synopsis of Algorithm

The Algorithm generates N*N anchor points and submits to the anchor channel.
R<sub>i</sub>,C<sub>j</sub>
where <sub>i</sub> is 1..N, <sub>j</sub> is 1..N

Each anchor is read from the channel and submited to go routine that  traverses the board using the PieceStack dedicated for the column of the given anchor.

The Traverse Logic for a given anchor.
when a anchor (R<sub>i</sub>,C<sub>j</sub>) is receied, a go routine is kick started with the aim to place Queens up to the row i. If the row i can be reached with placements it back tracks by forwarding the column position. If it exhausts all columns for rows before i, the go routine stops trying. 
so when the the last row Nth row anchor is received and the row can be reached, it is considered to be a solution as there will be N queens in place.

So the logic is as follows:
anchor go routine for (R<sub>i</sub>,C<sub>j</sub>)
Each anchor is submitted to a go routine. So there are approximately N<sup>4</sup> number of go routines. Most of these go routines will exit very quick if the predecessor anchor for the same column reached there before and done the work. 
so if the anchor for R<sub>i+m</sub> has reached before R<sub>i</sub> then the corresponding go routine exits straight away. if anchor for R<sub>n</sub> reaches before any of the other anchors for the same column then all subsequent anchors for the said column will exit immediately.

By implementing the go routines an efficiency of 3000X was achieved.

<pre>
Anchor is (R<sub>i</sub>,C<sub>j</sub>) 

if Col<sub>j</sub> already processed then
   exit
end
if R<sub>1</sub> then 
  place Queen in (R<sub>1</sub>,C<sub>j</sub>) 
  set (R<sub>0</sub>,C<sub></sub>) visited to True
end

if Stack is empty then
   place Queen in (R<sub>0</sub>,C<sub>j</sub>) 
end

lastQueen = lastPiece in Stack 
lastQueen is Queen(R<sub>si</sub>,C<sub>sj</sub>)

r = lastQueen Row + 1
colStart = 1
rowStop = i

for r <= rStop do
   c = colStart
   placed = false
   for c < N do
       if (r,c) has not been visited then
         visited[(r,c)] = true
         placed = Queen can be placed in (r,c) 
         if placed then
          push Queen(R<sub>i</sub>,C<sub>j</sub>)
          break
         end
       end
       c++
   done
   if placed then
      r++
      cStart = 0
   else
      if r <= rStop
        if Stack Size == 1 then
          break
        else
          pop last cell in Piece Stack
          r = popped piece row
          c = popped cell col + 1
        end
      else
        break
      end
   end
   if r < rowStop then
     c = j
     for r is r to rowStop do
       for c is 0 to N do
          visited[(r,c)] = true
       done
     done
   end
   if r == N then
      send the cells in PieceStack to Solutions List
      processed[Column] = true
   end
   
done
</pre>

### Summary
The Algorithm implemented here runs N*N go routines and traverses the stack for its column. To keep the integrity of the PieceStack the stack for the particular column is processed by a single anchor at a time. The go routine for other anchors for this column waits for the given anchor to traverse and backtrack (if needed) to the row of the given anchor.


[Markdown Syntax](https://stackedit.io/app#)