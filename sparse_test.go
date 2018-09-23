package golis_test

import (
	"testing"

	"github.com/Konstantin8105/golis"
	"gonum.org/v1/gonum/mat"
)

func TestSparseSquareMatrix(t *testing.T) {
	a := mat.NewDense(3, 3, []float64{
		8, 1, 6,
		3, 5, 7,
		4, 0, 2,
	})

	t.Run("Add", func(t *testing.T) {
		s := golis.NewSparseSquareMatrix(3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				s.Add(i, j, a.At(i, j)/2.0)
				s.Add(i, j, a.At(i, j)/2.0)
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%#v\n%#v", s, a)
		}
	})

	t.Run("Add reverse", func(t *testing.T) {
		s := golis.NewSparseSquareMatrix(3)
		for i := 2; i >= 0; i-- {
			for j := 2; j >= 0; j-- {
				s.Add(i, j, a.At(i, j)/2.0)
				s.Add(i, j, a.At(i, j)/2.0)
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%#v\n%#v", s, a)
		}
	})

	t.Run("Add random", func(t *testing.T) {
		s := golis.NewSparseSquareMatrix(3)
		for i := 2; i >= 0; i-- {
			for j := 0; j < 3; j++ {
				s.Add(i, j, a.At(i, j)/2.0)
				s.Add(i, j, a.At(i, j)/2.0)
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%#v\n%#v", s, a)
		}
	})

	t.Run("Set", func(t *testing.T) {
		s := golis.NewSparseSquareMatrix(3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				s.Set(i, j, a.At(i, j))
				s.Set(i, j, a.At(i, j))
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%#v\n%#v", s, a)
		}
	})

	t.Run("Set reverse", func(t *testing.T) {
		s := golis.NewSparseSquareMatrix(3)
		for i := 2; i >= 0; i-- {
			for j := 2; j >= 0; j-- {
				s.Set(i, j, a.At(i, j))
				s.Set(i, j, a.At(i, j))
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%#v\n%#v", s, a)
		}
	})

	t.Run("Set random", func(t *testing.T) {
		s := golis.NewSparseSquareMatrix(3)
		for i := 2; i >= 0; i-- {
			for j := 0; j < 3; j++ {
				s.Set(i, j, a.At(i, j))
				s.Set(i, j, a.At(i, j))
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%#v\n%#v", s, a)
		}
	})

	t.Run("Matrix with zero values", func(t *testing.T) {
		a := mat.NewDense(3, 3, make([]float64, 9))
		s := golis.NewSparseSquareMatrix(3)
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%#v\n%#v", s, a)
		}
	})

	t.Run("Sparse matrix", func(t *testing.T) {
		a := mat.NewDense(3, 3, make([]float64, 9))
		a.Set(1, 1, 42)
		s := golis.NewSparseSquareMatrix(3)
		s.Set(1, 1, 42)
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%#v\n%#v", s, a)
		}
	})
}

func isSame(s *golis.SparseSquareMatrix, a *mat.Dense) bool {
	r, c := s.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if s.At(i, j) != a.At(i, j) {
				return false
			}
		}
	}
	return true
}
