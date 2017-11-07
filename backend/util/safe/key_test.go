package safe

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestNewKey(t *testing.T) {
	_ = Suite(&KeySuite{})
	TestingT(t)
}

type KeySuite struct {
}

func (s *KeySuite) TestNewPassword(c *C) {
	k1 := NewKey("a", 0, []string{"c", "d"}, map[string]string{"e": "f"}, false)
	c.Assert(k1.String(), Equals, "[\"a\",0,[\"c\",\"d\"],{\"e\":\"f\"},false]")
}
