package golis

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Konstantin8105/errors"
	"gonum.org/v1/gonum/mat"
)

// LisPath is location of `lis` software.
// For example :
//	golis.LisPath = "/home/user/lis/bin/"
var LisPath string

// ErrorValue is error retirn value as result of `lis` software working
type ErrorValue int

func (g ErrorValue) Error() string {
	return fmt.Sprintf(
		"Error of `lis` software error: %s",
		errorStrings[int(g)])
}

// Constants of error values 'lis' software
const (
	IllOption ErrorValue = iota
	Breakdown
	OutOfMemory
	Maxiter
	NotImplemented
	ErrFileIO
)

var errorStrings = []string{
	"LIS_ILL_OPTION",
	"LIS_BREAKDOWN",
	"LIS_OUT_OF_MEMORY",
	"LIS_MAXITER",
	"LIS_NOT_IMPLEMENTED",
	"LIS_ERR_FILE_IO",
}

// Lsolve returns solution matrix of iterative solve for linear system.
//
//	A * x = b
//
// Where: A is matrix, b is right-hand vector.
//
// Description of rhsSetting and options, see in `lis` software documentation.
// Some examples:
//	options    = "-f quad"                    , Use quadriple precision
//	options    = "-i gmres -restart 20"       , Use solver GMRES with restart 20
//	options    = "-i bicgstab -maxiter 20000" , Use solver BiCGSTAB with max iteration 20000
//
func Lsolve(A, b mat.Matrix, options string) (
	solution mat.Matrix,
	rhistory []float64,
	output string,
	err error) {

	// check size of input Matrixs
	var et errors.Tree
	et.Name = "Check input matrix A and vector b"
	if r, c := A.Dims(); r != c {
		et.Add(fmt.Errorf("Matrix A is not square: [%d,%d]", r, c))
	}
	if r, c := b.Dims(); !(r > 0 && c == 1) {
		et.Add(fmt.Errorf("Vector b is not vertical vector: [%d,%d]", r, c))
	}
	{
		r, _ := A.Dims()
		if rb, _ := b.Dims(); r != rb {
			et.Add(fmt.Errorf("Amount of matrix and vector b is not same"))
		}
	}
	if et.IsError() {
		err = et
		return
	}

	// create a temp folder
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			err = fmt.Errorf("%v\nTemp folder: %v", err, tmpDir)
		}
	}()

	fn := func(name string) string {
		return filepath.Join(tmpDir, string(filepath.Separator), name)
	}

	// temp files
	var (
		inputFilename    = fn("input")
		solutionFilename = fn("solution.mtx")
		rhistoryFilename = fn("rhistory.txt")
	)

	inp := convertMatrixWithVector(A, b)
	err = ioutil.WriteFile(inputFilename, inp, 0644)
	if err != nil {
		return
	}

	// prepare arguments for `lis`
	args := []string{
		inputFilename,
		"0", // matrix b in inputFilename
		solutionFilename,
		rhistoryFilename,
	}
	args = append(args, strings.Split(options, " ")...)

	out, err := exec.Command(filepath.Join(LisPath, "lsolve"), args...).Output()
	if err != nil {
		err = fmt.Errorf("Error result of execute `lsolve`: %v", err)
		return
	}

	// Example of result parsing:
	// linear solver status  : normal end
	// linear solver status  : LIS_BREAKDOWN(code=2)
	for i, e := range errorStrings {
		if bytes.Contains(out, []byte(e)) {
			err = ErrorValue(i)
			return
		}
	}

	output = string(out)

	sol, err := ioutil.ReadFile(solutionFilename)
	if err != nil {
		return
	}
	solution, err = ParseSparseMatrix(sol)
	if err != nil {
		return
	}

	// parsing rhistory
	rhistory, err = parseRHistory(rhistoryFilename)
	if err != nil {
		return
	}

	return
}

// parseRHistory parsing rhs history
//
// Example:
//
//	1.000000e+00
//	0.000000e+00
func parseRHistory(rh string) (r []float64, err error) {
	b, err := ioutil.ReadFile(rh)
	if err != nil {
		return
	}

	lines := bytes.Split(b, []byte("\n"))
	for i := range lines {
		if len(lines[i]) == 0 {
			continue
		}
		s, err := strconv.ParseFloat(string(lines[i]), 64)
		if err != nil {
			err = fmt.Errorf("Cannot parse value `%v`: %v", string(lines[i]), err)
			return nil, err
		}
		r = append(r, s)
	}

	return
}
