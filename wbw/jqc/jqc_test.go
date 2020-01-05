package jqc

import (
	"github.com/magiconair/properties/assert"
	"github.com/zc310/utils"
	"testing"
)

func TestSfc(t *testing.T) {
	q, err := Match("2019068")
	assert.Equal(t, err, nil)
	utils.PrintJSON(q)
}

func TestSfcHisIssue(t *testing.T) {
	q, err := HisIssue()
	assert.Equal(t, err, nil)
	utils.PrintJSON(q)
}
