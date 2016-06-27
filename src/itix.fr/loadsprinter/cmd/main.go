package main

import "fmt"
import "os"
import "time"
import "itix.fr/loadsprinter/core"
import "itix.fr/loadsprinter/steps"

type MyVirtualUserFactory struct {
  count int
  scenario *core.Scenario
}

func (factory *MyVirtualUserFactory) CreateVirtualUser() *core.VirtualUser {
  name := fmt.Sprintf("vu-%03d", factory.count)
  factory.count++
  return core.NewVirtualUser(name, factory.scenario)
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
  /*vulogwr, err := os.Create("/tmp/vu.log")
  check(err)
  defer vulogwr.Close()*/
  csvwr, err := os.Create("/tmp/vu.csv")
  check(err)
  defer csvwr.Close()

  s1 := steps.NewWaitStep(300 * time.Millisecond)
  s2 := steps.NewLogStep("Hello World !")
  s3 := steps.NewFailStep("BLAAAAAAH")
  step1 := core.NewStep("wait_3s", true, s1)
  step2 := core.NewStep("log_hello", true, s2)
  step3 := core.NewStep("fail", false, s3)
  steps := []core.Step{ *step1, *step2, *step3 }
  scenario := core.NewScenario(steps, "test")

  f := &MyVirtualUserFactory{0, scenario}

  c := core.NewController(os.Stdout, os.Stdout, csvwr)
  wg, err := c.StartWith(5, f)
  wg.Wait()
}
