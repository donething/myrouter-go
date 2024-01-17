package httpclient

import (
	"github.com/donething/utils-go/dohttp"
)

// Client http 客户端
var Client = dohttp.New(true, false)
