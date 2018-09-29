package golis_test

import (
	"testing"

	"github.com/Konstantin8105/golis"
	"gonum.org/v1/gonum/mat"
)

func TestSparseMatrixSymmetric(t *testing.T) {
	a := mat.NewDense(3, 3, []float64{
		8, 1, 6,
		1, 5, 7,
		6, 7, 2,
	})

	t.Run("Add", func(t *testing.T) {
		s := golis.NewSparseMatrixSymmetric(3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if i > j {
					continue
				}
				s.Add(i, j, a.At(i, j)/2.0)
				s.Add(i, j, a.At(i, j)/2.0)
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Add up-down-up", func(t *testing.T) {
		s := golis.NewSparseMatrixSymmetric(3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if i > j {
					continue
				}
				s.Add(i, j, a.At(i, j)/2.0)
				s.Add(i, j, a.At(i, j)/2.0)
				s.Add(i, j, a.At(i, j)/2.0)
			}
		}
		for i := 2; i >= 0; i-- {
			for j := 2; j >= 0; j-- {
				if i > j {
					continue
				}
				s.Add(i, j, -a.At(i, j)/2.0)
				s.Add(i, j, 0.0)
				s.Add(i, j, -a.At(i, j)/2.0)
			}
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if i > j {
					continue
				}
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
		s := golis.NewSparseMatrixSymmetric(3)
		for i := 2; i >= 0; i-- {
			for j := 2; j >= 0; j-- {
				if i > j {
					continue
				}
				s.Add(i, j, a.At(i, j)/2.0)
				s.Add(i, j, a.At(i, j)/2.0)
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Add random", func(t *testing.T) {
		s := golis.NewSparseMatrixSymmetric(3)
		for i := 2; i >= 0; i-- {
			for j := 0; j < 3; j++ {
				if i > j {
					continue
				}
				s.Add(i, j, a.At(i, j)/2.0)
				s.Add(i, j, a.At(i, j)/2.0)
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Set", func(t *testing.T) {
		s := golis.NewSparseMatrixSymmetric(3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if i > j {
					continue
				}
				s.SetSym(i, j, a.At(i, j))
				s.SetSym(i, j, a.At(i, j))
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Set reverse", func(t *testing.T) {
		s := golis.NewSparseMatrixSymmetric(3)
		for i := 2; i >= 0; i-- {
			for j := 2; j >= 0; j-- {
				if i > j {
					continue
				}
				s.SetSym(i, j, a.At(i, j))
				s.SetSym(i, j, a.At(i, j))
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Set random", func(t *testing.T) {
		s := golis.NewSparseMatrixSymmetric(3)
		for i := 2; i >= 0; i-- {
			for j := 0; j < 3; j++ {
				if i > j {
					continue
				}
				s.SetSym(i, j, a.At(i, j))
				s.SetSym(i, j, a.At(i, j))
			}
		}
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Matrix with zero values", func(t *testing.T) {
		a := mat.NewDense(3, 3, make([]float64, 9))
		s := golis.NewSparseMatrixSymmetric(3)
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Sparse matrix", func(t *testing.T) {
		a := mat.NewDense(3, 3, make([]float64, 9))
		a.Set(1, 1, 42)
		s := golis.NewSparseMatrixSymmetric(3)
		s.SetSym(1, 1, 42)
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})

	t.Run("Transpose", func(t *testing.T) {
		a := mat.NewDense(3, 3, []float64{
			8, 1, 6,
			1, 5, 7,
			6, 7, 2,
		})
		s := golis.NewSparseMatrixSymmetric(3)
		for i := 0; i < 2; i++ {
			for j := 0; j < 3; j++ {
				if i > j {
					continue
				}
				s.SetSym(i, j, a.At(i, j))
				s.SetSym(i, j, a.At(i, j))
			}
		}
		st := s.T()
		stt := st.T()
		if !isSame(stt, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", stt, a)
		}
	})

	// t.Run("String empty sparse matrix", func(t *testing.T) {
	// 	s := golis.NewSparseMatrixSymmetric(3)
	// 	if len(s.String()) == 0 {
	// 		t.Fatalf("String for empty sparse matrix is empty")
	// 	}
	// })

	t.Run("SetZeroForRowColumn", func(t *testing.T) {
		a := mat.NewDense(3, 3, []float64{
			8, 1, 6,
			1, 5, 7,
			6, 7, 2,
		})
		s := golis.NewSparseMatrixSymmetric(3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if i > j {
					continue
				}
				s.Add(i, j, a.At(i, j))
			}
		}
		a.Set(0, 0, 0.0)
		a.Set(1, 0, 0.0)
		a.Set(2, 0, 0.0)
		a.Set(0, 1, 0.0)
		a.Set(0, 2, 0.0)
		s.SetZeroForRowColumn(0)
		if !isSame(s, a) {
			t.Fatalf("Value is not same:\n%s\n%#v", s, a)
		}
	})
}

// func isSame(s mat.Matrix, a mat.Matrix) bool {
// 	r, c := s.Dims()
// 	for i := 0; i < r; i++ {
// 		for j := 0; j < c; j++ {
// 			if s.At(i, j) != a.At(i, j) {
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }

// func ExampleSparseMatrix() {
// 	s := golis.NewSparseMatrix(3, 2)
// 	for i := 0; i < 3; i++ {
// 		for j := 0; j < 2; j++ {
// 			s.SetSym(i, j, float64(i+j*5))
// 		}
// 	}
// 	fmt.Fprintf(os.Stdout, "%s", s)
//
// 	// Output:
// 	// Amount of rows    :     3
// 	// Amount of columns :     2
// 	// row    column                value
// 	// 1      0      1.000000000000000e+00
// 	// 2      0      2.000000000000000e+00
// 	// 0      1      5.000000000000000e+00
// 	// 1      1      6.000000000000000e+00
// 	// 2      1      7.000000000000000e+00
// }

// func BenchmarkAt(b *testing.B) {
// 	sizes := []int{10, 20, 40, 80}
// 	for is := range sizes {
// 		s := golis.NewSparseMatrix(sizes[is], sizes[is])
// 		r, c := s.Dims()
// 		for i := 0; i < r; i++ {
// 			for j := 0; j < c; j++ {
// 				if r/2-4 < i && i < r/2+4 {
// 					if c/2-4 < j && j < c/2+4 {
// 						s.SetSym(i, j, float64(i+j*5))
// 					}
// 				}
// 			}
// 		}
// 		b.Run(fmt.Sprintf("ByRow   :%d", sizes[is]), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				for j := 0; j < r; j++ {
// 					_ = s.At(r/3, j)
// 				}
// 			}
// 		})
// 		b.Run(fmt.Sprintf("ByColumn:%d", sizes[is]), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				for j := 0; j < r; j++ {
// 					_ = s.At(j, c/3)
// 				}
// 			}
// 		})
// 		b.Run(fmt.Sprintf("OneCell :%d", sizes[is]), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				_ = s.At(r/3, r/3)
// 			}
// 		})
// 	}
// }

// func TestSparseMatrixPanics(t *testing.T) {
// 	for i, tc := range []struct{ r, c int }{
// 		{0, 0},
// 		{-1, 5},
// 		{5, -1},
// 		{-1, -1},
// 	} {
// 		t.Run(fmt.Sprintf("Panic%d", i), func(t *testing.T) {
// 			defer func() {
// 				r := recover()
// 				t.Logf("\n%v", r)
// 				if r == nil {
// 					t.Fatal("Haven`t panic for not valid data")
// 				}
// 			}()
// 			_ = golis.NewSparseMatrix(tc.r, tc.c)
// 		})
// 	}
//
// 	sp := golis.NewSparseMatrix(3, 2)
// 	for i, tc := range []struct{ r, c int }{
// 		{-1, 1},
// 		{1, -1},
// 		{-1, -1},
// 		{5, 1},
// 		{1, 5},
// 	} {
// 		t.Run(fmt.Sprintf("PanicAt%d", i), func(t *testing.T) {
// 			defer func() {
// 				r := recover()
// 				t.Logf("\n%v", r)
// 				if r == nil {
// 					t.Fatal("Haven`t panic for not valid data: ", tc)
// 				}
// 			}()
// 			_ = sp.At(tc.r, tc.c)
// 		})
// 	}
//
// 	for i, tc := range []struct{ r, c int }{
// 		{-1, 1},
// 		{1, -1},
// 		{-1, -1},
// 		{5, 1},
// 		{1, 5},
// 	} {
// 		t.Run(fmt.Sprintf("PanicSet%d", i), func(t *testing.T) {
// 			defer func() {
// 				r := recover()
// 				t.Logf("\n%v", r)
// 				if r == nil {
// 					t.Fatal("Haven`t panic for not valid data: ", tc)
// 				}
// 			}()
// 			sp.Set(tc.r, tc.c, 0)
// 		})
// 	}
//
// 	for i, tc := range []struct{ rc int }{
// 		{-1},
// 		{5},
// 	} {
// 		t.Run(fmt.Sprintf("PanicSetZeroRowAndColumn%d", i), func(t *testing.T) {
// 			defer func() {
// 				r := recover()
// 				t.Logf("\n%v", r)
// 				if r == nil {
// 					t.Fatal("Haven`t panic for not valid data: ", tc)
// 				}
// 			}()
// 			sp.SetZeroForRowColumn(tc.rc)
// 		})
// 	}
//
// 	t.Run("PanicSetInf-1", func(t *testing.T) {
// 		defer func() {
// 			r := recover()
// 			t.Logf("\n%v", r)
// 			if r == nil {
// 				t.Fatal("Haven`t panic for not valid data")
// 			}
// 		}()
// 		sp.Set(0, 0, math.Inf(-1))
// 	})
//
// 	t.Run("PanicSetInf+1", func(t *testing.T) {
// 		defer func() {
// 			r := recover()
// 			t.Logf("\n%v", r)
// 			if r == nil {
// 				t.Fatal("Haven`t panic for not valid data")
// 			}
// 		}()
// 		sp.Set(0, 0, math.Inf(1))
// 	})
//
// 	t.Run("PanicSetNan", func(t *testing.T) {
// 		defer func() {
// 			r := recover()
// 			t.Logf("\n%v", r)
// 			if r == nil {
// 				t.Fatal("Haven`t panic for not valid data")
// 			}
// 		}()
// 		sp.Set(0, 0, math.NaN())
// 	})
//
// 	t.Run("PanicAddInf-1", func(t *testing.T) {
// 		defer func() {
// 			r := recover()
// 			t.Logf("\n%v", r)
// 			if r == nil {
// 				t.Fatal("Haven`t panic for not valid data")
// 			}
// 		}()
// 		sp.Add(0, 0, math.Inf(-1))
// 	})
//
// 	t.Run("PanicAddInf+1", func(t *testing.T) {
// 		defer func() {
// 			r := recover()
// 			t.Logf("\n%v", r)
// 			if r == nil {
// 				t.Fatal("Haven`t panic for not valid data")
// 			}
// 		}()
// 		sp.Add(0, 0, math.Inf(1))
// 	})
//
// 	t.Run("PanicAddNan", func(t *testing.T) {
// 		defer func() {
// 			r := recover()
// 			t.Logf("\n%v", r)
// 			if r == nil {
// 				t.Fatal("Haven`t panic for not valid data")
// 			}
// 		}()
// 		sp.Add(0, 0, math.NaN())
// 	})
// }
