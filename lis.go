package golis

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

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

// solve : Ax=b, where A is matrix, b is vector
// TODO: add "option" description
// TODO: add description
func Lsolve(A, b Matrix, option int) (
	solution Matrix,
	rhistory Matrix,
	output string,
	err error) {

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

	// TODO : specific folder for lis application
	out, err := exec.Command("/home/lepricon/lis/bin/lsolve",
		inputFilename,
		fmt.Sprintf("%d", option),
		solutionFilename,
		rhistoryFilename,
		"-f", "quad", // double-double (quadruple) precision
	).Output()
	if err != nil {
		return
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

	// TODO: rhistory

	// TODO: Read line "linear solver status  : normal end"

	return
}
