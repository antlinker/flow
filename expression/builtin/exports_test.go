package builtin

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Reg()
	os.Exit(m.Run())
}
