package common

import (
	"testing"
)

func TestParallels8Driver_impl(t *testing.T) {
	var _ Driver = new(Parallels8Driver)
}
