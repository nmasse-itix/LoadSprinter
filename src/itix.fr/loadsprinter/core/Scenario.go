package core

import "log"
import "time"

type Scenario struct {
  steps []Step
  name string
}

func NewScenario(steps []Step, name string) *Scenario{
  return &Scenario{ steps: steps, name: name };
}

func (s *Scenario) Do(log *log.Logger, vuser *VirtualUser, stepsToController chan<- StepIteration, scenarioToController chan<- ScenarioIteration) error {
  log.Println("--> Scenario::Do()")

  // Prepare the statistics structure
  var stat ScenarioIteration
  stat.success = true;
  stat.vuser = vuser.name
  stat.scenario = s.name

  // Measure the beginning of the scenario
  start := time.Now()

  // Run the scenario
  var err error = nil
  for _, step := range s.steps {
    err = step.Do(log, vuser, s, stepsToController)
    if (err != nil) {
      if (step.required) {
        log.Printf("Mandatory step %v ended with error %v, stopping scenario !", step.name, err)
        stat.success = false;
        break
      } else {
        log.Printf("Optional step %v ended with error %v, continuing...", step.name, err)
      }
    }
  }

  // Compute the elapsed time since the beginning of this scenario
  elapsed := time.Since(start)
  stat.elapsed = elapsed

  // Send it to the controller
  scenarioToController <- stat

  log.Printf("<-- Scenario::Do() : err = %v", err)
  return err
}
