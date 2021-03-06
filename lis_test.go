package golis_test

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/Konstantin8105/golis"
	"gonum.org/v1/gonum/mat"
)

func init() {
	goPath := os.Getenv("GOPATH")
	lisPath := filepath.Join(goPath, "src/github.com/Konstantin8105/golis/bin/bin")

	golis.LisPath = lisPath
}

func TestLsolve(t *testing.T) {
	A := mat.NewDense(2, 2, []float64{
		1.0, 2.0,
		4.0, 1.0,
	})
	b := mat.NewDense(2, 1, []float64{
		4.0,
		9.0,
	})

	s, _, _, err := golis.Lsolve(A, b, "")
	if err != nil {
		t.Fatalf("Not correct result: %v", err)
	}

	if math.Abs(s.At(0, 0)-2) >= 1e-10 {
		t.Errorf("Element 0,0 is not correct : %v", s.At(0, 0))
	}
	if math.Abs(s.At(1, 0)-1) >= 1e-10 {
		t.Errorf("Element 1,0 is not correct : %v", s.At(1, 0))
	}
}

func TestLsolveQuad(t *testing.T) {
	A := mat.NewDense(2, 2, []float64{
		1.0, 2.0,
		1.0e-10, 1.0e-10,
	})
	b := mat.NewDense(2, 1, []float64{
		3.0,
		2.0e-10,
	})

	s, _, _, err := golis.Lsolve(A, b, "-f quad")
	if err != nil {
		t.Fatalf("Not correct result: %v", err)
	}

	if math.Abs(s.At(0, 0)-1) >= 1e-10 {
		t.Errorf("Element 0,0 is not correct : %v", s.At(0, 0))
	}
	if math.Abs(s.At(1, 0)-1) >= 1e-10 {
		t.Errorf("Element 1,0 is not correct : %v", s.At(1, 0))
	}
}

func TestLsolveOptions(t *testing.T) {
	A := mat.NewDense(2, 2, []float64{
		1.0, 2.0,
		4.0, 1.0,
	})
	b := mat.NewDense(2, 1, []float64{
		4.0,
		9.0,
	})

	options := []string{
		// Empty
		"",

		// Quadriple precision
		"-f quad",

		// Solvers
		"-i cg",                   // CG solver
		"-i cg -maxiter 20000",    // CG solver
		"-i bicg",                 // BiCG solver
		"-i cgs",                  // CGS
		"-i bicgstab",             // BiCGSTAB
		"-i bicgstabl",            // BiCGSTAB(l)
		"-i bicgstabl -ell 2",     // BiCGSTAB(l)
		"-i bicgstabl -ell 3",     // BiCGSTAB(l)
		"-i gpbicg",               // GPBiCG
		"-i tfqmr",                // TFQMR
		"-i orthomin",             // Orthomin(m)
		"-i orthomin -restart 20", // Orthomin(m)
		"-i gmres",                // GMRES
		"-i gmres -restart 20",    // GMRES
		"-i jacobi",               // Jacobi
		"-i gs",                   // Gauss-Seidel
		"-i sor",                  // SOR
		"-i sor -omega 1.5",       // SOR
		"-i bicgsafe",             // BiCGSafe
		"-i cr",                   // CR
		"-i bicr",                 // BiCR
		"-i crs",                  // CRS
		"-i bicrstab",             // BiCRSTAB
		"-i gpbicr",               // GPBiCR
		"-i bicrsafe",             // BiCRSafe
		"-i fgmres",               // FGMRES(m)
		"-i fgmres -restart 30",   // FGMRES(m)
		"-i idrs",                 // IDR(s)
		"-i idrs -irestart 2",     // IDR(s)
		"-i idr1",                 // IDR(1)
		"-i minres",               // MINRES
		"-i cocg",                 // COCG
		"-i cocr",                 // COCR

		// Solver with Preconditioners
		"-i bicrstab -p none",                       // none
		"-i bicrstab -p jacobi",                     // Jacobi
		"-i bicrstab -p ilu",                        // ILU(k)
		"-i bicrstab -p ilu -ilu_fill 0",            // ILU(k)
		"-i bicrstab -p ssor",                       // SSOR
		"-i bicrstab -p ssor -ssor_omega 1.5",       // SSOR
		"-i bicrstab -p hybrid",                     // Hybrid
		"-i bicrstab -p hybrid -hybrid_i bicg",      // Hybrid
		"-i bicrstab -p hybrid -hybrid_maxiter 20",  // Hybrid
		"-i bicrstab -p hybrid -hybrid_tol  1e-4",   // Hybrid
		"-i bicrstab -p hybrid -hybrid_omega  1.4",  // Hybrid
		"-i bicrstab -p hybrid -hybrid_ell 2 ",      // Hybrid
		"-i bicrstab -p hybrid -hybrid_restart 30 ", // Hybrid
		"-i bicrstab -p is",                         // I+S
		"-i bicrstab -p is -is_alpha 1.0",           // I+S
		"-i bicrstab -p is -is_m 3",                 // I+S
		"-i bicrstab -p sainv",                      // SAINV
		"-i bicrstab -p sainv -sainv_drop 0.05",     // SAINV
		"-i bicrstab -p saamg",                      // SA-AMG
		"-i bicrstab -p saamg -saamg_unsym false",   // SA-AMG
		"-i bicrstab -p saamg -saamg_theta 0.05",    // SA-AMG
		"-i bicrstab -p iluc",                       // Crout ILU
		"-i bicrstab -p iluc -iluc_drop 0.05",       // Crout ILU
		"-i bicrstab -p iluc -iluc_rate 5.0",        // Crout ILU
		"-i bicrstab -p ilut",                       // ILUT
		"-i bicrstab -adds true",                    // Additive Schwarz
		"-i bicrstab -adds true -adds_iter 1",       // Additive Schwarz
	}

	for _, opt := range options {
		t.Run(fmt.Sprintf("Option%s", opt), func(t *testing.T) {
			s, _, _, err := golis.Lsolve(A, b, opt)
			if err != nil {
				t.Log(err)
				return
			}

			if math.Abs(s.At(0, 0)-2) >= 1e-10 {
				t.Errorf("Element 0,0 is not correct : %v", s.At(0, 0))
			}
			if math.Abs(s.At(1, 0)-1) >= 1e-10 {
				t.Errorf("Element 1,0 is not correct : %v", s.At(1, 0))
			}
		})
	}
}

