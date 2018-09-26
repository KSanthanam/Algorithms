# Algorithms 
## NQueens
Given a Chessboard of size N, find possible Solutions where N queens are placed on the board

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

Anchors for each Column

Column 0 Anchors are (0,0),(1,0),(2,0),(3,0)
Following are the stages processAnchor goes through

<table>
<tr><th>Stage 1 </th><th>Stage 2</th><th>Stage 3</th></tr>
<tr><td>

 | Q | 0 | 1 | 2 | 3 |
 |:--|:--|:--|:--|:--|
 | 0 | Q |   |   |   | 
 | 1 |   |   | Q |   |
 | 2 |   |   |   |   |
 | 3 |   |   |   |   |

</td><td>

 | Q | 0 | 1 | 2 | 3 |
 |:--|:--|:--|:--|:--|
 | 0 | Q |   |   |   | 
 | 1 |   |   | Q |   |
 | 2 |   |   |   |   |
 | 3 |   |   |   |   |

</td><td>

 | Q | 0 | 1 | 2 | 3 |
 |:--|:--|:--|:--|:--|
 | 0 | Q |   |   |   | 
 | 1 |   |   |   | Q |
 | 2 |   | Q |   |   |
 | 3 |   |   |   |   |

</td></tr></table>

Column 1 Anchors are (0,1),(1,1),(2,1),(3,1)
Following are the stages processAnchor goes through

<table>
<tr><th>Stage 1 </th></tr>
<tr><td>

 | Q | 0 | 1 | 2 | 3 |
 |:--|:--|:--|:--|:--|
 | 0 |   | Q |   |   | 
 | 1 |   |   |   | Q |
 | 2 | Q |   |   |   |
 | 3 |   |   | Q |   |

</td></tr></table>


Column 2 Anchors are (0,2),(1,2),(2,2),(3,2)
Following are the stages processAnchor goes through

<table>
<tr><th>Stage 1 </th></tr>
<tr><td>

 | Q | 0 | 1 | 2 | 3 |
 |:--|:--|:--|:--|:--|
 | 0 |   |   | Q |   | 
 | 1 | Q |   |   |   |
 | 2 |   |   |   | Q |
 | 3 |   | Q |   |   |

</td></tr></table>

Column 3 Anchors are (0,3),(1,3),(2,3),(3,3)
Following are the stages processAnchor goes through

<table>
<tr><th>Stage 1 </th><th>Stage 2</th><th>Stage 3</th></tr>
<tr><td>

 | Q | 0 | 1 | 2 | 3 |
 |:--|:--|:--|:--|:--|
 | 0 |   |   |   | Q | 
 | 1 | Q |   |   |   |
 | 2 |   |   | Q |   |
 | 3 |   |   |   |   |

</td><td>

 | Q | 0 | 1 | 2 | 3 |
 |:--|:--|:--|:--|:--|
 | 0 |   |   |   | Q | 
 | 1 |   | Q |   |   |
 | 2 |   |   |   |   |
 | 3 |   |   |   |   |

</td><td>

 | Q | 0 | 1 | 2 | 3 |
 |:--|:--|:--|:--|:--|
 | 0 |   |   |   | Q | 
 | 1 |   |   |   |   |
 | 2 |   |   |   |   |
 | 3 |   |   |   |   |

</td></tr></table>


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
for row = 1 to N do
  for col = 1 to N do
    send Anchor(R<sub>row</sub>,C<sub>col</sub>) to Anchors Channel
  done
done

for each Anchor in Anchor Channel processAnchor(Anchor)

processAnchor where Anchor is Anchor(R<sub>i</sub>,C<sub>j</sub>) 
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
Usage: bin/algos -debug=false -level=1 -size=4
where debug is true/false
      level is debug level
      size is N

The Algorithm implemented here runs N*N go routines and traverses the stack for its column. To keep the integrity of the PieceStack the stack for the particular column is processed by a single anchor at a time. The go routine for other anchors for this column waits for the given anchor to traverse and backtrack (if needed) to the row of the given anchor.

### Results
Found 2 solutions for Board 4x4  in 468.802µs

Found 10 solutions for Board 10x10  in 4.869127ms

Found 15 solutions for Board 15x15  in 107.873553ms

Found 19 solutions for Board 19x19  in 2.172028907s

Found 20 solutions for Board 20x20  in 22.19845266s

Found 21 solutions for Board 21x21  in 14.093889882s

Found 22 solutions for Board 22x22  in 3m8.988134421s

Found 25 solutions for Board 25x25  in 1m24.612503551s

      1st Solution for 28th Col took      9m34.124887762s and is 30 long
      2nd Solution for 29th Col took     11m48.831587143s and is 30 long
      3rd Solution for 27th Col took     24m27.320351249s and is 30 long
      4th Solution for 25th Col took     30m11.097575441s and is 30 long
      5th Solution for 26th Col took     32m11.920358783s and is 30 long
      6th Solution for 24th Col took      45m17.77867037s and is 30 long
      7th Solution for 23rd Col took     46m49.624517124s and is 30 long
      8th Solution for 17th Col took     56m56.317268522s and is 30 long
      9th Solution for 21st Col took     1h1m8.825636808s and is 30 long
     10th Solution for 14th Col took   1h46m37.190281609s and is 30 long
     11th Solution for 15th Col took   1h54m47.518960064s and is 30 long
     12th Solution for 20th Col took   1h55m48.047302863s and is 30 long
     13th Solution for 16th Col took   1h57m27.719530822s and is 30 long
     14th Solution for 19th Col took   1h59m40.918951655s and is 30 long
     15th Solution for 12th Col took    2h0m31.994823474s and is 30 long
     16th Solution for 11th Col took     2h2m15.24557813s and is 30 long
     17th Solution for  9th Col took    2h3m43.981316128s and is 30 long
     18th Solution for 22nd Col took     2h9m9.944309036s and is 30 long
     19th Solution for 18th Col took    2h16m56.77916977s and is 30 long
     20th Solution for 10th Col took   2h31m29.901413772s and is 30 long
     21st Solution for  6th Col took   2h32m17.004094943s and is 30 long
     22nd Solution for 13th Col took   2h43m46.517150932s and is 30 long
     23rd Solution for  7th Col took     3h47m9.75321572s and is 30 long
     24th Solution for  0th Col took   5h48m43.115539259s and is 30 long
     25th Solution for  3rd Col took   5h49m27.725350314s and is 30 long
     26th Solution for  1st Col took    6h15m1.320787215s and is 30 long
     27th Solution for  2nd Col took   6h17m13.277885383s and is 30 long