package common

import (
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-plugins/config/source/consul/v2"
	"strconv"
)

func GetConsulConfig(host string, port int64, prefix string) (config.Config, error) {
	consulSource := consul.NewSource(
		// 设置配置中心的地址
		consul.WithAddress(host+":"+strconv.FormatInt(port, 10)),
		// 设置前缀
		consul.WithPrefix(prefix),
		// 是否移除前缀，这里设置为true
		consul.StripPrefix(true),
	)
	config, err := config.NewConfig()
	if err != nil {
		return config, err
	}
	err = config.Load(consulSource)
	return config, err
}
