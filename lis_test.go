package golis_test

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Konstantin8105/golis"
	"gonum.org/v1/gonum/mat"
)

var lisPath string = "/home/konstantin/lis/bin/"

func TestLsolve(t *testing.T) {
	// change location of lis software
	golis.LisPath = lisPath

	A := mat.NewDense(2, 2, []float64{
		1.0, 2.0,
		4.0, 1.0,
	})
	b := mat.NewDense(2, 1, []float64{
		4.0,
		9.0,
	})

	s, _, _, err := golis.Lsolve(A, b, "", "")
	if err != nil {
		t.Errorf("Not correct result: %v", err)
	}

	if math.Abs(s.At(0, 0)-2) >= 1e-10 {
		t.Errorf("Element 0,0 is not correct : %v", s.At(0, 0))
	}
	if math.Abs(s.At(1, 0)-1) >= 1e-10 {
		t.Errorf("Element 1,0 is not correct : %v", s.At(1, 0))
	}
}

func TestLsolveQuad(t *testing.T) {
	// change location of lis software
	golis.LisPath = lisPath

	A := mat.NewDense(2, 2, []float64{
		1.0, 2.0,
		1.0e-10, 1.0e-10,
	})
	b := mat.NewDense(2, 1, []float64{
		3.0,
		2.0e-10,
	})

	s, _, _, err := golis.Lsolve(A, b, "", "-f quad")
	if err != nil {
		t.Errorf("Not correct result: %v", err)
	}

	if math.Abs(s.At(0, 0)-1) >= 1e-10 {
		t.Errorf("Element 0,0 is not correct : %v", s.At(0, 0))
	}
	if math.Abs(s.At(1, 0)-1) >= 1e-10 {
		t.Errorf("Element 1,0 is not correct : %v", s.At(1, 0))
	}
}

func TestLsolveOptions(t *testing.T) {
	// change location of lis software
	golis.LisPath = lisPath

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
	}

	for _, opt := range options {
		t.Run(fmt.Sprintf("Option%s", opt), func(t *testing.T) {
			s, _, _, err := golis.Lsolve(A, b, "", opt)
			if err != nil {
				t.Logf("Not correct result: %v", err)
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
	// change location of lis software
	golis.LisPath = lisPath

	tcs := []struct {
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
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Fail%d", i), func(t *testing.T) {
			A := mat.NewDense(2, 2, tc.a)
			B := mat.NewDense(2, 1, tc.b)

			_, _, _, err := golis.Lsolve(A, B, "", "")
			if err == nil {
				t.Fatalf("Haven`t error : %v", err)
			}
		})
	}
}

func TestTodo(t *testing.T) {
	// Show all to do`s in comment code
	source, err := filepath.Glob(fmt.Sprintf("./%s", "*.go"))
	if err != nil {
		t.Fatal(err)
	}

	var amount int

	for i := range source {
		t.Run(source[i], func(t *testing.T) {
			file, err := os.Open(source[i])
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			pos := 0
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				pos++
				index := strings.Index(line, "//")
				if index < 0 {
					continue
				}
				if !strings.Contains(strings.ToUpper(line), "TODO") {
					continue
				}
				t.Logf("%13s:%-4d %s", source[i], pos, line[index:])
				amount++
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		})
	}
	if amount > 0 {
		t.Logf("Amount TODO: %d", amount)
	}
}

func TestFmt(t *testing.T) {
	// Show all fmt`s in comments code
	source, err := filepath.Glob(fmt.Sprintf("./%s", "*.go"))
	if err != nil {
		t.Fatal(err)
	}

	var amount int

	for i := range source {
		t.Run(source[i], func(t *testing.T) {
			file, err := os.Open(source[i])
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			pos := 1
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				pos++
				index := strings.Index(line, "//")
				if index < 0 {
					continue
				}
				if !strings.Contains(line, "fmt.") {
					continue
				}
				t.Logf("%d %s", pos, line[index:])
				amount++
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		})
	}
	if amount > 0 {
		t.Logf("Amount commented fmt: %d", amount)
	}
}

func TestDebug(t *testing.T) {
	// Show all debug information in code
	source, err := filepath.Glob(fmt.Sprintf("./%s", "*.go"))
	if err != nil {
		t.Fatal(err)
	}

	for i := range source {
		t.Run(source[i], func(t *testing.T) {
			file, err := os.Open(source[i])
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			pos := 1
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				pos++
				if !strings.Contains(line, "fmt"+"."+"Print") {
					continue
				}
				t.Errorf("Debug line: %d in file %s", pos, source[i])
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		})
	}
}
