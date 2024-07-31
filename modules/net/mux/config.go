package mux

/*
   @Author: orbit-w
   @File: config
   @2024 7月 周日 19:23
*/

type ClientConfig struct {
	MaxVirtualConns uint32 //最大流数
}

const (
	maxVirtualConns = 200
)

func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		MaxVirtualConns: maxVirtualConns,
	}
}

func parseConfig(params ...ClientConfig) ClientConfig {
	if len(params) == 0 {
		return DefaultClientConfig()
	}

	conf := params[0]
	if conf.MaxVirtualConns <= 0 {
		conf.MaxVirtualConns = maxVirtualConns
	}
	return conf
}
