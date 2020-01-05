package sfc

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/zc310/utils"
)

func TestSfc(t *testing.T) {
	q, err := Sfc("2019112", false)
	assert.Equal(t, err, nil)
	utils.PrintJSON(q)
}
func TestSfcHisIssue(t *testing.T) {
	q, err := SfcHisIssue()
	assert.Equal(t, err, nil)
	utils.PrintJSON(q)
}
