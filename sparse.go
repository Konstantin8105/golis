package golis

import (
	"math"
	"sort"
)

type triple struct {
	position int64   // position matrix element (row + column * size)
	d        float64 // data
}

type SparseSquareMatrix struct {
	size int // size of matrix
	data struct {
		ts          []triple // non-zero value in matrix
		amountAdded int      // amount unsorted of triples
	}
}

// NewSparseSquareMatrix return new sparse square matrix
func NewSparseSquareMatrix(size int) *SparseSquareMatrix {
	m := new(SparseSquareMatrix)
	m.size = size
	m.data.ts = make([]triple, 0, size)
	return m
}

// At returns the value of a matrix element at row i, column j.
// It will panic if i or j are out of bounds for the matrix.
func (m *SparseSquareMatrix) At(r, c int) float64 {
	m.check(r, c)
	m.compress()

	// calculate position
	position := int64(r) + int64(c)*int64(m.size)

	// binary search of position
	index := sort.Search(len(m.data.ts), func(i int) bool {
		return m.data.ts[i].position >= position
	})

	if m.data.ts[index].position == position {
		return m.data.ts[index].d
	}

	return 0.0
}

func (m *SparseSquareMatrix) Set(r, c int, value float64) {
	m.check(r, c)
	checkValue(value)
	m.compress()

	// calculate position
	position := int64(r) + int64(c)*int64(m.size)

	// binary search of position
	index := sort.Search(len(m.data.ts), func(i int) bool {
		return m.data.ts[i].position >= position
	})

	if m.data.ts[index].position == position {
		m.data.ts[index].d = value
		return
	}

	m.data.ts = append(m.data.ts, triple{position: position, d: value})
	m.data.amountAdded++
}

// checkValue is panic if value is not correct: NaN or infinity.
func checkValue(v float64) {
	if math.IsNaN(v) {
		panic("value is NaN")
	}
	if math.IsInf(v, 0) {
		panic("value is infinity")
	}
}

func (m *SparseSquareMatrix) check(r, c int) {
	// TODO : add tree error panic
	if r < 0 || r >= m.size {
		panic("index out of range")
	}
	if c < 0 || c >= m.size {
		panic("index out of range")
	}
}

// compress triples data. Example of triples:
// [row column data]
// Before compression: [1 1 0.1] [1 2 0.5] [1 1 0.5]
// Intermediante     : [1 1 0.6] [1 2 0.5] [1 1 0.0]
// After  compression: [1 1 0.6] [1 2 0.5]
func (m *SparseSquareMatrix) compress() {
	if m.data.amountAdded == 0 {
		// compression is no need
		return
	}

	// sort by position
	sort.SliceStable(m.data.ts, func(i, j int) bool {
		return m.data.ts[i].position < m.data.ts[j].position
	})

	// summarize element with same indexes row, column and add 0.0 in old element
	for i := 1; i < len(m.data.ts); i++ {
		if m.data.ts[i-1].r != m.data.ts[i].r {
			continue
		}
		if m.data.ts[i-1].c != m.data.ts[i].c {
			continue
		}
		// triples element i-1 and i have same row and column
		m.data.ts[i-1].d += m.data.ts[i].d
		m.data.ts[i].d = 0.0
	}

	// moving data for avoid elements with 0.0 values
	for i := 0; i < len(m.data.ts); i++ {
	}

	// cut triple slice by nonzero elements
	// TODO
}

// Add is alternative of pattern m.Set(r,c, someValue + m.At(r,c)).
// Addition value to matrix element
func (m *SparseSquareMatrix) Add(r, c int, value float64) {
	m.check(r, c)
	checkValue(value)
	position := int64(r) + int64(c)*int64(m.size) // calculate position
	m.data.ts = append(m.data.ts, triple{position: position, d: value})
	m.data.amountAdded++
	if m.data.amountAdded > size {
		m.compress()
	}
}

// Dims returns the dimensions of a Matrix.
// Where: r - amount of rows, c - amount of columns.
func (m *SparseSquareMatrix) Dims() (r, c int) {
	return m.size, m.size
}
