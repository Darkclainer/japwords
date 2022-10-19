package jisho

import "fmt"

type LemmaBatchError struct {
	Errs []error
}

func (e *LemmaBatchError) Error() string {
	return fmt.Sprintf("%d lemma parsing failed", len(e.Errs))
}

type LemmaError struct {
	ID  int
	Err error
}

func (e *LemmaError) Error() string {
	return fmt.Sprintf("lemma parsing %d failed: %s", e.ID, e.Err)
}

func (e *LemmaError) Unwrap() error {
	return e.Err
}
