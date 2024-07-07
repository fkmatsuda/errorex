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
	"encoding/json"
	"fmt"
	"reflect"
)

var (
	errorCodes = make(map[string]errorCodeRegistry)
)

const (
	// ErrCodeUnknownError is the errorex code for when the errorex code is unknown
	ErrCodeUnknownError = "errorex.000"
	// ErrCodeNotRegistered is the errorex code for when the errorex code is not registered
	ErrCodeNotRegistered = "errorex.001"
	// ErrCodeAlreadyRegistered is the errorex code for when the errorex code is already registered
	ErrCodeAlreadyRegistered = "errorex.002"
	// ErrDetailTypeMismatch is the errorex code for when the errorex detail type does not match the registered type
	ErrDetailTypeMismatch = "errorex.003"
)

// UnknownErrorDetail is the type of the detail of an unknown errorex
type UnknownErrorDetail struct {
	Detail string `json:"detail"`
}

// ErrorEXDetail is the type of the detail of an errorex
type ErrorEXDetail struct {
	Code string `json:"code"`
}

// ErrorEXDetailTypeMismatch is the type of the detail of an errorex
type ErrorEXDetailTypeMismatch struct {
	ExpectedType string `json:"expectedType"`
	ActualType   string `json:"actualType"`
}

func init() {
	// Register the errorex codes
	RegisterErrorCode(ErrCodeUnknownError, "Unknown errorex", UnknownErrorDetail{})
	RegisterErrorCode(ErrCodeNotRegistered, "Errorex code not registered", ErrorEXDetail{})
	RegisterErrorCode(ErrCodeAlreadyRegistered, "Errorex code already registered", ErrorEXDetail{})
	RegisterErrorCode(ErrDetailTypeMismatch, "Errorex detail type mismatch", ErrorEXDetailTypeMismatch{})
}

// ErrorConstructor is a function that creates an errorEX
type ErrorConstructor[T any] func(code string, detail T) EX

// EX is a custom errorex type with additional information
type EX interface {
	error
	// Code is the errorex code
	Code() string
	// Detail returns the detail of the error
	Detail() any
}

type ex struct {
	code   string
	detail any
}

type errorCodeRegistry struct {
	code        string
	description string
	detailType  reflect.Type
}

// RegisterErrorCode registers errorex codes to prevent repeats
func RegisterErrorCode[T any](code string, description string, detail T) {
	// Prevent repeats
	if _, ok := errorCodes[code]; ok {
		// Fatal errorex
		panic(New(ErrCodeAlreadyRegistered, ErrorEXDetail{Code: code}))
	}
	// Register the errorex code
	registry := errorCodeRegistry{
		code:        code,
		description: description,
		detailType:  reflect.TypeOf(detail),
	}
	errorCodes[code] = registry
}

// Code returns the errorex code
func (e *ex) Code() string {
	return e.code
}

// Detail returns the errorex detail
func (e *ex) Detail() any {
	return e.detail
}

// Error returns the errorex message
func (e *ex) Error() string {
	detailJSON, err := json.Marshal(e.detail)
	if err != nil {
		return fmt.Sprintf(`{"code": "%s", "detail": "failed to marshal detail: %v"}`, e.code, err)
	}
	return fmt.Sprintf(`{"code": "%s", "detail": %s}`, e.code, string(detailJSON))
}

// New returns a new errorex.EX
// Code is the errorex code.
// Detail is the errorex detail.
func New[T any](code string, detail T) EX {
	// Check if the code exists
	var (
		errorRegistry errorCodeRegistry
		ok            bool
	)
	if errorRegistry, ok = errorCodes[code]; !ok {
		// Fatal errorex
		panic(New(ErrCodeNotRegistered, ErrorEXDetail{Code: code}))
	}
	// Check if the detail type matches the registered type
	if reflect.TypeOf(detail) != errorRegistry.detailType {
		// Fatal errorex
		panic(New(ErrDetailTypeMismatch, ErrorEXDetailTypeMismatch{
			ExpectedType: errorRegistry.detailType.String(),
			ActualType:   reflect.TypeOf(detail).String(),
		}))
	}
	return &ex{
		code:   code,
		detail: detail,
	}
}

// Is checks if the errorex is of type EX and if the code matches
func Is(err error, code string) bool {
	// Check if the error code is registered
	if _, ok := errorCodes[code]; !ok {
		// Fatal errorex
		panic(New(ErrCodeNotRegistered, ErrorEXDetail{Code: code}))
	}
	// Check if the error is nil
	if err == nil {
		return false
	}
	// Check if the error has a method Code
	errorValue := reflect.ValueOf(err)
	if errorValue.Kind() != reflect.Ptr {
		return false
	}
	codeMethod := errorValue.MethodByName("Code")
	if !codeMethod.IsValid() {
		return false
	}
	codeValue := codeMethod.Call([]reflect.Value{})
	if len(codeValue) != 1 {
		return false
	}
	if codeValue[0].Kind() != reflect.String {
		return false
	}
	return codeValue[0].String() == code
}
