package shared

import "fmt"

type NotFoundError struct {
	// The thing that couldn't be found
	Thing string
	// The source of the error (the function that returned it and what file that function is in)
	Source Source
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%v could not find %v", e.Source, e.Thing)
}
