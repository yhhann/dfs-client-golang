package client

import (
	"fmt"
	"path/filepath"
)

func getCacheFilePath(baseDir string, domain int64, fn string) string {
	result := filepath.Join(baseDir, fmt.Sprintf("g%d", domain%10000), fmt.Sprintf("%d", domain))

	if len(fn) > 0 {
		dummy := "00"
		dummy2 := "00"
		if len(fn) > 2 {
			dummy = fn[0:2]
		}
		if len(fn) > 4 {
			dummy2 = fn[2:4]
		}
		result = filepath.Join(result, dummy, dummy2, fn)
	}
	return result
}
