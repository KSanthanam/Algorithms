package com.algorithms.nqueens
import scala.collection.mutable.HashMap
import scala.collection.mutable.MutableList
import java.util.Calendar
import java.util.concurrent.TimeUnit

object NQueens  {
    val cols: HashMap[Int,Solution] = HashMap()
    def GetSolutions(N: Integer) = {
        val start = Calendar.getInstance
        val anchors = (0 to N-1).toVector.flatMap(r => (0 to N-1).map(c => (r,c))).par
        val solutions = anchors.map(processAnchor(N))
        val result = solutions.filter(_.length == N).toList
        val now = Calendar.getInstance
        val duration = now.getTimeInMillis - start.getTimeInMillis        
        println(duration, "milliseconds")
        result
    }
    def processAnchor(size: Int)(anchor: (Int, Int)) = {
        val col = anchor._2
        val soln = cols.get(col) match {
            case Some(sol) => sol
            case None => Solution(size,col,0,HashMap[(Int,Int),Boolean](),MutableList[(Int,Int)]())
        }
        val traversed = soln.progressAnchor(anchor)
        val solution = if ((anchor._1 == (size-1)) && (traversed.length == size)) traversed.toList else MutableList[(Int,Int)]().toList
        solution
    }
}

case class Solution(size: Int, col: Int,var processed: Int, visited: HashMap[(Int,Int),Boolean], var solution: MutableList[(Int,Int)]) {
    def sameRow(left:(Int,Int),right: (Int,Int)) = left._1 == right._1
    def sameCol(left:(Int,Int),right: (Int,Int)) = left._2 == right._2
    def diagnalForward(left:(Int,Int),right: (Int,Int)) = left._1 - left._2 == right._1 - right._2
    def diagnalBackward(left:(Int,Int),right: (Int,Int)) = left._1 + left._2 == right._1 + right._2
    def inPath(left:(Int,Int),right: (Int,Int)) = sameRow(left,right) || 
                                                  sameCol(left,right) ||
                                                  diagnalForward(left,right) ||
                                                  diagnalBackward(left,right)                                                   
    def isInSolnPath(pos: (Int,Int))(solution: MutableList[(Int, Int)]) = solution.map(queen => inPath(pos,queen)).fold(false)(_ || _)
    def isInPath(pos: (Int,Int)) = isInSolnPath(pos)(solution)
    def getColsForSize(r: Int)(size: Int) = (0 to size - 1).toList.map(c => (r,c))
    def getCols(r: Int) = getColsForSize(r)(size)
    def getPosForSize(r: Int, c: Int)(size: Int) = (c to size - 1).toList.map(c => (r,c))
    def getPos(r: Int, c: Int) = getPosForSize(r,c)(size)
    def setVisited(pos: (Int, Int)) = visited.put(pos,true)
    def notVisited(pos: (Int,Int)) = !visited.getOrElse(pos,false)
    def alreadyVisited(pos: (Int, Int)) = visited.getOrElse(pos,false)
    def isInPathOrVisited(pos: (Int,Int)) = isInPath(pos) || alreadyVisited(pos)
    def crawlRow(r: Int) = {
        val nextCol = getCols(r).takeWhile(isInPathOrVisited).length
        val nextAction = if (nextCol < size) {
            val nextCell = (r, nextCol)
            visited(nextCell) = true
            solution += nextCell
            processed = r
            Action.NextRow
            } else {
                Action.PopAPiece
            }
        nextAction
    }

    def progressAnchor(anchor: (Int, Int))  = {
        def stepForward() : Action.Value = {
            if (solution.length == 0) solution += ((0,anchor._2))
            val result = if (processed >= anchor._1) {
                Action.Traversed  
            } else {
                // (solution.length to size-1).toList.foreach(clearRow)
                crawlRow(solution.length)
            }
            result
        }
        def stepBackward() : Action.Value = {
            def clearRow(r: Int) = (0 to size-1).toList.foreach(col => visited.remove((r,col)))
            val result = if ((solution.length <= 1) || (solution.length >= size)) {
                Action.Traversed
            } else {
                val lastPiece = solution.takeRight(1).get(0).get
                solution = solution.dropRight(1)
                val col = lastPiece._2
                val row = solution.length
                (col+1 to size-1).toList.foreach(c => visited.remove((row,c)))
                (row+1 to size-1).toList.flatMap(r => (0 to size-1).toList.map(c => (r,c))).foreach(pos => visited.remove(pos))
                crawlRow(solution.length)
            }
            result
        }
        var next = Action.NextRow
        while (next != Action.Traversed) {
            // next match {
            //     case Action.NextRow => {
            //         println(anchor,"NextRow", processed, solution, visited)
            //     }
            //     case Action.PopAPiece => {
            //         println(anchor,"PopAPiece", processed, solution, visited)
            //     }
            // }
            next = next match {
                case Action.NextRow => stepForward()
                case Action.PopAPiece => stepBackward()
            }
        }
        solution
    }
}

// val soln = new Solution(4,0,0,HashMap[(Int,Int),Boolean](),MutableList[(Int, Int)]())
// val anchor = (2,3)
// soln.progressAnchor(anchor)
object Action extends Enumeration {
    val NoAction, NextRow, PopAPiece, Traversed = Value
}

