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

package models

import "time"

var tableNamePrefix = "chaosmeta_platform_"

type BaseTimeModel struct {
	CreateTime time.Time `json:"create_time" orm:"column(create_time);auto_now_add;type(datetime)"`
	UpdateTime time.Time `json:"update_time" orm:"column(update_time);auto_now;type(datetime)"`
}

func (bt *BaseTimeModel) SetCreateTime(time time.Time) {
	bt.CreateTime = time
}

func (bt *BaseTimeModel) SetUpdateTime(time time.Time) {
	bt.UpdateTime = time
}
