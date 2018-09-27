package golis

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// LisPath is location of `lis` software.
// For example :
//	golis.LisPath = "/home/user/lis/bin/"
var LisPath string

// Matrix interface is must comparable with gonum.mat.Matrix interface
type Matrix interface {
	// Dims returns the dimensions of a Matrix.
	// Where: r - amount of rows, c - amount of columns.
	Dims() (r, c int)

	// At returns the value of a matrix element at row i, column j.
	// It will panic if i or j are out of bounds for the matrix.
	At(i, j int) float64
}

// convertMatrixWithVector - convert matrix and vector to byte slice in
// Matrix Market format
// See description:
// https://math.nist.gov/MatrixMarket/formats.html
//
// Coordinate Format for Sparse Matrices
// Format of MM        : coordinate
// Type of output data : matrix with vector
// Type of values      : real
// Type of matrix      : general
func convertMatrixWithVector(A, b Matrix) []byte {
	var buf bytes.Buffer

	buf.WriteString("%%MatrixMarket matrix coordinate real general\n")

	rA, cA := A.Dims()
	rb, cb := b.Dims()

	if cb != 1 {
		panic(fmt.Errorf("Input `b` is not vector: [%d,%d]", rb, cb))
	}

	// amount of non-zero values
	var nonZeros int
	for i := 0; i < rA; i++ {
		for j := 0; j < cA; j++ {
			if A.At(i, j) != 0.0 {
				nonZeros++
			}
		}
	}
	// write sizes
	// TODO: is need "1 0"
	buf.WriteString(fmt.Sprintf("%d %d %d 1 0\n", rA, cA, nonZeros))

	// write matrix A
	for i := 0; i < rA; i++ {
		for j := 0; j < cA; j++ {
			if A.At(i, j) != 0.0 {
				buf.WriteString(fmt.Sprintf("%d %d %20.16e\n", i+1, j+1, A.At(i, j)))
			}
		}
	}
	// write vector b
	for i := 0; i < rb; i++ {
		buf.WriteString(fmt.Sprintf("%d %20.16e\n", i+1, b.At(i, 0)))
	}

	return buf.Bytes()
}

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
// solve : Ax=b, where A is matrix, b is vector
// TODO: add "option" description
// TODO: add description
func Lsolve(A, b Matrix, rhsSetting, options string) (
	solution Matrix,
	rhistory []float64,
	output string,
	err error) {

	// check size of input Matrixs
	if r, c := A.Dims(); r != c && r > 0 {
		err = fmt.Errorf("Matrix A is not square: [%d,%d]", r, c)
		return
	}
	if r, c := b.Dims(); !(r > 0 && c == 1) {
		err = fmt.Errorf("Vector b is not vertical vector: [%d,%d]", r, c)
		return
	}
	{
		r, _ := A.Dims()
		if rb, _ := b.Dims(); r != rb {
			err = fmt.Errorf("Amount of matrix and vector b is not same")
			return
		}
	}

	if rhsSetting == "" {
		rhsSetting = "0"
	}

	// TODO: add error tree

	// create a temp folder
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return
	}

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
		rhsSetting,
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
