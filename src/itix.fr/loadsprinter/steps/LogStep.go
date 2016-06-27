package steps;

import "log"

type LogStep struct {
  message string
}

func NewLogStep(message string) *LogStep {
  return &LogStep{ message: message };
}

func (ls *LogStep) Do(log *log.Logger) error {
  log.Println(ls.message)
  return nil;
}
