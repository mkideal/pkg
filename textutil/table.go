package textutil

import (
	"bytes"
	"fmt"
	"io"
)

const (
	borderCross = "+"
	borderRow   = "-"
	borderCol   = "|"
)

// Table represents a string-matrix
type Table interface {
	RowCount() int
	ColCount() int
	Get(i, j int) string
}

// WriteTable formats table to writer
func WriteTable(w io.Writer, table Table) {
	rowCount, colCount := table.RowCount(), table.ColCount()
	if rowCount <= 0 || colCount <= 0 {
		return
	}
	widthArray := make([]int, colCount)
	for j := 0; j < colCount; j++ {
		maxWidth := 0
		for i := 0; i < rowCount; i++ {
			width := len(table.Get(i, j))
			if i == 0 || width > maxWidth {
				maxWidth = width
			}
		}
		widthArray[j] = maxWidth
	}
	rowBorder := rowBorderLine(widthArray)
	fmt.Fprint(w, rowBorder)
	for i := 0; i < rowCount; i++ {
		fmt.Fprint(w, "\n")
		writeTableRow(w, table, i, widthArray)
		fmt.Fprint(w, "\n")
		fmt.Fprint(w, rowBorder)
	}
	fmt.Fprint(w, "\n")
}

func rowBorderLine(widthArray []int) string {
	buf := bytes.NewBufferString(borderCross)
	for _, width := range widthArray {
		repeatWriteString(buf, borderRow, width+2)
		buf.WriteString(borderCross)
	}
	return buf.String()
}

func writeTableRow(w io.Writer, table Table, rowIndex int, widthArray []int) {
	fmt.Fprint(w, borderCol)
	colCount := table.ColCount()
	for j := 0; j < colCount; j++ {
		fmt.Fprint(w, " ")
		format := fmt.Sprintf("%%-%ds", widthArray[j]+1)
		fmt.Fprintf(w, format, table.Get(rowIndex, j))
		fmt.Fprint(w, borderCol)
	}
}

func repeatWriteString(w io.Writer, s string, count int) {
	for i := 0; i < count; i++ {
		fmt.Fprint(w, s)
	}
}

// TableView represents a view of table, it implements Table interface, too
type TableView struct {
	table              Table
	rowIndex, colIndex int
	rowCount, colCount int
}

func (tv TableView) RowCount() int {
	return tv.rowCount
}

func (tv TableView) ColCount() int {
	return tv.colCount
}

func (tv TableView) Get(i, j int) string {
	return tv.table.Get(tv.rowCount+i, tv.colCount+j)
}

// ClipTable creates a view of table
func ClipTable(table Table, i, j, m, n int) Table {
	minR, minC := i, j
	maxR, maxC := i+m, j+n
	if minR < 0 || minC < 0 || minR > maxR || minC > maxC || maxR >= table.RowCount() || maxC >= table.ColCount() {
		panic("out of bound")
	}
	return &TableView{table, i, j, m, n}
}

// TableWithHeader add header for table
type TableWithHeader struct {
	table  Table
	header []string
}

func (twh TableWithHeader) RowCount() int { return twh.table.RowCount() + 1 }
func (twh TableWithHeader) ColCount() int { return twh.table.ColCount() }
func (twh TableWithHeader) Get(i, j int) string {
	if i == 0 {
		return twh.header[j]
	}
	return twh.table.Get(i-1, j)
}

func AddTableHeader(table Table, header []string) Table {
	return &TableWithHeader{table, header}
}

// 2-Array string
type StringMatrix [][]string

func (m StringMatrix) RowCount() int { return len(m) }
func (m StringMatrix) ColCount() int {
	if len(m) == 0 {
		return 0
	}
	return len(m[0])
}
func (m StringMatrix) Get(i, j int) string { return m[i][j] }
