package cnv

import (
	"strconv"
)

func Int64(v any) int64 {
	switch v := v.(type) {
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return int64(v)
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	default:
		return 0
	}
}

func Int64Slice[T any](v []T) []int64 {
	dst := []int64{}
	for _, v := range v {
		dst = append(dst, Int64(v))
	}
	return dst
}
