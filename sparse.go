package golis

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"strconv"
)

type triple struct {
	position int64   // position matrix element (row + column * size)
	d        float64 // data
}

// TODO add research for finding limit size
// TODO create garantee for memory = amount of non-zero element + size
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
	m := new(SparseMatrix)
	m.r = r
	m.c = c
	m.data.ts = make([]triple, 0, r*c) // TODO: may be size must be more
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

// checkValue is panic if value is not correct: NaN or infinity.
func checkValue(v float64) {
	if math.IsNaN(v) {
		panic("value is NaN")
	}
	if math.IsInf(v, 0) {
		panic("value is infinity")
	}
}

func (m *SparseMatrix) check(r, c int) {
	// TODO : add tree error panic
	if r < 0 || r >= m.r {
		panic("index out of range")
	}
	if c < 0 || c >= m.c {
		panic("index out of range")
	}
}

// compress triples data. Example of triples:
// [row column data]
// Before compression: [1 1 0.1] [1 2 0.5] [1 1 0.5]
// Intermediante     : [1 1 0.6] [1 2 0.5] [1 1 0.0]
// After  compression: [1 1 0.6] [1 2 0.5]
func (m *SparseMatrix) compress() {
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

	// cut triple slice by nonzero elements
	for i := len(m.data.ts) - 1; i >= 0; i-- {
		if math.Abs(m.data.ts[i].d) != 0.0 || i == 0 {
			m.data.ts = m.data.ts[:i+1]
			break
		}
	}

	m.data.amountAdded = 0

	// check result of compression
	for i := 1; i < len(m.data.ts); i++ {
		if m.data.ts[i-1].position != m.data.ts[i].position {
			continue
		}
		// not correct compression
		panic(fmt.Errorf("Not correct compresstion: same position\n%v", m.stringByColumn()))
	}
}

// TODO : add standart output for sparse matrix
// TODO : add example
func (m *SparseMatrix) stringByColumn() string {
	s := "\n"
	for i := range m.data.ts {
		s += fmt.Sprintf("%5d) %5d %10.9e\n",
			i, m.data.ts[i].position, m.data.ts[i].d)
	}
	return s
}

// Add is alternative of pattern m.Set(r,c, someValue + m.At(r,c)).
// Addition value to matrix element
func (m *SparseMatrix) Add(r, c int, value float64) {
	m.check(r, c)
	checkValue(value)
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

// Dims returns the dimensions of a Matrix.
// Where: r - amount of rows, c - amount of columns.
func (m *SparseMatrix) Dims() (r, c int) {
	return m.r, m.c
}

// TODO: add function of matrix : get Min and Max absolute value for checking singular
// TODO: need research of memory for operation Add

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

	v = new(SparseMatrix)

	lines := bytes.Split(b, []byte("\n"))

	// TODO: check vector
	// TODO: check real

	// convert size of vector
	if s, err := strconv.ParseInt(string(lines[1]), 10, 64); err != nil {
		err = fmt.Errorf("Cannot parse size `%v`: %v", string(lines[1]), err)
		return nil, err
	} else {
		v.r = int(s)
		v.c = 0
	}

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
		if s, err := strconv.ParseInt(string(pars[0]), 10, 64); err != nil {
			err = fmt.Errorf("Cannot parse index `%v`: %v", string(pars[0]), err)
			return nil, err
		} else {
			t.position = s - 1 // in MatrixMarket index from 1, but not zero
		}
		// parse value
		for pos := 1; pos < len(pars); pos++ {
			if len(pars[pos]) == 0 {
				continue
			}
			if s, err := strconv.ParseFloat(string(pars[pos]), 64); err != nil {
				err = fmt.Errorf("Cannot parse value `%v`: %v", string(pars[pos]), err)
				return nil, err
			} else {
				t.d = s
			}
		}
		v.data.ts = append(v.data.ts, t)
	}

	return v, nil
}
