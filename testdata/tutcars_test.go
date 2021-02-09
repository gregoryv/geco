package testdata

import "testing"

type tutCar struct {
	*testing.T
	*Car
}

func (me *Car) tu(t *testing.T) *tutCar {
	return &tutCar{T: t, Car: me}
}

func (me *tutCar) shouldStart(key int) {
	err := me.Start(key)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *tutCar) mustStart(key int) {
	err := me.Start(key)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}
