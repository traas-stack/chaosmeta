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

package user

import (
	"chaosmeta-platform/config"
	"fmt"
	"testing"
	"time"
)

func init() {
	config.DefaultRunOptIns = &config.Config{
		SecretKey: "samson",
	}
}

func TestAuthentication_VerifyToken(t *testing.T) {
	authentication := Authentication{}
	tocken, err := authentication.GenerateToken("samson", "admin", 1*time.Minute)
	fmt.Println(tocken)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := authentication.VerifyToken(tocken)
	fmt.Sprintln(*claims, err)
}
