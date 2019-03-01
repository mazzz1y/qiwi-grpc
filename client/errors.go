// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package client

import "fmt"

var (
	// ErrSyntaxError error with bad request
	ErrSyntaxError = fmt.Errorf("invalid syntax")
	// ErrTokenInvalid error
	ErrTokenInvalid = fmt.Errorf("invalid token")
	// ErrNotFound error
	ErrNotFound = fmt.Errorf("not found")
	// ErrRateLimitReached error
	ErrRateLimitReached = fmt.Errorf("too many requests")
	// ErrForbidden error
	ErrIncorrectWebHookParams = fmt.Errorf("incorrect webhook params")
	// ErrForbidden error
	ErrForbidden = fmt.Errorf("forbidden")
	// ErrServerError error
	ErrServerError = fmt.Errorf("server error")
)

var (
	codeToError = map[int]error{
		400: ErrSyntaxError,
		401: ErrTokenInvalid,
		403: ErrForbidden,
		404: ErrNotFound,
		422: ErrIncorrectWebHookParams,
		423: ErrRateLimitReached,
		500: ErrServerError,
	}
)
