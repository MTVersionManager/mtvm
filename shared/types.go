package shared

type SuccessMsg string

type Source struct {
	File     string
	Function string
}

func (s Source) String() string {
	return s.File + ":" + s.Function
}
