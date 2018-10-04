import org.scalatest._
import com.algorithms._

class AlgorithmsSpec extends FunSuite with DiagrammedAssertions {
    test("Solution for 4 should be as below") {  
        val N = 4
        val result = List[List[(Int,Int)]](List[(Int,Int)]((0,1), (1,3), (2,0), (3,2)), List[(Int,Int)]((0,2), (1,0), (2,3), (3,1)))
        assert(Algorithms.getNQueens(N) == result)
    }
    test("Solution for 10 should be as below") {  
        val N = 10
        val result = List(List((0,0), (1,2), (2,5), (3,7), (4,9), (5,4), (6,8), (7,1), (8,3), (9,6)), List((0,1), (1,3), (2,5), (3,7), (4,9), (5,0), (6,2), (7,4), (8,6), (9,8)), List((0,2), (1,0), (2,5), (3,8), (4,4), (5,9), (6,7), (7,3), (8,1), (9,6)), List((0,3), (1,0), (2,4), (3,7), (4,9), (5,2), (6,6), (7,8), (8,1), (9,5)), List((0,4), (1,0), (2,3), (3,8), (4,6), (5,1), (6,9), (7,2), (8,5), (9,7)), List((0,5), (1,0), (2,2), (3,9), (4,7), (5,1), (6,3), (7,8), (8,6), (9,4)), List((0,6), (1,0), (2,2), (3,5), (4,7), (5,9), (6,3), (7,8), (8,4), (9,1)), List((0,7), (1,0), (2,2), (3,5), (4,8), (5,6), (6,9), (7,3), (8,1), (9,4)), List((0,8), (1,0), (2,2), (3,7), (4,5), (5,1), (6,9), (7,4), (8,6), (9,3)), List((0,9), (1,0), (2,3), (3,5), (4,2), (5,8), (6,1), (7,7), (8,4), (9,6)))
        assert(Algorithms.getNQueens(N) == result)
    }
}