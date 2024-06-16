/*
 Copyright (c) 2024 fkmatsuda <fabio@fkmatsuda.dev>

 Permission is hereby granted, free of charge, to any person obtaining a copy of
 this software and associated documentation files (the "Software"), to deal in
 the Software without restriction, including without limitation the rights to
 use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 the Software, and to permit persons to whom the Software is furnished to do so,
 subject to the following conditions:

 The above copyright notice and this permission notice shall be included in all
 copies or substantial portions of the Software.

 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package errorex

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterErrorCode(t *testing.T) {
	t.Run("should register an error code", func(t *testing.T) {
		code := "test.code"
		description := "test description"
		detail := struct{ Message string }{}

		assert.NotPanics(t, func() {
			RegisterErrorCode(code, description, detail)
		})
	})

	t.Run("should panic if error code is already registered", func(t *testing.T) {
		code := ErrCodeAlreadyRegistered
		description := "test description"
		detail := struct{ Code string }{Code: code}

		expectedMessage := New(code, ErrorEXDetail{Code: code}).Error()
		assert.PanicsWithError(t, expectedMessage, func() {
			RegisterErrorCode(code, description, detail)
		})
	})
	t.Run("should create a new EX error with the given code and detail", func(t *testing.T) {
		code := "test.code"
		detail := struct{ Message string }{Message: "test detail"}

		ex := New(code, detail)

		assert.Equal(t, code, ex.Code())
		assert.Equal(t, detail, ex.Detail())
	})

}

func TestIs(t *testing.T) {

	RegisterErrorCode("test.is", "test description", struct{ Message string }{})

	t.Run("should return true if the error is of type EX and has the matching code", func(t *testing.T) {
		code := "test.is"
		ex := New(code, struct{ Message string }{Message: "test detail"})

		assert.True(t, Is(ex, code))
	})

	t.Run("should return false if the error is not of type EX", func(t *testing.T) {
		err := fmt.Errorf("some error")
		code := "test.is"

		assert.False(t, Is(err, code))
	})

	t.Run("should panic if the error code is not registered", func(t *testing.T) {
		ex := New("test.is", struct{ Message string }{})

		assert.Panics(t, func() {
			Is(ex, "unregistered.code")
		})
	})

	t.Run("should panic if the type of the provided detail is not the same type as specified in the record", func(t *testing.T) {
		panicDetail := ErrorEXDetailTypeMismatch{
			ExpectedType: "struct { Message string }",
			ActualType:   "struct { Label string }",
		}
		expectedMessage := New("errorex.003", panicDetail).Error()
		assert.PanicsWithError(t, expectedMessage, func() {
			New("test.is", struct{ Label string }{Label: "test detail"})
		})
	})

}

func TestEX_Error(t *testing.T) {
	t.Run("should return error message in JSON format", func(t *testing.T) {
		code := "test.format"

		RegisterErrorCode(code, "test description", struct{ Message string }{})

		detail := struct{ Message string }{Message: "test detail"}
		ex := New(code, detail)

		expected := `{"code": "test.format", "detail": {"Message":"test detail"}}`
		assert.JSONEq(t, expected, ex.Error())
	})
}

// Additional test cases should be added for other methods and error scenarios.
