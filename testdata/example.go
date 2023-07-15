package testdata

//go:generate mkgetters -t Car,Boat -w getters.go
type Car struct {
	Name string

	model string
	make  int
}

type Boat struct {
	color int
}
