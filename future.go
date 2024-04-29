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
	"context"
	"errors"
	"fmt"
)

// ErrNotReady is returned when the future is not complete.
var ErrNotReady = errors.New("future not ready")

// Future represents a read-only view of the result of an asynchronous operation.
//
// Construct a future either using [NewAsync] or from a [Promise] with [Promise.Future].
type Future[R any] Promise[R]

// NewAsync runs fn asynchronously, immediately returning a [Future] that can be used to retrieve the
// eventual result. This allows separating computation from evaluating the result.
func NewAsync[R any](fn func() (R, error)) *Future[R] {
	p := Promise[R]{}
	go p.Do(fn)

	return p.Future()
}

// Await returns the cached result or blocks until a result is available or the context is canceled.
func (f *Future[R]) Await(ctx context.Context) (R, error) {
	select { // wait for future completion or context cancel
	case <-f.Done():
		return f.value, f.err

	case <-ctx.Done():
		return *new(R), fmt.Errorf("future await: %w", ctx.Err())
	}
}

// Try returns the cached result when ready, [ErrNotReady] otherwise.
func (f *Future[R]) Try() (R, error) {
	if f == nil || !f.done.Closed() {
		return *new(R), ErrNotReady
	}

	return f.value, f.err
}

// Done returns a channel that is closed when the future is complete.
// It enables the use of future values in select statements.
func (f *Future[_]) Done() <-chan struct{} {
	if f == nil {
		return nil
	}

	return f.done.Done()
}

func (f *Future[R]) String() string {
	if f == nil {
		return "Future <nil>"
	}

	if !f.done.Closed() {
		return "Future pending"
	}

	return fmt.Sprintf("Future resolved: %v, %v", f.value, f.err)
}
