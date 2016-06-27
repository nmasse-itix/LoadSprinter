package core;

import "log"
import "sync"

type VirtualUser struct {
  name string
  scenario *Scenario
}

func NewVirtualUser(name string, scenario *Scenario) *VirtualUser {
  return &VirtualUser{ name: name, scenario: scenario }
}

func (vu *VirtualUser) Init(log *log.Logger) error {
  log.Println("--> VirtualUser::Init()")
  // TODO
  log.Println("<-- VirtualUser::Init()")
  return nil
}

func (vu *VirtualUser) Run(log *log.Logger, wg *sync.WaitGroup, start <-chan int, stepsToController chan<- StepIteration, scenarioToController chan<- ScenarioIteration) error {
  log.Println("--> VirtualUser::Run()")

  // Make sure we notify we are done when this method is finished
  defer wg.Done()

  // Wait for the signal to start
  log.Println("Waiting for the signal to start")
  <- start
  log.Println("GOOOOOO !")

  for {
    err := vu.scenario.Do(log, vu, stepsToController, scenarioToController);
    if (err != nil) {
      log.Printf("Scenario ended with error %v, ", err)

    }
  }
  log.Println("<-- VirtualUser::Run()")
  return nil
}
