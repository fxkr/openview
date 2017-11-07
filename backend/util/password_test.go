package util

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestPassword(t *testing.T) {
	_ = Suite(&PasswordSuite{})
	TestingT(t)
}

type PasswordSuite struct {
}

func (s *PasswordSuite) TestNewPassword(c *C) {
	p1, err := NewPassword()
	c.Assert(err, IsNil)
	p2, err := NewPassword()
	c.Assert(err, IsNil)
	c.Assert(len(p1), Equals, 32)
	c.Assert(len(p2), Equals, 32)
	c.Assert(p1, Not(Equals), p2)
}
