package mody

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testInner struct {
	ValueStr  string
	ValueInt  int
	ValueBool bool
}

func TestUpdate_General(t *testing.T) {
	type test struct {
		FieldA testInner
	}

	v := &test{
		FieldA: testInner{
			ValueStr:  "abc",
			ValueInt:  1,
			ValueBool: false,
		},
	}

	assert.ErrorIs(t, Update(test{}, "", 0), ErrTypeMustBePointer)
	assert.ErrorIs(t, Update(v, "NotExist", "abc"), ErrFieldNotExistent)
	assert.ErrorIs(t, Update(v, "FieldA.NotExist", "abc"), ErrFieldNotExistent)

	assert.ErrorIs(t, Update(v, "FieldA.", "abc"), ErrFieldNotExistent)
	assert.ErrorIs(t, Update(v, ".ValueStr", "abc"), ErrFieldNotExistent)
	assert.ErrorIs(t, Update(v, "...", "abc"), ErrFieldNotExistent)
	assert.ErrorIs(t, Update(v, "", "abc"), ErrFieldNotExistent)
}

func TestUpdate_InnerValue(t *testing.T) {
	type test struct {
		FieldA testInner
	}

	v := &test{
		FieldA: testInner{
			ValueStr:  "abc",
			ValueInt:  1,
			ValueBool: false,
		},
	}

	assert.ErrorIs(t, Update(test{}, "", 0), ErrTypeMustBePointer)
	assert.ErrorIs(t, Update(v, "NotExist", "abc"), ErrFieldNotExistent)
	assert.ErrorIs(t, Update(v, "FieldA.NotExist", "abc"), ErrFieldNotExistent)

	assert.ErrorIs(t, Update(v, "FieldA.ValueStr", 1), ErrTypeMissmatch)
	assert.Nil(t, Update(v, "FieldA.ValueStr", "def"))
	assert.Equal(t, v.FieldA.ValueStr, "def")

	assert.ErrorIs(t, Update(v, "FieldA.ValueInt", "def"), ErrTypeMissmatch)
	assert.Nil(t, Update(v, "FieldA.ValueInt", 2))
	assert.Equal(t, v.FieldA.ValueInt, 2)

	assert.ErrorIs(t, Update(v, "FieldA.ValueBool", "def"), ErrTypeMissmatch)
	assert.Nil(t, Update(v, "FieldA.ValueBool", true))
	assert.Equal(t, v.FieldA.ValueBool, true)
}

func TestUpdate_InnerPtr(t *testing.T) {
	type test struct {
		FieldA *testInner
	}

	v := &test{
		FieldA: &testInner{
			ValueStr:  "abc",
			ValueInt:  1,
			ValueBool: false,
		},
	}

	assert.ErrorIs(t, Update(v, "FieldA.ValueStr", 1), ErrTypeMissmatch)
	assert.Nil(t, Update(v, "FieldA.ValueStr", "def"))
	assert.Equal(t, v.FieldA.ValueStr, "def")

	assert.ErrorIs(t, Update(v, "FieldA.ValueInt", "def"), ErrTypeMissmatch)
	assert.Nil(t, Update(v, "FieldA.ValueInt", 2))
	assert.Equal(t, v.FieldA.ValueInt, 2)

	assert.ErrorIs(t, Update(v, "FieldA.ValueBool", "def"), ErrTypeMissmatch)
	assert.Nil(t, Update(v, "FieldA.ValueBool", true))
	assert.Equal(t, v.FieldA.ValueBool, true)
}

func TestUpdateJson(t *testing.T) {
	type test struct {
		FieldA testInner
	}

	v := &test{
		FieldA: testInner{
			ValueStr:  "abc",
			ValueInt:  1,
			ValueBool: false,
		},
	}

	assert.ErrorIs(t, UpdateJson(v, "FieldA.ValueStr", "1"), ErrTypeMissmatch)
	assert.Nil(t, UpdateJson(v, "FieldA.ValueStr", `"def"`))
	assert.Equal(t, v.FieldA.ValueStr, "def")

	assert.ErrorIs(t, UpdateJson(v, "FieldA.ValueInt", `"def"`), ErrTypeMissmatch)
	assert.Nil(t, UpdateJson(v, "FieldA.ValueInt", "2"))
	assert.Equal(t, v.FieldA.ValueInt, 2)

	assert.ErrorIs(t, UpdateJson(v, "FieldA.ValueBool", `"def"`), ErrTypeMissmatch)
	assert.Nil(t, UpdateJson(v, "FieldA.ValueBool", "true"))
	assert.Equal(t, v.FieldA.ValueBool, true)
}

func TestCatch(t *testing.T) {
	assert.Nil(t, Catch(func() {
		fmt.Println("hey ho")
	}))
	assert.EqualError(t, Catch(func() {
		panic("test error")
	}), "test error")
}
