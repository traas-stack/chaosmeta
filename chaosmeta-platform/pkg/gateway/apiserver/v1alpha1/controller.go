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

package v1alpha1

import (
	"chaosmeta-platform/util/errors"
	"chaosmeta-platform/util/log"
	beego "github.com/beego/beego/v2/server/web"
)

type BeegoOutputController struct{}

func (c BeegoOutputController) Error(bc *beego.Controller, err error) {
	log.Error(err)
	bc.Data["json"] = errors.ErrServer().WithMessage(err.Error())
	bc.ServeJSON()
}

func (c BeegoOutputController) ErrUnauthorized(bc *beego.Controller, err error) {
	log.Error(err)
	bc.Data["json"] = errors.ErrUnauthorized().WithMessage(err.Error())
	bc.ServeJSON()
}

func (c BeegoOutputController) ErrorWithMessage(bc *beego.Controller, msg string) {
	log.Error(msg)
	bc.Data["json"] = errors.ErrServer().WithMessage(errors.ErrServer().WithMessage(msg).Error())
	bc.ServeJSON()
}

func (c BeegoOutputController) ErrorWithData(bc *beego.Controller, data interface{}) {
	bc.Data["json"] = errors.ErrServer().WithData(data)
	bc.ServeJSON()
}

func (c BeegoOutputController) Success(bc *beego.Controller, data interface{}) {
	bc.Data["json"] = errors.OK().WithData(data)
	bc.ServeJSON()
}

func (c BeegoOutputController) SuccessNoData(bc *beego.Controller) {
	bc.Data["json"] = errors.OK().CleanData()
	bc.ServeJSON()
}
