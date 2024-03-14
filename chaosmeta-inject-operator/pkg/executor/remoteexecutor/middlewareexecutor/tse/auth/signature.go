package auth

import (
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/config"
	"gitlab.alipay-inc.com/mist-sdk/mist_sdk_go/mist"
	"log"
)

type MistClient struct {
	Config config.MistConfig
}

func (r *MistClient) MistConfig() (string, string) {
	config := mist.NewMistConfig()
	// 线下/线上AntVipUrl和BkmiUrl具体配置参看上面的 参数说明
	config.SetAntVipUrl(r.Config.AntVipUrl)
	config.SetBkmiUrl(r.Config.BkmiUrl)

	config.SetAppName(r.Config.AppName)
	config.SetTenant(r.Config.Tenant)

	// 线下/线上Mode具体配置参看上面的 参数说明
	config.SetMode(r.Config.Mode)

	client := mist.NewMistClient(config)
	ak, sk, _, _, err := client.GetSecretInfo(r.Config.SecretName)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	return ak, sk
}
