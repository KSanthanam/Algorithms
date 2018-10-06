import org.scalatest._
import com.algorithms._

class AlgorithmsSpec extends FunSuite with DiagrammedAssertions {
    test("NQueens for 4 should be as below") {  
        val N = 4
        val result = List[List[(Int,Int)]](List[(Int,Int)]((0,1), (1,3), (2,0), (3,2)), List[(Int,Int)]((0,2), (1,0), (2,3), (3,1)))
        assert(Algorithms.getNQueens(N) == result)
    }
    test("NQueens for 10 should be as below") {  
        val N = 10
        val result = List(List((0,0), (1,2), (2,5), (3,7), (4,9), (5,4), (6,8), (7,1), (8,3), (9,6)), List((0,1), (1,3), (2,5), (3,7), (4,9), (5,0), (6,2), (7,4), (8,6), (9,8)), List((0,2), (1,0), (2,5), (3,8), (4,4), (5,9), (6,7), (7,3), (8,1), (9,6)), List((0,3), (1,0), (2,4), (3,7), (4,9), (5,2), (6,6), (7,8), (8,1), (9,5)), List((0,4), (1,0), (2,3), (3,8), (4,6), (5,1), (6,9), (7,2), (8,5), (9,7)), List((0,5), (1,0), (2,2), (3,9), (4,7), (5,1), (6,3), (7,8), (8,6), (9,4)), List((0,6), (1,0), (2,2), (3,5), (4,7), (5,9), (6,3), (7,8), (8,4), (9,1)), List((0,7), (1,0), (2,2), (3,5), (4,8), (5,6), (6,9), (7,3), (8,1), (9,4)), List((0,8), (1,0), (2,2), (3,7), (4,5), (5,1), (6,9), (7,4), (8,6), (9,3)), List((0,9), (1,0), (2,3), (3,5), (4,2), (5,8), (6,1), (7,7), (8,4), (9,6)))
        assert(Algorithms.getNQueens(N) == result)
    }
    /*
KnapSack 0/1 Problem Example
m = 8
n = 4
P = {1,2,5,6}
W = {2,3,4,5}
*/
    test("Knapsack for  should be as below") {  
        val m = 8
        val pieces = Array[(Int,Int)]((2,1),(3,2),(4,5),(5,6))
        val result = List((4,5),(3,2))
        assert(Algorithms.getKnapsack(pieces,m) == result)
    }
}