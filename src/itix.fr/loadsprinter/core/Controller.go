package core

import "io"
import "log"
import "fmt"
import "sync"
import "time"
import "encoding/csv"

type Controller struct {
  vuserLog *log.Logger
  controllerLog *log.Logger
  users []*VirtualUser
  wg sync.WaitGroup
  csv *csv.Writer
  stepStats chan StepIteration
  scenarioStats chan ScenarioIteration
}

type VirtualUserFactory interface {
  CreateVirtualUser() *VirtualUser
}

type StepIteration struct {
  vuser string
  scenario string
  step string
  success bool
  elapsed time.Duration
}

func (s StepIteration) String() string {
  return fmt.Sprintf("step: %v > %v > %v: success = %v, elapsed = %v", s.vuser, s.scenario, s.step, s.success, s.elapsed)
}

type ScenarioIteration struct {
  vuser string
  scenario string
  success bool
  elapsed time.Duration
}

func (s ScenarioIteration) String() string {
  return fmt.Sprintf("scenario: %v > %v: success = %v, elapsed = %v", s.vuser, s.scenario, s.success, s.elapsed)
}

func NewController(vuserLogFile io.Writer, controllerLogFile io.Writer, csvfile io.Writer) *Controller {
  c := Controller{}
  c.vuserLog = log.New(vuserLogFile, "", log.Lshortfile | log.LstdFlags)
  c.controllerLog = log.New(controllerLogFile, "", log.Lshortfile | log.LstdFlags)
  c.csv = csv.NewWriter(csvfile)
  c.stepStats = make(chan StepIteration)
  c.scenarioStats = make(chan ScenarioIteration)
  return &c;
}

func (c *Controller) StartWith(n int, f VirtualUserFactory) (wg *sync.WaitGroup, err error) {
  c.controllerLog.Printf("--> Controller::StartWith(%v, %v)", n, f)
  c.wg.Add(n) // Initialize the WaitGroup with the number of routines to create

  // Create a channel that will be used to start all VirtualUsers together
  start := make(chan int)

  for i := 0; i < n; i++ {
    log := *c.vuserLog
    vu := f.CreateVirtualUser()
    log.SetPrefix(fmt.Sprintf("%v: ", vu.name))
    c.users = append(c.users, vu)
    vu.Init(&log)
    go vu.Run(&log, &c.wg, start, c.stepStats, c.scenarioStats)
    c.controllerLog.Printf("Added one more VirtualUser with name = %v", vu.name)
  }

  // Start collecting results as soon as virtual users are started
  go func() { <- start; c.GatherResults() }()

  c.controllerLog.Println("Ready ? Set !")
  time.Sleep(1 * time.Second)
  // Start all VirtualUsers together
  close(start)
  c.controllerLog.Println("Go !")

  c.controllerLog.Println("<-- Controller::StartWith()")
  return &c.wg, nil
}

func (c *Controller) GatherResults() {
  c.controllerLog.Println("--> Controller::GatherResults()")

  // Gather results every second
  timer := make(chan int)
  go func() {
    for {
      time.Sleep(1 * time.Second)
      timer <- 0
    }
  }()

  for {
    c.controllerLog.Println("Waiting for message")
    select {
    case <- timer:
      c.controllerLog.Println("Time's up !")
    case stepStat := <- c.stepStats:
      c.controllerLog.Println(stepStat)
    case scenarioStat := <- c.scenarioStats:
      c.controllerLog.Println(scenarioStat)
    }
    c.controllerLog.Println("Received a message")
  }

  c.controllerLog.Println("<-- Controller::GatherResults()")
}