func TestLsolveFail(t *testing.T) {
	for i, tc := range []struct {
		a, b []float64
	}{
		// TODO: need internal checking
		// {
		// 	a: []float64{
		// 		1.0, 2.0,
		// 		0.0, 0.0,
		// 	},
		// 	b: []float64{
		// 		3.0,
		// 		0.0,
		// 	},
		// },
		// TODO: need internal checking
		// {
		// 	a: []float64{
		// 		1.0, 2.0,
		// 		1.0e-30, 1.0e-30,
		// 	},
		// 	b: []float64{
		// 		3.0,
		// 		2.0e-30,
		// 	},
		// },
		{
			a: []float64{
				1.0, 0.0,
				1.0, 0.0,
			},
			b: []float64{
				3.0,
				1.0,
			},
		},
	} {
		t.Run(fmt.Sprintf("Fail%d", i), func(t *testing.T) {
			A := mat.NewDense(2, 2, tc.a)
			B := mat.NewDense(2, 1, tc.b)

			_, _, _, err := golis.Lsolve(A, B, "")
			t.Logf("\n%v", err)
			if err == nil {
				t.Fatalf("Haven`t error : %v", err)
			}
		})
	}

	for i, tc := range []struct {
		rA, cA, rB, cB int
	}{
		{2, 2, 1, 1},
		{2, 2, 3, 1},
		{2, 3, 3, 1},
		{3, 3, 1, 3},
	} {
		t.Run(fmt.Sprintf("FailSize%d", i), func(t *testing.T) {
			A := golis.NewSparseMatrix(tc.rA, tc.cA)
			B := golis.NewSparseMatrix(tc.rB, tc.cB)

			_, _, _, err := golis.Lsolve(A, B, "")
			t.Logf("\n%v", err)
			if err == nil {
				t.Fatalf("Haven`t error : %v", err)
			}
		})
	}
}

func BenchmarkLsolve(b *testing.B) {
	size := 500

	A := mat.NewDense(size, size, nil)
	B := mat.NewDense(size, 1, nil)
	for i := 0; i < size; i++ {
		B.Set(i, 0, float64(i))
		for j := 0; j < size; j++ {
			if j > i-5 && j < i+5 {
				A.Set(i, j, float64(i+5*j))
			}
		}
	}

	var err error
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, err = golis.Lsolve(A, B, "-f quad -maxiter 40000")
		if err != nil {
			panic(err)
		}
	}
}
