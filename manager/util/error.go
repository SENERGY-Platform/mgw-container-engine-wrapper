/*
 * Copyright 2022 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import "fmt"

type Error struct {
	code int
	msg  string
	err  error
}

func NewError(code int, msg string, err error) error {
	return &Error{
		code: code,
		msg:  msg,
		err:  err,
	}
}

func (e *Error) Error() string {
	if e.msg != "" && e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	} else if e.msg != "" {
		return e.msg
	} else {
		return e.err.Error()
	}
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Code() int {
	return e.code
}
