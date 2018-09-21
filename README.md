# Algorithms 
## NQueens
Given a Chessboard of size N find possible Solutions where N queens are placed on the board

### Time complexity
The **Time Complexity** by Brute force logic is N^N as one needs to place the queen in (a,b) and check N-1 places to validate.

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
There are few solutions with inefficient timecomplexity. The worst being Q(N^N)

The solution implemented is the most efficient possible in Golang using back propagation, channel and go routines.

#### Synopsis of Algorithm

The Algorithm generates N*N anchor points and submits to the anchor channel.

$$
i = 1..N, j =1..N
(\R_i,C_j)
$$

Each anchor is read from the channel and submited to go routine that  traverses the board using the PieceStack dedicated for the column of the given anchor.

The Traverse Logic for a given anchor.
when a anchor  $$ (\R_i,C_j) $$ is receied, a go routine is kick started with the aim to place Queens up to the row i. If the row i can be reached with placements it back tracks by forwarding the column position. If it exhausts all columns for rows before i, the go routine stops trying. 
so when the the last row Nth row anchor is received and the row can be reached, it is considered to be a solution as there will be N queens in place.

So the logic is as follows:
<pre>
anchor go routine for $$ (\R_i,C_j) $$
if the i == 0 insert the anchor into the piece stack for j.

else

r = i + 1
colStart = 0
for r <= i do
c = colStart
   placed = false
   for c < N do
       if (r,c) has not been visited then
         visited[(r,c)] = true
         if can Queen be placed in (r,c) then
             placed = true
         end
       end
       c++
   done
   if placed then
      r++
      cStart = 0
   else
      if r <= i
        pop last cell in Piece Stack
        r = last cell row
        c = last cell col + 1
      else
        break
      end
   end
   if r < j then
     for r is r to j do
       for c is 0 to N do
          visited[(r,c)] = true
       done
     done
   end
   if r == N-1 then
      send the cells in PieceStack to Solutions 
   end
done
</pre>


[Markdown Syntax](https://stackedit.io/app#)