package sfc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetZcMatchText(t *testing.T) {
	s, err := GetZcMatchObject("20002")
	assert.Equal(t, err, nil)
	t.Log("\n")
	t.Log(s)
	t.Log(s.Odds())
}
