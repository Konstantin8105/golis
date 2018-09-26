package golis_test

import (
	"fmt"
	"testing"

	"github.com/Konstantin8105/golis"
	"gonum.org/v1/gonum/mat"
)

func TestSparseMatrix(t *testing.T) {
	a := mat.NewDense(3, 3, []float64{
		8, 1, 6,
		3, 5, 7,
		4, 0, 2,
	})

	t.Run("Add", func(t *testing.T) {
		s := golis.NewSparseMatrix(3, 3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				s.Add(i, j, a.At(i, j)/2.0)
				s.Add(i, j, a.At(i, j)/2.0)
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Add up-down-up", func(t *testing.T) {
		s := golis.NewSparseMatrix(3, 3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				s.Add(i, j, a.At(i, j)/2.0)
				s.Add(i, j, a.At(i, j)/2.0)
				s.Add(i, j, a.At(i, j)/2.0)
			}
		}
		for i := 2; i >= 0; i-- {
			for j := 2; j >= 0; j-- {
				s.Add(i, j, -a.At(i, j)/2.0)
				s.Add(i, j, 0.0)
				s.Add(i, j, -a.At(i, j)/2.0)
			}
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				s.Add(i, j, -a.At(i, j)/2.0)
				s.Add(i, j, a.At(i, j)/2.0)
				s.Add(i, j, a.At(i, j)/2.0)
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Add reverse", func(t *testing.T) {
		s := golis.NewSparseMatrix(3, 3)
		for i := 2; i >= 0; i-- {
			for j := 2; j >= 0; j-- {
				s.Add(i, j, a.At(i, j)/2.0)
				s.Add(i, j, a.At(i, j)/2.0)
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Add random", func(t *testing.T) {
		s := golis.NewSparseMatrix(3, 3)
		for i := 2; i >= 0; i-- {
			for j := 0; j < 3; j++ {
				s.Add(i, j, a.At(i, j)/2.0)
				s.Add(i, j, a.At(i, j)/2.0)
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Set", func(t *testing.T) {
		s := golis.NewSparseMatrix(3, 3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				s.Set(i, j, a.At(i, j))
				s.Set(i, j, a.At(i, j))
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Set reverse", func(t *testing.T) {
		s := golis.NewSparseMatrix(3, 3)
		for i := 2; i >= 0; i-- {
			for j := 2; j >= 0; j-- {
				s.Set(i, j, a.At(i, j))
				s.Set(i, j, a.At(i, j))
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Set random", func(t *testing.T) {
		s := golis.NewSparseMatrix(3, 3)
		for i := 2; i >= 0; i-- {
			for j := 0; j < 3; j++ {
				s.Set(i, j, a.At(i, j))
				s.Set(i, j, a.At(i, j))
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Matrix with zero values", func(t *testing.T) {
		a := mat.NewDense(3, 3, make([]float64, 9))
		s := golis.NewSparseMatrix(3, 3)
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Sparse matrix", func(t *testing.T) {
		a := mat.NewDense(3, 3, make([]float64, 9))
		a.Set(1, 1, 42)
		s := golis.NewSparseMatrix(3, 3)
		s.Set(1, 1, 42)
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Rectange horizontal matrix", func(t *testing.T) {
		a := mat.NewDense(3, 2, []float64{
			8, 1, 6,
			3, 5, 7,
		})
		s := golis.NewSparseMatrix(3, 2)
		for i := 0; i < 3; i++ {
			for j := 0; j < 2; j++ {
				s.Set(i, j, a.At(i, j))
				s.Set(i, j, a.At(i, j))
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Rectange vertical matrix", func(t *testing.T) {
		a := mat.NewDense(2, 3, []float64{
			8, 1,
			3, 5,
			2, 6,
		})
		s := golis.NewSparseMatrix(2, 3)
		for i := 0; i < 2; i++ {
			for j := 0; j < 3; j++ {
				s.Set(i, j, a.At(i, j))
				s.Set(i, j, a.At(i, j))
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})
}

func isSame(s *golis.SparseMatrix, a *mat.Dense) bool {
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

func TestParse(t *testing.T) {
	b := []byte(`%%MatrixMarket vector coordinate real general
3
1  -5.49999999999999822364e+00
2   2.49999999999999955591e+00
3   4.99999999999999911182e+00`)

	s, err := golis.ParseSparseMatrix(b)
	if err != nil {
		t.Fatalf("Cannot parse : %v", err)
	}

	a := mat.NewDense(3, 1, []float64{
		-5.49999999999999822364e+00,
		2.49999999999999955591e+00,
		4.99999999999999911182e+00,
	})
	if !isSame(s, a) {
		t.Fatalf("Value is not same:\n%#v\n%#v", s, a)
	}
}

func ExampleString() {
	s := golis.NewSparseMatrix(3, 2)
	for i := 0; i < 3; i++ {
		for j := 0; j < 2; j++ {
			s.Set(i, j, float64(i+j*5))
		}
	}
	fmt.Printf("%s", s)

	// Output:
	// Amount of rows    :     3
	// Amount of columns :     2
	// row    column                value
	// 1      0      1.000000000000000e+00
	// 2      0      2.000000000000000e+00
	// 0      1      5.000000000000000e+00
	// 1      1      6.000000000000000e+00
	// 2      1      7.000000000000000e+00
}

func BenchmarkAt(b *testing.B) {
	sizes := []int{10, 20, 40, 80}
	for i := range sizes {
		b.Run(fmt.Sprintf("%d", sizes[i]), func(b *testing.B) {
			b.StopTimer()
			s := golis.NewSparseMatrix(sizes[i], sizes[i])
			r, c := s.Dims()
			for i := 0; i < r; i++ {
				for j := 0; j < c; j++ {
					s.Set(i, j, float64(i+j*5))
				}
			}
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < r; j++ {
					_ = s.At(j, j)
				}
			}
		})
	}
}
