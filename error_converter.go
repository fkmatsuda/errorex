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

// ErrorConverter defines the interface for handlers in the chain of responsibility
// that attempt to convert an error into an error compatible with errorex.EX.
// If the handler cannot convert the error, it delegates the task to the next handler in the chain.
// A large number of handlers in the chain of responsibilities can cause delays in the process, so I do not recommend having a single chain for the entire application. Instead, have multiple chains with handlers that make sense in their contexts.
type ErrorConverter interface {
	// ConvertError tries to convert an error into another error compatible with errorex.EX.
	// If it cannot handle the conversion, it attempts to delegate the task to the next handler in the chain.
	// It returns an error which will be nil if conversion is not possible by any handler.
	ConvertError(err error) EX

	// SetNext sets the next handler in the chain.
	SetNext(next ErrorConverter)
}

// BaseErrorConverter provides a basic implementation of the ErrorConverter interface,
// holding a reference to the next handler in the chain.
type BaseErrorConverter struct {
	next ErrorConverter
}

// SetNext sets the next handler in the chain.
func (b *BaseErrorConverter) SetNext(next ErrorConverter) {
	b.next = next
}

// ConvertError in BaseErrorConverter should be overridden by concrete handlers.
// This default implementation delegates the conversion task to the next handler.
func (b *BaseErrorConverter) ConvertError(err error) EX {
	if b.next != nil {
		return b.next.ConvertError(err)
	}
	return nil // Return nil if there is no next handler to delegate to.
}

type unknownErrorConverter struct {
	BaseErrorConverter
}

func (u *unknownErrorConverter) ConvertError(err error) EX {
	return New(ErrCodeUnknownError, UnknownErrorDetail{Detail: err.Error()})
}

// NewUnknownErrorConverter creates a new unknownErrorConverter
// this converter will convert any error into an unknown error with the message of the error as the detail.
// This converter should be used as the last handler in the chain.
func NewUnknownErrorConverter() ErrorConverter {
	return &unknownErrorConverter{}
}

// exErrorConverter checks if the error passed implements EX, and if so, returns the error itself.
// Otherwise, it attempts to delegate the conversion to the next handler in the chain.
type exErrorConverter struct {
	BaseErrorConverter
}

// ConvertError checks if the error passed as a parameter implements EX.
// If so, it returns the error. If not, it attempts to delegate the conversion to the next handler.
func (c *exErrorConverter) ConvertError(err error) EX {
	if ex, ok := err.(EX); ok {
		return ex // Returns the parameter value if it is already an EX.
	}
	// Delegates to the next handler in the chain if this is not an EX error.
	return c.BaseErrorConverter.ConvertError(err)
}

// NewEXErrorConverter creates a new exErrorConverter
// this converter will check if the error passed implements EX, and if so, returns the error itself, otherwise it attempts to delegate the conversion to the next handler in the chain.
// This converter should be used as the first handler in the chain.
func NewEXErrorConverter(next ErrorConverter) ErrorConverter {
	converter := &exErrorConverter{}
	converter.SetNext(next)
	return converter
}
