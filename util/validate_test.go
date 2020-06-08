package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueIsZeroValue(t *testing.T) {
	v := ValueValidationImpl{}
	v.Load(1, "")

	assert.Equal(t, true, v.IfHasNilValue(), "wrong ")

	v.Load("")
	assert.Equal(t, true, v.IfHasNilValue(), "wrong ")

	var int64test int64 = 124555
	v.Load(int64test)
	assert.Equal(t, false, v.IfHasNilValue(), "wrong ")
	v.Load(&int64test)
	assert.Equal(t, false, v.IfHasNilValue(), "wrong ")

	var int64nil int64
	v.Load(int64nil)
	assert.Equal(t, true, v.IfHasNilValue(), "wrong ")
	v.Load(&int64nil)
	assert.Equal(t, false, v.IfHasNilValue(), "wrong ")

	v.Load(&int64nil, 2)
	assert.Equal(t, false, v.IfHasNilValue(), "wrong ")

}
