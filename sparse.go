package golis

import (
	"fmt"
	"math"
	"sort"

	"github.com/Konstantin8105/errors"
	"gonum.org/v1/gonum/mat"
)

// guarantee SparseMatrix have interface of gonum.mat.Matrix
var _ mat.Matrix = (*SparseMatrix)(nil)

type triple struct {
	position int64   // position matrix element (row + column * size)
	d        float64 // data
}

// byTriple implements sort.Interface based on the position field.
type byTriple []triple

func (a byTriple) Len() int           { return len(a) }
func (a byTriple) Less(i, j int) bool { return a[i].position < a[j].position }
func (a byTriple) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// TODO add research for finding limit size
// TODO create guarantee for memory = amount of non-zero element + size
// TODO use memory blocks for triples separate by size L2 cache

// SparseMatrix is struct of sparse matrix
type SparseMatrix struct {
	r    int // amount of matrix rows
	c    int // amount of matrix columns
	data struct {
		ts          []triple // non-zero value in matrix
		amountAdded int      // amount unsorted of triples
	}
}

// NewSparseMatrix return new sparse square matrix
func NewSparseMatrix(r, c int) *SparseMatrix {
	var et errors.Tree
	et.Name = "Check size of matrix"
	if r < 0 {
		et.Add(fmt.Errorf("Size of rows cannot be less zero : %d", r))
	}
	if r == 0 {
		et.Add(fmt.Errorf("Size of rows cannot be zero"))
	}
	if c < 0 {
		et.Add(fmt.Errorf("Size of columns cannot be less zero : %d", c))
	}
	if c == 0 {
		et.Add(fmt.Errorf("Size of columns cannot be zero"))
	}
	if et.IsError() {
		panic(et)
	}

	m := new(SparseMatrix)
	m.r = r
	m.c = c
	// allocate memory for triplets
	switch {
	case r == 1: // vector
		m.data.ts = make([]triple, 0, c/2)

	case c == 1: // vector
		m.data.ts = make([]triple, 0, r/2)

	case r == c: // square matrix
		m.data.ts = make([]triple, 0, r)

	default:
		m.data.ts = make([]triple, 0, c)
	}
	return m
}

// At returns the value of a matrix element at row i, column j.
// It will panic if i or j are out of bounds for the matrix.
func (m *SparseMatrix) At(r, c int) float64 {
	m.check(r, c)
	m.compress()

	// calculate position
	position := int64(r) + int64(c)*int64(m.r)

	// binary search of position
	index := sort.Search(len(m.data.ts), func(i int) bool {
		return m.data.ts[i].position >= position
	})

	if index < len(m.data.ts) && m.data.ts[index].position == position {
		return m.data.ts[index].d
	}

	return 0.0
}

// Set set value in sparse matrix by address [r,c].
// If r,c outside of matrix, then create a panic.
// If value is not valid, then create panic.
func (m *SparseMatrix) Set(r, c int, value float64) {
	m.check(r, c)
	checkValue(value)
	m.compress()

	// calculate position
	position := int64(r) + int64(c)*int64(m.r)

	// binary search of position
	index := sort.Search(len(m.data.ts), func(i int) bool {
		return m.data.ts[i].position >= position
	})

	if index < len(m.data.ts) && m.data.ts[index].position == position {
		m.data.ts[index].d = value
		return
	}

	// TODO: append can multiply memory by 2 - it is not effective
	m.data.ts = append(m.data.ts, triple{position: position, d: value})
	m.data.amountAdded++
}

// SetZeroForRowColumn set zero for all matrix element on
// row and column `rc`
func (m *SparseMatrix) SetZeroForRowColumn(rc int) {
	m.check(rc, rc)
	for i := range m.data.ts {
		if int(m.data.ts[i].position%int64(m.r)) == rc {
			// zero on rows
			m.data.ts[i].d = 0.0
			continue
		}
		if int(m.data.ts[i].position/int64(m.r)) == rc {
			// zero on columns
			m.data.ts[i].d = 0.0
			continue
		}
	}
	m.data.amountAdded = -1
}

// checkValue is panic if value is not correct: NaN or infinity.
func checkValue(v float64) {
	if math.IsNaN(v) {
		panic("Value is not valid : NaN")
	}
	if math.IsInf(v, 0) {
		panic("Value is not valid : infinity")
	}
}

func (m *SparseMatrix) check(r, c int) {
	var et errors.Tree
	et.Name = "Check input indexes of element"

	if r < 0 {
		et.Add(fmt.Errorf("Index of rows cannot be less zero : %d", r))
	}
	if r >= m.r {
		et.Add(fmt.Errorf("Index of rows is outside of matrix: %d of %d", r, m.r))
	}
	if c < 0 {
		et.Add(fmt.Errorf("Index of columns cannot be less zero : %d", c))
	}
	if c >= m.c {
		et.Add(fmt.Errorf("Index of columns is outside of matrix: %d of %d", c, m.c))
	}
	if et.IsError() {
		panic(et)
	}
}

