package log

import (
	"chaosmeta-platform/config"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
)

const TraceIdKey = "TraceId"

func Init() {
	if err := logs.SetLogger(logs.AdapterFile, fmt.Sprintf(`{"filename":"%s","daily":%t,"maxdays":%d}`, config.DefaultRunOptIns.Log.Filename, config.DefaultRunOptIns.Log.Daily, config.DefaultRunOptIns.Log.MaxDays)); err != nil {
		panic(any(fmt.Sprintf("set logger error: %s", err.Error())))
	}
}
