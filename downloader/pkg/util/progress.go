package util

import (
	"fmt"
)

func PrintProgress(current, total int) {
	percent := (float64(current) / float64(total)) * 100
	fmt.Printf("\rProgress: %.2f%%", percent)
}
