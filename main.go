package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Matrix struct {
	rows int
	cols int
	data []float64
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	fmt.Printf("%v: %v\n", msg, time.Since(start))
}
func print_matrix(mat Matrix) {
	for i := 0; i < mat.rows; i++ {
		for j := 0; j < mat.cols; j++ {
			fmt.Printf("%g ", mat.data[i*mat.rows+j])
		}
		fmt.Printf("\n")
	}

}
func fill_index(mat Matrix) {
	for i := 0; i < mat.rows; i++ {
		for j := 0; j < mat.cols; j++ {
			mat.data[i*mat.rows+j] = float64(i*mat.rows + j)
		}
	}

}
func fill_value(mat Matrix, val float64) {
	for i := 0; i < mat.rows; i++ {
		for j := 0; j < mat.cols; j++ {
			mat.data[i*mat.rows+j] = val
		}
	}

}

func matmul_basic(m1, m2 Matrix) *Matrix {
	result := &Matrix{rows: m1.rows, cols: m2.cols, data: make([]float64, m1.rows*m2.cols)}
	return matmul_chunk(&m1, &m2, result, 0, m1.rows, nil)
}
func matmul_chunk(m1, m2, result *Matrix, i_min, i_max int, wg *sync.WaitGroup) *Matrix {
	if wg == nil {
		defer wg.Done()
	}
	for i := i_min; i < i_max; i++ {
		for j := 0; j < m2.cols; j++ {
			for k := 0; k < m2.rows; k++ {
				result.data[i*m2.rows+j] += m1.data[i*m1.rows+k] * m2.data[k*m2.rows+j]
			}
		}
	}
	return result
}
func matmul_parallel(m1 Matrix, m2 Matrix) *Matrix {
	defer duration(track("parallel"))
	result := &Matrix{rows: m1.rows, cols: m2.cols, data: make([]float64, m1.rows*m2.cols)}
	wg := &sync.WaitGroup{}
	wg.Add(4)
	rows := m1.rows
	go matmul_chunk(&m1, &m2, result, 0, rows/4, wg)
	go matmul_chunk(&m1, &m2, result, rows/4, rows/2, wg)
	go matmul_chunk(&m1, &m2, result, rows/2, rows*3/4, wg)
	go matmul_chunk(&m1, &m2, result, rows*3/4, rows, wg)
	wg.Wait()
	return result
}

func init_mat(rows int, cols int) *Matrix {
	var ret *Matrix = new(Matrix)
	ret.rows = rows
	ret.cols = cols
	ret.data = make([]float64, rows*cols)

	return ret
}

func main() {
	fmt.Println("NumCPU", runtime.NumCPU())
	fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))
	var m1 Matrix = *init_mat(90, 90)
	var m2 Matrix = *init_mat(90, 90)
	fill_value(m1, 2)
	fill_index(m2)
	matmul_parallel(m1, m2)
	matmul_basic(m1, m2)
}
