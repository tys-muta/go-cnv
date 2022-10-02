package cnv

import (
	"fmt"
)

func String(src any) string {
	switch v := src.(type) {
	case int, int32, int64, uint, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%s", v)
	}
}
