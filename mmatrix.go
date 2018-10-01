package golis

import (
	"bytes"
	"fmt"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

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
func convertMatrixWithVector(A, b mat.Matrix) []byte {
	// TODO : add to specific package mmatrix
	var buf bytes.Buffer

	buf.WriteString("%%MatrixMarket matrix coordinate real general\n")

	rA, cA := A.Dims()
	rb, cb := b.Dims()

	if cb != 1 {
		panic(fmt.Errorf("Input `b` is not vector: [%d,%d]", rb, cb))
	}

	// amount of non-zero values
	var nonZeros int
	// TODO add optimization for SparseMatrix
	switch v := A.(type) {
	case *SparseMatrix:
		nonZeros = len(v.data.ts)
	default:
		for i := 0; i < rA; i++ {
			for j := 0; j < cA; j++ {
				if A.At(i, j) != 0.0 {
					nonZeros++
				}
			}
		}
	}
	// write sizes
	// add string "1 0" for indicate that is matrix with vector
	buf.WriteString(fmt.Sprintf("%d %d %d 1 0\n", rA, cA, nonZeros))

	// write matrix A
	// TODO add optimization for SparseMatrix
	switch v := A.(type) {
	case *SparseMatrix:
		v.compress()
		for i := range v.data.ts {
			r := int(v.data.ts[i].position % int64(v.r))
			c := int(v.data.ts[i].position / int64(v.r))
			buf.WriteString(fmt.Sprintf("%d %d %20.16e\n", r+1, c+1, v.data.ts[i].d))
		}
	default:
		for i := 0; i < rA; i++ {
			for j := 0; j < cA; j++ {
				if A.At(i, j) != 0.0 {
					buf.WriteString(fmt.Sprintf("%d %d %20.16e\n", i+1, j+1, A.At(i, j)))
				}
			}
		}
	}
	// write vector b must be Dense
	for i := 0; i < rb; i++ {
		buf.WriteString(fmt.Sprintf("%d %20.16e\n", i+1, b.At(i, 0)))
	}

	return buf.Bytes()
}

// ParseSparseMatrix returns sparse matrix parsed from byte slice in
// MatrixMarket format and error, if exist
//
// Example:
//
//  %%MatrixMarket vector coordinate real general
//  3
//  1  -5.49999999999999822364e+00
//  2   2.49999999999999955591e+00
//  3   4.99999999999999911182e+00
//
func ParseSparseMatrix(b []byte) (v *SparseMatrix, err error) {

	// TODO add optimization for SparseMatrix
	v = new(SparseMatrix)

	lines := bytes.Split(b, []byte("\n"))

	// TODO: check vector
	// TODO: check real

	// convert size of vector
	s, err := strconv.ParseInt(string(lines[1]), 10, 64)
	if err != nil {
		err = fmt.Errorf("Cannot parse size `%v`: %v", string(lines[1]), err)
		return nil, err
	}
	v.r = int(s)
	v.c = 1

	v.data.ts = make([]triple, 0, v.r)

	// convert values
	for i := range lines {
		if i < 2 {
			continue
		}
		if len(bytes.TrimSpace(lines[i])) == 0 {
			continue
		}
		pars := bytes.Split(lines[i], []byte(" "))
		var t triple
		// parse index
		s, err := strconv.ParseInt(string(pars[0]), 10, 64)
		if err != nil {
			err = fmt.Errorf("Cannot parse index `%v`: %v", string(pars[0]), err)
			return nil, err
		}
		t.position = s - 1 // in MatrixMarket index from 1, but not zero

		// parse value
		for pos := 1; pos < len(pars); pos++ {
			if len(pars[pos]) == 0 {
				continue
			}
			s, err := strconv.ParseFloat(string(pars[pos]), 64)
			if err != nil {
				err = fmt.Errorf("Cannot parse value `%v`: %v", string(pars[pos]), err)
				return nil, err
			}
			t.d = s
		}
		v.data.ts = append(v.data.ts, t)
	}

	// compress
	v.data.amountAdded = -1
	v.compress()

	return v, nil
}
