package golis_test

import (
	"fmt"
	"testing"

	"github.com/Konstantin8105/golis"
	"gonum.org/v1/gonum/mat"
)

func TestLis(t *testing.T) {
	A := mat.NewDense(2, 2, []float64{
		1.0, 2.0,
		4.0, 1.0,
	})
	b := mat.NewDense(2, 1, []float64{
		4.0,
		9.0,
	})

	s, r, o, e := golis.Lsolve(A, b, 0)
	fmt.Println(s, r, o, e)
}
