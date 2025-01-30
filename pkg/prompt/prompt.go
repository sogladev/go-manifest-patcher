package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var ErrUserCancelled = errors.New("operation cancelled by user")

func PromptyN(message string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", message)

	input, err := reader.ReadString('\n')
	if err != nil {
		return ErrUserCancelled
	}

	input = strings.ToLower(strings.TrimSpace(input))
	if input != "y" {
		return ErrUserCancelled
	}

	return nil
}
