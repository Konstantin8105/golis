package golis

//
// // guarantee have interface of gonum.mat.Matrix
// var _ mat.Matrix = (*SkylineSymmetricMatrix)(nil)
//
// // SkylineSymmetricMatrix storage symmetric matrix in skyline format
// type SkylineSymmetricMatrix struct {
// 	size int // amount of element of square matrix
// 	sc   []skylineColumn
// }
//
// // skylineColumn is struct in skyline format of column.
// //
// // Example :
// // We have column with size from 0 to diagonal is 10.
// //
// // Element of column [ 0 0 0 0 0 1 2 3 4 5 ].
// // So in skylineColumn strorage {azv: 6, d:[1 2 3 4 5]}
// //
// // Element of column [ 1 2 3 4 5 6 7 8 9 10].
// // So in skylineColumn strorage {azv: 0, d:[1 2 3 4 5 6 7 8 9 10]}
// type skylineColumn struct {
// 	azv int       // amount zero values in column or position of first non-zero row in column
// 	d   []float64 // data from azv row to diagonal
// }
//
// func NewSkylineSymmetricMatrix(size int) *SkylineSymmetricMatrix {
// 	var et errors.Tree
// 	et.Name = "Check size of matrix"
// 	if size < 0 {
// 		et.Add(fmt.Errorf("Size of matrix cannot be less zero : %d", size))
// 	}
// 	if size == 0 {
// 		et.Add(fmt.Errorf("Size of matrix cannot be zero"))
// 	}
// 	if et.IsError() {
// 		panic(et)
// 	}
//
// 	m := new(SparseMatrix)
// 	m.size = size
// 	m.sc = make([]skylineColumn, size)
// 	// allocate memory
// 	minimalWidth := 6
// 	for i := range m.sc {
// 		m.sc[i].azv = i
// 		m.sc[i].d = make([]float64, 0, minimalWidth)
// 	}
// }
