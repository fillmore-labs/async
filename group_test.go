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
	"testing"

	"fillmore-labs.com/async"
	"fillmore-labs.com/async/group"
	"github.com/stretchr/testify/assert"
)

func TestDoAsync(t *testing.T) {
	t.Parallel()

	// given
	var g group.Group

	// when
	f := async.DoAsync(&g, func() (int, error) { return 0, errTest })

	err := g.Wait()
	_, err1 := f.Try()

	// then
	assert.ErrorIs(t, err, errTest)
	assert.ErrorIs(t, err1, errTest)
}
