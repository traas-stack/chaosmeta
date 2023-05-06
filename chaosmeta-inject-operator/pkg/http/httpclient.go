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

package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type HTTPClient struct {
	Client *http.Client
}

func (h *HTTPClient) Post(ctx context.Context, url string, data []byte) ([]byte, error) {
	logger := log.FromContext(ctx)
	logger.Info("request: " + string(data))

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("new requset error: %s", err.Error())
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client send requset error: %s", err.Error())
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %s", err.Error())
	}

	logger.Info("response: " + string(res))
	return res, nil
}

func (h *HTTPClient) Get(ctx context.Context, url string) ([]byte, error) {
	logger := log.FromContext(ctx)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("new requset error: %s", err.Error())
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client send requset error: %s", err.Error())
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %s", err.Error())
	}

	logger.Info("response: " + string(res))
	return res, nil
}