// compress triples data. Example of triples:
// [row column data]
// Before compression: [1 1 0.1] [1 2 0.5] [1 1 0.5]
// Intermediante     : [1 1 0.6] [1 2 0.5] [1 1 0.0]
// After  compression: [1 1 0.6] [1 2 0.5]
func (m *SparseMatrix) compress() {
	// check only with zero for force compression in
	// parsing case
	if m.data.amountAdded == 0 {
		// compression is no need
		return
	}

	// sort by position
	sort.Sort(byTriple(m.data.ts))

	// summarize element with same indexes row, column and add 0.0 in old element
	for i := 1; i < len(m.data.ts); i++ {
		if m.data.ts[i-1].position != m.data.ts[i].position {
			continue
		}
		nonZero := i - 1
		for ; i < len(m.data.ts); i++ {
			if m.data.ts[nonZero].position != m.data.ts[i].position {
				break
			}
			// triples element i-1 and i have same row and column
			m.data.ts[nonZero].d += m.data.ts[i].d // TODO: add float64 limit checking
			m.data.ts[i].d = 0.0
		}
	}

	// moving data for avoid elements with 0.0 values
	var nonZeroPos int
	for zeroPos := 0; zeroPos < len(m.data.ts); zeroPos++ {
		// find position of zero value triple
		if math.Abs(m.data.ts[zeroPos].d) != 0.0 {
			continue
		}

		// find next non-zero value triple
		if nonZeroPos < zeroPos {
			nonZeroPos = zeroPos
		}
		for ; nonZeroPos < len(m.data.ts); nonZeroPos++ {
			if math.Abs(m.data.ts[nonZeroPos].d) != 0.0 {
				break
			}
		}
		if nonZeroPos >= len(m.data.ts) {
			break
		}

		// move value
		m.data.ts[zeroPos] = m.data.ts[nonZeroPos]
		m.data.ts[nonZeroPos].d = 0.0
	}

	{
		// cut triple slice by nonzero elements
		var cut int
		for cut = len(m.data.ts) - 1; cut >= 0; cut-- {
			if math.Abs(m.data.ts[cut].d) != 0.0 {
				break
			}
		}
		m.data.ts = m.data.ts[:cut+1]
	}

	m.data.amountAdded = 0

	// Only for debuging:
	// // check result of compression
	// for i := 1; i < len(m.data.ts); i++ {
	// 	if m.data.ts[i-1].position != m.data.ts[i].position {
	// 		continue
	// 	}
	// 	// not correct compression
	// 	panic(fmt.Errorf("Not correct compresstion: same position\n%s",
	// 		m.String()))
	// }
}

// TODO : fmt.Formatted

// String return standard golis string of sparse matrix
func (m *SparseMatrix) String() string {
	m.compress()
	s := "\n"
	s += fmt.Sprintf("Amount of rows    : %5d\n", m.r)
	s += fmt.Sprintf("Amount of columns : %5d\n", m.c)
	s += fmt.Sprintf("%-6s %-6s %20s\n", "row", "column", "value")
	if len(m.data.ts) == 0 {
		return s
	}
	pos := 0
	for c := 0; c < m.c; c++ {
		for r := 0; r < m.r; r++ {
			position := int64(r) + int64(c)*int64(m.r) // calculate position
			if m.data.ts[pos].position == position {
				s += fmt.Sprintf("%-6d %-6d %-20.15e\n",
					r, c, m.data.ts[pos].d)
				pos++
			}
			if pos >= len(m.data.ts) {
				goto end
			}
		}
	}
end:
	return s
}

// Add is alternative of pattern m.Set(r,c, someValue + m.At(r,c)).
// Addition value to matrix element
func (m *SparseMatrix) Add(r, c int, value float64) {
	m.check(r, c)
	checkValue(value)
	if math.Abs(value) == 0.0 { // no need addition zero value
		return
	}
	position := int64(r) + int64(c)*int64(m.r) // calculate position
	// TODO: append can multiply memory by 2 - it is not effective
	m.data.ts = append(m.data.ts, triple{position: position, d: value})
	m.data.amountAdded++
	max := m.c
	if m.r > m.c {
		max = m.r
	}
	if m.data.amountAdded > max {
		m.compress()
	}
}

// T returns the transpose of the Matrix. Whether T returns a copy of the
// underlying data is implementation dependent.
// This method may be implemented using the Transpose type, which
// provides an implicit matrix transpose.
func (m *SparseMatrix) T() mat.Matrix {
	m.compress()
	out := new(SparseMatrix)
	out.r = m.c
	out.c = m.r
	out.data.ts = make([]triple, 0, len(m.data.ts))
	pos := 0
	for c := 0; c < m.c; c++ {
		for r := 0; r < m.r; r++ {
			position := int64(r) + int64(c)*int64(m.r) // calculate position
			if m.data.ts[pos].position == position {
				out.Add(c, r, m.data.ts[pos].d)
				pos++
			}
			if pos >= len(m.data.ts) {
				goto end
			}
		}
	}
end:
	out.data.amountAdded = -1
	out.compress()
	return out
}

// Dims returns the dimensions of a Matrix.
// Where: r - amount of rows, c - amount of columns.
func (m *SparseMatrix) Dims() (r, c int) {
	return m.r, m.c
}

// TODO: add function of matrix : get Min and Max absolute value for checking singular
// TODO: need research of memory for operation Add
