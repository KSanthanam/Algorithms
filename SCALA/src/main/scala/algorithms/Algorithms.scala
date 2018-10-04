package com.algorithms
import com.algorithms.nqueens._

object Algorithms extends App {
    def getNQueens(N: Integer) : List[List[(Int,Int)]] = NQueens.GetSolutions(N)
}