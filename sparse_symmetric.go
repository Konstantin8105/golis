package golis

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// guarantee SparseMatrix have interface of gonum.mat.Matrix
var _ mat.MutableSymmetric = (*SparseMatrixSymmetric)(nil)

// SparseMatrixSymmetric is struct of sparse matrix
type SparseMatrixSymmetric struct {
	s *SparseMatrix
}

// NewSparseMatrixSymmetric return new sparse square matrix
func NewSparseMatrixSymmetric(size int) *SparseMatrixSymmetric {
	m := new(SparseMatrixSymmetric)
	m.s = NewSparseMatrix(size, size)
	return m
}

// At returns the value of a matrix element at row i, column j.
// It will panic if i or j are out of bounds for the matrix.
func (m *SparseMatrixSymmetric) At(r, c int) float64 {
	if r > c {
		return m.s.At(c, r)
	}
	return m.s.At(r, c)
}

// Dims returns the dimensions of a Matrix.
// Where: r - amount of rows, c - amount of columns.
func (m *SparseMatrixSymmetric) Dims() (r, c int) {
	return m.s.Dims()
}

// T returns the transpose of the Matrix. Whether T returns a copy of the
// underlying data is implementation dependent.
// This method may be implemented using the Transpose type, which
// provides an implicit matrix transpose.
func (m *SparseMatrixSymmetric) T() mat.Matrix {
	c := new(SparseMatrixSymmetric)
	c.s = new(SparseMatrix)
	c.s.data.ts = make([]triple, len(m.s.data.ts))
	copy(c.s.data.ts, m.s.data.ts)
	return c
}

// Symmetric returns the number of rows/columns in the matrix.
func (m *SparseMatrixSymmetric) Symmetric() int {
	return m.s.r
}

// Set set value in sparse matrix by address [r,c].
// If r,c outside of matrix, then create a panic.
// If value is not valid, then create panic.
func (m *SparseMatrixSymmetric) SetSym(r, c int, value float64) {
	if r > c {
		panic(fmt.Errorf("SparseMatrixSymmetric have only upper value: %d <= %d", r, c))
	}
	m.s.Set(r, c, value)
}

// Add is alternative of pattern m.Set(r,c, someValue + m.At(r,c)).
// Addition value to matrix element
func (m *SparseMatrixSymmetric) Add(r, c int, value float64) {
	if r > c {
		panic(fmt.Errorf("SparseMatrixSymmetric have only upper value: %d <= %d", r, c))
	}
	m.s.Add(r, c, value)
}

// SetZeroForRowColumn set zero for all matrix element on
// row and column `rc`
func (m *SparseMatrixSymmetric) SetZeroForRowColumn(rc int) {
	m.s.SetZeroForRowColumn(rc)
}

func (m *SparseMatrixSymmetric) String() string {
	return m.s.String()
}
