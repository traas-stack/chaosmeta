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

package errutil

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmeta/chaosmetad/pkg/log"
	"os"
)

const (
	NoErr = iota
	BadArgsErr
	DBErr
	InjectErr
	InternalErr
	RecoverErr
	UnknownErr
)

const (
	ExpectedErr = 99
)

func SolveErr(ctx context.Context, code int, msg string) {
	if code == NoErr {
		log.GetLogger(ctx).Debug(msg)
		os.Exit(NoErr)
	}
	log.GetLogger(ctx).Error(msg)
	os.Exit(code)
}

func ExitExpectedErr(msg string) {
	fmt.Printf("[error]%s\n", msg)
	os.Exit(ExpectedErr)
}
