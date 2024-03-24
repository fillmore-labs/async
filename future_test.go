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

package async_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"fillmore-labs.com/async"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

var errTest = errors.New("test error")

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestAsyncValue(t *testing.T) {
	t.Parallel()

	// given
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// when
	f := async.NewAsync(func() (int, error) { return 1, nil })
	val, err := f.Await(ctx)

	// then
	if assert.NoError(t, err) {
		assert.Equal(t, 1, val)
	}
}

func TestAsyncError(t *testing.T) {
	t.Parallel()

	// given
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// when
	f := async.NewAsync(func() (int, error) { return 0, errTest })
	_, err := f.Await(ctx)

	// then
	assert.ErrorIs(t, err, errTest)
}

func TestCancellation(t *testing.T) {
	t.Parallel()

	// given
	run := make(chan struct{})
	f := async.NewAsync(func() (int, error) {
		<-run

		return 1, nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// when
	_, err := f.Await(ctx)
	close(run)

	// then
	assert.ErrorIs(t, err, context.Canceled)
}

func TestMultiple(t *testing.T) {
	t.Parallel()

	// given
	const iterations = 1_000
	const concurrency = 4

	// when
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for i := 0; i < iterations; i++ {
		p := async.Promise[int]{}
		f := p.Future()

		var values [concurrency]int
		var errs [concurrency]error

		var wg sync.WaitGroup
		wg.Add(concurrency)
		for c := 0; c < concurrency; c++ {
			go func(i int) {
				defer wg.Done()
				values[i], errs[i] = f.Await(ctx)
			}(c)
		}
		p.Resolve(i)
		wg.Wait()

		// then
		for c := 0; c < concurrency; c++ {
			if assert.NoError(t, errs[c]) {
				assert.Equal(t, i, values[c])
			}
		}
	}
}

func TestTry(t *testing.T) {
	t.Parallel()

	// given
	p := async.Promise[int]{}
	f := p.Future()

	// when
	_, err1 := f.Try()
	_ = f.Done()
	_, err2 := f.Try()

	p.Resolve(1)
	value3, err3 := f.Try()
	value4, err4 := f.Try()

	// then
	assert.ErrorIs(t, err1, async.ErrNotReady)
	assert.ErrorIs(t, err2, async.ErrNotReady)
	if assert.NoError(t, err3) {
		assert.Equal(t, 1, value3)
	}
	if assert.NoError(t, err4) {
		assert.Equal(t, 1, value4)
	}
}

func TestNil(t *testing.T) {
	t.Parallel()

	// given
	var p *async.Promise[int]
	f := p.Future()

	// when
	p.Resolve(1)
	p.Reject(errTest)
	p.Do(func() (int, error) { panic("should not be called") })

	_, err := f.Try()
	done := f.Done()

	// then
	assert.ErrorIs(t, err, async.ErrNotReady)
	assert.Nil(t, done)
}

func TestPromise_String(t *testing.T) {
	t.Parallel()

	pending := func() *async.Promise[int] {
		p := async.Promise[int]{}

		return &p
	}

	pending2 := func() *async.Promise[int] {
		p := async.Promise[int]{}
		_ = p.Future().Done()

		return &p
	}

	resolved := func() *async.Promise[int] {
		p := async.Promise[int]{}
		p.Resolve(1)

		return &p
	}

	tests := []struct {
		name string
		p    *async.Promise[int]
		want string
	}{
		{"Nil", (*async.Promise[int])(nil), "Promise <nil>"},
		{"Pending", pending(), "Promise pending"},
		{"PendingWithDone", pending2(), "Promise pending"},
		{"Resolved", resolved(), "Promise resolved: 1, <nil>"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equalf(t, tt.want, tt.p.String(), "String()")
		})
	}
}

func TestFuture_String(t *testing.T) {
	t.Parallel()

	pending := func() *async.Future[int] {
		p := async.Promise[int]{}

		return p.Future()
	}

	pending2 := func() *async.Future[int] {
		p := async.Promise[int]{}
		_ = p.Future().Done()

		return p.Future()
	}

	resolved := func() *async.Future[int] {
		p := async.Promise[int]{}
		p.Reject(errTest)

		return p.Future()
	}

	tests := []struct {
		name string
		f    *async.Future[int]
		want string
	}{
		{"Nil", (*async.Future[int])(nil), "Future <nil>"},
		{"Pending", pending(), "Future pending"},
		{"PendingWithDone", pending2(), "Future pending"},
		{"Resolved", resolved(), "Future resolved: 0, test error"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equalf(t, tt.want, tt.f.String(), "String()")
		})
	}
}
