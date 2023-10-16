/*
 * Copyright 2022-2023 Chaos Meta Authors.
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

package compress

import (
	"bytes"
	"compress/zlib"
	"io"
)

func DoZlibCompress(src string) (string, error) {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	if _, err := w.Write([]byte(src)); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return in.String(), nil
}

func DoZlibUnCompress(compressSrc string) (string, error) {
	b := bytes.NewReader([]byte(compressSrc))
	var out bytes.Buffer
	r, err := zlib.NewReader(b)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(&out, r); err != nil {
		return "", err
	}

	return out.String(), nil
}
