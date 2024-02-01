package logger

import "github.com/donething/utils-go/dolog"

// openwrt 中可能缺少`zoneinfo`包，时区就会为`UTC`时区
var Info, Warn, Error = dolog.InitLog(dolog.DefaultFlag)
