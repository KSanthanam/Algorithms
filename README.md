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

By implementing the go routines and channels, an efficiency of **10000x** was achieved.

<pre>
Anchor is Anchor(R<sub>i</sub>,C<sub>j</sub>) 
  action = nextRow(anchor)
  loop
    switch action {
      NextRow:
        nextRow(anchor)
      PopStack:
        popStack(anchor)
      Traversed:
        return
    }

nextRow function returns action
  if stackSize == 0 then
    stack = [(1,j)]
    stack.visited[1,j] = true
  end
  if stackSize >= i then
    stack.traversed[i] = true
    return Trversed
  end
  r = stackSize
  c = 1
  while c <= N do
    while c < N and stack.visited[(r,c)] do
      c++
    done
    while c < N && can not place (r,c) given stack of queens do
      stack.visited[(r,c)] = true
      c++
    done
    if c < N && can place (r,c) given stack of queens then
      add (r,c) to stack
      stack.visited[(r,c)] = true
      if stackSize  >= i then 
        stack.traversed[i] = true
        return Traversed
      end
      return NextRow
    end
  done
  return PopStack

popRow function returns action
  if stack.traversed[i] then
    return Traversed
  end
  if stackSize <= 1 then
    stack.traversed[i] = true
    return Traversed
  end
  r = stackSize 
  for rows r to N do
    for col 1 to N do
      stack.visited[(r,c)] = false 
    done
  done
  while r <= i and stackSize > 1 do
    r = lastPiece's row
    c = lastPiece's col + 1
    pop lastPiece from stack
    while c < N && cant place (r,c) for given queens in stack do
      stack.visited[(r,c)] = true
      c++
    done
    if c < N && can place (r,c) for given queens in stack then
      add (r,c) to stack
      if stackSize  >= i then 
        stack.traversed[i] = true
        return Traversed
      end
      return NextRow
    end
    if c >= N then 
      return PopStack
    end
  done
  stack.traversed[i] = true
  return Traversed
</pre>

### Summary
The Algorithm implemented here runs N*N go routines and traverses the stack for its column. To keep the integrity of the PieceStack the stack for the particular column is processed by a single anchor at a time. The go routine for other anchors for this column waits for the given anchor to traverse and backtrack (if needed) to the row of the given anchor.
