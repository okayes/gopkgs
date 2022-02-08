package serverinfo

import (
    "testing"
)

func TestIP(t *testing.T) {
    ip, err := IP()
    t.Log(ip, err)
}
