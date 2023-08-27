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

package json

import (
	"encoding/json"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
)

func MargePatch(originalObj, updatedObj interface{}) ([]byte, error) {
	originalJSON, err := json.Marshal(originalObj)
	if err != nil {
		return nil, err
	}

	updatedJSON, err := json.Marshal(updatedObj)
	if err != nil {
		return nil, err
	}

	data, err := jsonpatch.CreateMergePatch(originalJSON, updatedJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to marge patch data, error: %s", err)
	}

	return data, nil
}
