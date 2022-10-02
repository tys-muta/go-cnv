package cnv

func Merge[T any](maps ...map[string]T) map[string]T {
	dst := map[string]T{}

	for _, m := range maps {
		for k, v := range m {
			dst[k] = v
		}
	}

	return dst
}
