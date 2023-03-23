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
	"reflect"
	"runtime"
	"strings"

	json "github.com/bytedance/sonic"
)

const (
	defaultDumpJSONValue = ""
)

// DumpJSON marshals a value to a JSON string, if it fails, it will return an empty string.
func DumpJSON(v any) string {
	b, err := json.MarshalString(v)
	if err != nil {
		return defaultDumpJSONValue
	}
	return b
}

// TernaryOp is a helper function to do a ternary operation.
func TernaryOp[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

// GetFunctionName uses reflect to retrieve the name of the function passed into the higher-order function.
func GetFunctionName(f any) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	if strings.Contains(name, "/") {
		v := strings.Split(name, "/")
		if len(v) == 0 {
			return name
		}
		return v[len(v)-1]
	}
	return name
}
