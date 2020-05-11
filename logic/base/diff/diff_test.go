package diff

import (
	"fmt"
	"testing"
)

func TestDiff(t *testing.T) {
	str := DoTextDiff("abced","adbck")
	fmt.Println(str)
}
