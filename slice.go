package cnv

import (
	"math"

	"github.com/tys-muta/go-ers"
)

func Map[Src, Dst any](sources []Src, mapper func(Src) (Dst, error)) ([]Dst, error) {
	dst := []Dst{}
	for _, src := range sources {
		if v, err := mapper(src); err != nil {
			return nil, ers.W(err)
		} else {
			dst = append(dst, v)
		}
	}
	return dst, nil
}

func ChunkWithLength[T any](slice []T, length int) [][]T {
	total := len(slice)
	if length <= 0 {
		return [][]T{slice}
	}
	if total <= 0 {
		return [][]T{}
	}
	if total <= length {
		return [][]T{slice}
	}

	end := length
	outputs := make([][]T, 0)
	for start := 0; start < total && end <= total; end += int(math.Min(float64(length), float64(total-end))) {
		outputs = append(outputs, slice[start:end])
		start += length
	}

	return outputs
}

func ChunkWithSize[T any](slice []T, size int) [][]T {
	remaining := len(slice)
	offset := 0
	outputs := [][]T{}
	for current := 0; remaining != 0; current += offset {
		if remaining < size {
			offset = remaining
		} else {
			offset = size
		}
		outputs = append(outputs, slice[current:current+offset])
		remaining -= offset
	}

	return outputs
}
