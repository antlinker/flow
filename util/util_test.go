package util

import (
	"fmt"
	"testing"
)

func TestStringToInt(t *testing.T) {
	r, err := StringToInt("1.0")
	fmt.Println(r, err.Error())
}
