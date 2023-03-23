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

package heimdall

import (
	"context"
	"testing"

	"github.com/bytedance/heimdall/constants"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Foo   string
	World int64
}

func TestStructGzip(t *testing.T) {
	ctx := context.Background()
	populatedStruct := &testStruct{Foo: "Bar", World: 42}
	compressedData, _ := CompressStruct(ctx, populatedStruct, constants.GzipCompressionType)

	newStruct := &testStruct{}
	err := DecompressStruct(ctx, compressedData, newStruct, constants.GzipCompressionType)
	assert.Nil(t, err)
	assert.Equal(t, populatedStruct, newStruct)
}

func TestStructSnappy(t *testing.T) {
	ctx := context.Background()
	populatedStruct := &testStruct{Foo: "Bar", World: 42}
	compressedData, _ := CompressStruct(ctx, populatedStruct, constants.SnappyCompressionType)

	newStruct := &testStruct{}
	err := DecompressStruct(ctx, compressedData, newStruct, constants.SnappyCompressionType)
	assert.Nil(t, err)
	assert.Equal(t, populatedStruct, newStruct)
}

func TestStructNoCompression(t *testing.T) {
	ctx := context.Background()
	populatedStruct := &testStruct{Foo: "Bar", World: 42}
	compressedData, _ := CompressStruct(ctx, populatedStruct, constants.NoCompressionType)

	newStruct := &testStruct{}
	err := DecompressStruct(ctx, compressedData, newStruct, constants.NoCompressionType)
	assert.Nil(t, err)
	assert.Equal(t, populatedStruct, newStruct)
}
