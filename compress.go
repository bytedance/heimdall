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
	"bytes"
	"compress/gzip"
	"context"

	"github.com/bytedance/heimdall/constants"
	json "github.com/bytedance/sonic"
	"github.com/golang/snappy"
	"github.com/pkg/errors"
)

const (
	defaultCompressionLevel = gzip.DefaultCompression
)

// CompressStruct converts a struct to JSON and compresses it according to chosen compression library.
func CompressStruct(ctx context.Context, s any, compressionLibrary constants.CompressionLibraryType) ([]byte, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal struct")
	}
	switch compressionLibrary {
	case constants.GzipCompressionType:
		return gzipCompression(b)
	case constants.SnappyCompressionType:
		return snappyCompression(b)
	default:
		return b, nil
	}
}

func gzipCompression(data []byte) ([]byte, error) {
	var buffer bytes.Buffer

	writer, err := gzip.NewWriterLevel(&buffer, defaultCompressionLevel)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create new gzip writer")
	}

	_, err = writer.Write(data)
	if err != nil {
		return nil, errors.Wrap(err, "cannot write to gzip")
	}

	err = writer.Close()
	if err != nil {
		return nil, errors.Wrap(err, "cannot close writer")
	}

	return buffer.Bytes(), nil
}

func snappyCompression(data []byte) ([]byte, error) {
	var buffer bytes.Buffer

	writer := snappy.NewBufferedWriter(&buffer)

	_, err := writer.Write(data)

	if err != nil {
		return nil, errors.Wrap(err, "cannot write to snappy")
	}

	err = writer.Close()
	if err != nil {
		return nil, errors.Wrap(err, "cannot close writer")
	}

	return buffer.Bytes(), nil
}
