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

// Set is a generic set implementation.
type Set[T comparable] map[T]struct{}

// New returns a new set based on the specified generic type (must be comparable).
func New[T comparable]() Set[T] {
	return Set[T]{}
}

// Add adds an item to the set. It returns the same set to allow easy chaining of Adds.
func (s Set[T]) Add(item T) Set[T] {
	s[item] = struct{}{}
	return s
}

// Remove removes an item from the set. It will throw an error if the item is not present.
func (s Set[T]) Remove(item T) bool {
	if !s.Contains(item) {
		return false
	}
	delete(s, item)
	return true
}

// MustRemove removes an item from the set. However, it assumes that the item is present and it must be removed.
// It will panic if the item is not present.
func (s Set[T]) MustRemove(item T) {
	ok := s.Remove(item)
	if !ok {
		panic("failed to remove from set")
	}
}

// Contains simply checks if the item exists in the set.
func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}
