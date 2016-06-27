package main

import "fmt"
import "time"
//import "sync"

type Thing struct {
  a int
  b string
}

func (t Thing) String() string {
  return fmt.Sprintf("%v-%v", t.b, t.a)
}

func Do(a int, b string, ch chan Thing) {
  var t Thing
  t.a = a
  t.b = b
  ch <- t
}

func Doit(id string, chan1 chan Thing, chan2 chan Thing) {
  for {
    var t Thing
    t.a = 666
    t.b = fmt.Sprintf("%v-evil", id)

    Do(1, fmt.Sprintf("%v-one", id), chan1)
    time.Sleep(300 * time.Millisecond)
    Do(2, fmt.Sprintf("%v-two", id), chan1)
    Do(3, fmt.Sprintf("%v-three", id), chan1)
    time.Sleep(50 * time.Millisecond)
    Do(4, fmt.Sprintf("%v-four", id), chan1)

    chan2 <- t
  }
}

func main() {
  chan1 := make(chan Thing)
  chan2 := make(chan Thing)

  go Doit("A", chan1, chan2);
  go Doit("B", chan1, chan2);
  go Doit("C", chan1, chan2);
  go Doit("D", chan1, chan2);
  go Doit("E", chan1, chan2);

  for {
    time.Sleep(1*time.Second)
    select {
    case x := <- chan1:
      fmt.Printf("1: %v\n", x)
    case x := <- chan2:
      fmt.Printf("2: %v\n", x)
    }
  }
}
