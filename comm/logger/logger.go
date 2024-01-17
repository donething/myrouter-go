package logger

import "github.com/donething/utils-go/dolog"

var Info, Warn, Error = dolog.InitLog(dolog.DefaultFlag)
