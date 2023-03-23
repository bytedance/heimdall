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

package helpers

import (
	"crypto/sha512"
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCacheKey(t *testing.T) {
	tests := []struct {
		name         string
		req          *testReqStruct
		functionName string
		softTTL      time.Duration
		hardTTL      time.Duration
		key          string
		err          bool
		version      string
	}{
		{
			name: "normal",
			req: &testReqStruct{
				Foo: "bar",
			},
			functionName: "helloWorld",
			softTTL:      1 * time.Second,
			hardTTL:      2 * time.Second,
			key:          hash(`helloWorld:{"foo":"bar"}:1:2:v1.0.0`),
			err:          false,
			version:      "v1.0.0",
		},
		{
			name:         "nil req",
			req:          nil,
			functionName: "worldHello",
			softTTL:      2 * time.Second,
			hardTTL:      3 * time.Second,
			key:          hash("worldHello:null:2:3:v1.0.0"),
			err:          false,
			version:      "v1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := GenerateCacheKey(tt.req, tt.functionName, tt.softTTL, tt.hardTTL, tt.version)
			assert.Equal(t, tt.err, err != nil)
			assert.Equal(t, tt.key, key)
		})
	}
}

type testReqStruct struct {
	Foo string `json:"foo"`
}

func hash(v string) string {
	hasher := sha512.New()
	hasher.Write([]byte(v))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
