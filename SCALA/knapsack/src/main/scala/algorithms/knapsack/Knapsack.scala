package com.algorithms.knapsack
import scala.util.Sorting
import scala.collection.mutable.MutableList

object Knapsack {
    def GetKnapsack(pieces: Array[(Int,Int)], maxWeight: Int) = {
        val objects = pieces.map(piece => Piece(piece._1,piece._2))
        Sorting.quickSort(objects)(WeightOrdering)
        val ordered = Array[Piece](Piece(0,0)) ++ objects
        var Values = (0 to ordered.length-1).toArray.map(row => (0 to maxWeight).toArray.map(col => 0))
        var max = MaxPosition(0,0,0)
        def placeValue(r: Int, w: Int) : Int = {
            val o = ordered(r)
            val preProfit = if ((w-o.weight) >= 0) Values(r-1)(w-o.weight) else 0
            val lastProfit = Values(r-1)(w)
            val result = if (w > o.weight) {
                val currentProfit = o.profit + preProfit
                if (lastProfit > currentProfit) lastProfit else currentProfit
            } else {
                lastProfit
            }
            result
        }
        def processRow(i: Int, w: Int) = {
            val result = if ((i == 0) || (w == 0)) {
            0
            } else {
                placeValue(i,w)
            }
            if (result > max.profit) {
                max = MaxPosition(i,w,result)
            }
            Values(i)(w) = result
        }
        (0 to ordered.length-1).toList.foreach(i => (0 to maxWeight).toList.foreach(w => processRow(i,w)) )

        val selected = MutableList[(Int,Int)]()
	    var profitLeft = max.profit
        var maxi = max.row
        var	maxw = max.weight
        while ((profitLeft > 0) && (maxi > 0) && (maxw > 0)) {
            if (Values(maxi)(maxw) > Values(maxi-1)(maxw)) {
                selected += ((ordered(maxi).weight,ordered(maxi).profit))
                profitLeft -= ordered(maxi).profit
                maxw -= ordered(maxi).weight
            }
            maxi -= 1
        }
        selected
    }
}

case class Piece(weight: Int, profit: Int) 
case class MaxPosition(row: Int, weight: Int, profit: Int)
object WeightOrdering extends Ordering[Piece] {
    def compare(a: Piece, b: Piece) = a.weight compare b.weight
}
