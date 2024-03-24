// Copyright 2023-2024 Oliver Eikemeier. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package async

import (
	"fmt"
	"sync/atomic"
)

// Promise defines the common operations for resolving a [Future] to its final value.
//
// An empty value is valid and must not be copied after creation. One of the resolving
// operations may be called once from any goroutine, all subsequent calls will panic.
type Promise[R any] struct {
	_     noCopy
	done  atomic.Pointer[chan struct{}] // lazy chan, signals when future has completed
	value R                             // result value, protected by done
	err   error                         // result error, protected by done
}

// Future returns a [Future] for this promise.
func (p *Promise[R]) Future() *Future[R] {
	return (*Future[R])(p)
}

// Resolve resolves the promise with a value.
func (p *Promise[R]) Resolve(value R) {
	if p == nil {
		return
	}

	p.value = value
	p.close()
}

// Reject breaks the promise with an error.
func (p *Promise[R]) Reject(err error) {
	if p == nil {
		return
	}

	p.err = err
	p.close()
}

// Do runs fn synchronously, fulfilling the promise once it completes.
func (p *Promise[R]) Do(fn func() (R, error)) {
	if p == nil {
		return
	}

	p.value, p.err = fn()
	p.close()
}

func (p *Promise[R]) close() {
	if done := p.done.Swap(&closedChan); done != nil {
		close(*done)
	}
}

func (p *Promise[R]) String() string {
	if p == nil {
		return "Promise <nil>"
	}

	if done := p.done.Load(); done != nil {
		select {
		case <-*done:
			return fmt.Sprintf("Promise resolved: %v, %v", p.value, p.err)
		default:
		}
	}

	return "Promise pending"
}
