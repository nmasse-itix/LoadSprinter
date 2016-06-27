package core

import "log"
import "time"

type StepImpl interface {
  Do(log *log.Logger) error
}

type Step struct {
  name string
  required bool
  impl StepImpl
}

func NewStep(name string, required bool, impl StepImpl) *Step {
  return &Step{ name: name, required: required, impl: impl }
}

func (s *Step) Do(log *log.Logger, vuser *VirtualUser, scenario *Scenario, stepsToController chan<- StepIteration) error {
  log.Printf("--> Step::Do(%v)", s.name)

  // Prepare the statistics structure
  var stat StepIteration
  stat.success = true
  stat.vuser = vuser.name
  stat.scenario = scenario.name
  stat.step = s.name

  // Measure the beginning of this step
  start := time.Now()

  err := s.impl.Do(log)
  if (err != nil) {
    log.Printf("Step %v ended with error %v", s.name, err)
    stat.success = false;
  }

  // Compute the elapsed time since the beginning of this step
  elapsed := time.Since(start)
  stat.elapsed = elapsed

  // Send it to the controller
  stepsToController <- stat

  log.Printf("<-- Step::Do(%v)", s.name)
  return err;
}
