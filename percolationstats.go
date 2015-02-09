package main

import(
  "github.com/iloire/percolation"
  "math/rand"
  "time"
  "fmt"
  "flag"
  "runtime"
)

var nFlag = flag.Int("n", 10, "size")
var tFlag = flag.Int("t", 100, "repetitions")
var tCPUS = flag.Int("cpus", 4, "number of CPUs")

func newRndCell(n int) percolation.Cell {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  i:= uint(r.Intn(n))
  j:= uint(r.Intn(n))
  return percolation.Cell{i, j}
}

func calcPercolation(n int, c chan float32){
  count:=0
  p:= new(percolation.Percolation)
  p.Initialize(uint(n))

  for !p.Percolates() {
    cell:=newRndCell(n)
    if !p.IsOpen(cell){
      p.Open(cell)
      count++
    }
  }
  c <- float32(count) / float32(n*n)
}

func main(){

  start := time.Now()

  flag.Parse()

  var sum float32
  var n = *nFlag
  var t = *tFlag
  var cpus = *tCPUS

  runtime.GOMAXPROCS(cpus)

  c:= make(chan float32)
  for i:=0; i<t; i++ {
    go calcPercolation(n, c)
  }

  for i:=0; i<t; i++ {
    sum = sum + <-c
  }
  
  fmt.Printf("\n mean: %v", sum / float32(t))
  fmt.Printf("\n took: %v", time.Since(start))

  //TODO: stddev, confidence interval, etc.
}