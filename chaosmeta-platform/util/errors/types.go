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

package errors

import (
	"net/http"
)

func OK() Error {
	return NewError(http.StatusOK, http.StatusText(http.StatusOK), 0)
}

func ErrServer() Error {
	return NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), 2)
}

func ErrParam() Error {
	return NewError(http.StatusBadRequest, "the user is already disabled", 2)
}

func ErrSignParam() Error {
	return NewError(http.StatusForbidden, http.StatusText(http.StatusForbidden), 2)
}

func ErrUnauthorized() Error {
	return NewError(http.StatusUnauthorized, "Unauthorized: Invalid username or password.", 2)
}

func ErrNotFound() Error {
	return NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound), 2)
}
