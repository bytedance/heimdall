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
	"fmt"
	"time"

	json "github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

// GenerateCacheKey generates a cache key for a given function name and request encoded in SHA512. With the following format:
// SHA512(functionName:marshalledRequest:softTTL:hardTTL)
func GenerateCacheKey(req any, functionName string, softTTL, hardTTL time.Duration, version string) (string, error) {
	var (
		err           error
		marshalledReq string
	)

	if req == nil {
		marshalledReq = ""
	} else {
		marshalledReq, err = json.ConfigStd.MarshalToString(req) // ensures Map's keys are sorted for unique key generation
		if err != nil {
			return "", errors.Wrap(err, "")
		}
	}

	key := constructUnhashedKey(functionName, marshalledReq, softTTL, hardTTL, version)

	hasher := sha512.New()
	hasher.Write([]byte(key))
	hashedKey := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return hashedKey, nil
}

func constructUnhashedKey(functionName string, marshalledReq string, softTTL, hardTTL time.Duration, version string) string {
	return fmt.Sprintf("%v:%v:%d:%d:%v", functionName, marshalledReq, int64(softTTL.Seconds()), int64(hardTTL.Seconds()), version)
}
