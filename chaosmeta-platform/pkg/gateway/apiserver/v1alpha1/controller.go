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
