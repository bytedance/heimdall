/*
 * Copyright 2022 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddSet(t *testing.T) {
	set := New[int64]()
	tests := []struct {
		name      string
		itemToAdd int64
		noInSet   int
	}{
		{
			name:      "add pos",
			itemToAdd: 1,
			noInSet:   1,
		}, {
			name:      "add neg",
			itemToAdd: -1,
			noInSet:   2,
		}, {
			name:      "add same",
			itemToAdd: 1,
			noInSet:   2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			set.Add(test.itemToAdd)
			assert.Equal(t, test.noInSet, len(set))
		})
	}
}

func TestRemoveSet(t *testing.T) {
	set := New[int64]().Add(1).Add(-1)
	tests := []struct {
		name         string
		itemToRemove int64
		noInSet      int
		ok           bool
	}{
		{
			name:         "remove pos",
			itemToRemove: 1,
			noInSet:      1,
			ok:           true,
		}, {
			name:         "remove neg",
			itemToRemove: -1,
			noInSet:      0,
			ok:           true,
		}, {
			name:         "remove zero",
			itemToRemove: 1,
			noInSet:      0,
			ok:           false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ok := set.Remove(test.itemToRemove)
			assert.Equal(t, test.ok, ok)
			assert.Equal(t, test.noInSet, len(set))
		})
	}
}
