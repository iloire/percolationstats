package main

import(
  "github.com/iloire/percolation"
  "math/rand"
  "time"
  "fmt"
  "flag"
  "runtime"
)

var concurrency = 8

var nFlag = flag.Int("n", 200, "size")
var tFlag = flag.Int("t", 10, "repetitions")
var tCPUSFlag = flag.Int("cpus", 1, "number of CPUs")

var semaphore = make(chan int, concurrency)

func newRndCell(r *rand.Rand, n int) percolation.Cell {
  i:= uint(r.Intn(n))
  j:= uint(r.Intn(n))
  return percolation.Cell{i, j}
}

func calcPercolation(n int, c chan float32){
  semaphore <- 1
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  p:= new(percolation.Percolation)
  p.Initialize(uint(n))

  count:=0
  for !p.Percolates() {
    cell:=newRndCell(r, n)
    if !p.IsOpen(cell){
      p.Open(cell)
      count++
    }
  }
  c <- float32(count) / float32(n*n)
  <-semaphore
}

func main(){
  
  start := time.Now()

  flag.Parse()

  runtime.GOMAXPROCS(*tCPUSFlag)

  c:= make(chan float32)
  for i:=0; i<*tFlag; i++ {
    go calcPercolation(*nFlag, c)
  }

  var sum float32
  for i:=0; i<*tFlag; i++ {
    sum = sum + <-c
  }
  
  fmt.Printf("\n mean: %v", sum / float32(*tFlag))
  fmt.Printf("\n took: %v", time.Since(start))

  //TODO: stddev, confidence interval, etc.
}