package steps;

import "log"
import "errors"

type FailStep struct {
  message string
}

func NewFailStep(message string) *FailStep {
  return &FailStep{ message: message };
}

func (fs *FailStep) Do(log *log.Logger) error {
  log.Printf("error: %v", fs.message)
  return errors.New(fs.message);
}
