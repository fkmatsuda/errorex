/*
 *   Copyright (c) 2024 fkmatsuda <fabio@fkmatsuda.dev>
 *   All rights reserved.

 *   Permission is hereby granted, free of charge, to any person obtaining a copy
 *   of this software and associated documentation files (the "Software"), to deal
 *   in the Software without restriction, including without limitation the rights
 *   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *   copies of the Software, and to permit persons to whom the Software is
 *   furnished to do so, subject to the following conditions:

 *   The above copyright notice and this permission notice shall be included in all
 *   copies or substantial portions of the Software.

 *   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *   SOFTWARE.
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

	t.Run("should return false if the error is nil", func(t *testing.T) {
		assert.False(t, Is(nil, "test.is"))
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

func TestEX_New(t *testing.T) {
	t.Run("should panic if error code is not registered", func(t *testing.T) {
		assert.Panics(t, func() {
			New("unregistered.code", struct{ Message string }{})
		})
	})
}

func TestEX_Converter(t *testing.T) {
	t.Run("should return an unknown error converter", func(t *testing.T) {
		converter := NewUnknownErrorConverter()
		assert.NotNil(t, converter)
		assert.IsType(t, &unknownErrorConverter{}, converter)
	})

	t.Run("should convert an error into an unknown error", func(t *testing.T) {
		converter := NewUnknownErrorConverter()
		err := fmt.Errorf("test error")
		expectedMessage := New(ErrCodeUnknownError, UnknownErrorDetail{Detail: err.Error()}).Error()
		assert.Equal(t, expectedMessage, converter.ConvertError(err).Error())
	})

	t.Run("should test a second converter", func(t *testing.T) {
		RegisterErrorCode(ErrCodeMockError, "test description", MockErrorDetail{})

		converter := &mockErrorConverter{}
		converter.SetNext(NewUnknownErrorConverter())

		err := fmt.Errorf("test error")
		expectedMessage := converter.ConvertError(err)
		assert.True(t, Is(expectedMessage, ErrCodeMockError))

		err = fmt.Errorf("unknown error")
		unknownError := converter.ConvertError(err)
		assert.True(t, Is(unknownError, ErrCodeUnknownError))
	})

	t.Run("base converter should return nil if there is no next handler", func(t *testing.T) {
		converter := BaseErrorConverter{}

		err := fmt.Errorf("test error")
		assert.Nil(t, converter.ConvertError(err))
	})

	t.Run("base converter should return UnknownError with UnknowErrorConverter in the chain", func(t *testing.T) {
		converter := BaseErrorConverter{}
		converter.SetNext(NewUnknownErrorConverter())

		err := fmt.Errorf("test error")
		unknownError := converter.ConvertError(err)
		assert.True(t, Is(unknownError, ErrCodeUnknownError))
	})

}

// Mocks

const ErrCodeMockError = "ErrCodeMockError"

type MockErrorDetail struct {
	Detail string `json:"detail"`
}

type mockErrorConverter struct {
	BaseErrorConverter
}

func (m *mockErrorConverter) ConvertError(err error) EX {
	if err.Error() != "test error" {
		return m.next.ConvertError(err)
	}
	return New(ErrCodeMockError, MockErrorDetail{Detail: err.Error()})
}
