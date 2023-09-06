package internal

import (
	"github.com/coffeehc/boot/configuration"
)

func GetServiceInfo() configuration.ServiceInfo {
	return configuration.ServiceInfo{
		ServiceName: "xNetForward",
		Version:     "0.0.1",
		Descriptor:  "网络转发",
		APIDefine:   "",
	}
}
