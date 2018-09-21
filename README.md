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

The Algorithm generates N*N anchor points 
$$\(r_i,c_j) 
where i = 1..N and j = 1..N. 

[Markdown Syntax](https://stackedit.io/app#)