ThisBuild / scalaVersion := "2.12.6"
ThisBuild / organization := "com.algorithms"


val scalaTest = "org.scalatest" %% "scalatest" % "3.0.5"

lazy val algorithms = (project in file("."))
.aggregate(nqueens,knapsack)
.dependsOn(nqueens,knapsack)
  .settings(
    name := "Algorithms",
    libraryDependencies += scalaTest % Test,
)

lazy val nqueens = (project in file ("nqueens"))
    .settings(
        name := "NQueens",
libraryDependencies += scalaTest % Test,
    )

lazy val knapsack = (project in file ("knapsack"))
    .settings(
        name := "Knapsack",
libraryDependencies += scalaTest % Test,
    )