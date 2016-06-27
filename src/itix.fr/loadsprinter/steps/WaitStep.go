package steps;

import "time"
import "log"

type WaitStep struct {
  duration time.Duration
}

func NewWaitStep(duration time.Duration) *WaitStep {
  return &WaitStep{ duration: duration };
}

func (ws *WaitStep) Do(log *log.Logger) error {
  log.Printf("Sleeping during %v", ws.duration)
  time.Sleep(ws.duration)
  log.Println("Woken up !")
  return nil
}
