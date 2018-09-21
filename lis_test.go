package golis_test

import (
	"fmt"
	"testing"

	"github.com/Konstantin8105/golis"
	"gonum.org/v1/gonum/mat"
)

func TestLsolve(t *testing.T) {
	A := mat.NewDense(2, 2, []float64{
		1.0, 2.0,
		4.0, 1.0,
	})
	b := mat.NewDense(2, 1, []float64{
		4.0,
		9.0,
	})

	s, r, o, err := golis.Lsolve(A, b, 0)
	fmt.Println(s, r, o, err)
}

func TestLsolveSingular(t *testing.T) {
	A := mat.NewDense(2, 2, []float64{
		1.0, 2.0,
		1.e-30, 1e-30,
	})
	b := mat.NewDense(2, 1, []float64{
		3.0,
		2.0e-30,
	})

	s, r, o, err := golis.Lsolve(A, b, 0)
	fmt.Println(s, r, o, err)
	if err == nil {
		t.Fatalf("Haven`t error : %v", err)
	}
}

func TestLsolveFail(t *testing.T) {
	A := mat.NewDense(2, 2, []float64{
		1.0, 2.0,
		1.0, 2.0,
	})
	b := mat.NewDense(2, 1, []float64{
		4.0,
		4.0,
	})

	s, r, o, err := golis.Lsolve(A, b, 0)
	fmt.Println(s, r, o, err)
	if err == nil {
		t.Fatalf("Haven`t error : %v", err)
	}
}
