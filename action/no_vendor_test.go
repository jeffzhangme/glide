package action

import (
	"testing"

	"github.com/jeffzhangme/glide/msg"
)

func TestNoVendor(t *testing.T) {
	msg.Default.PanicOnDie = true
	NoVendor("../testdata/nv", false, false)
}
