package com.algorithms
import com.algorithms.nqueens._
import com.algorithms.knapsack._

object Algorithms extends App {
    def getNQueens(N: Integer) : List[List[(Int,Int)]] = NQueens.GetSolutions(N)
    def getKnapsack(pieces: Array[(Int,Int)], max: Int) = Knapsack.GetKnapsack(pieces, max)
}