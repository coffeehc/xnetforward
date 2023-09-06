package forwardservice

import (
	"context"
	"github.com/coffeehc/base/log"
	"github.com/coffeehc/commons/asyncservice"
	"github.com/coffeehc/xnetforward/internal/services/configservice"
	"go.uber.org/zap"
)

type Service interface {
}

func newService(ctx context.Context) Service {
	configservice.EnablePlugin(ctx)
	configService := configservice.GetService()
	asyncservice.EnablePlugin(ctx)
	impl := &serviceImpl{
		configService: configService,
		forwards:      make(map[string]*Listener),
	}
	impl.configService.RegisterOnConfigChange("forward", impl.onConfigChange)
	return impl
}

type serviceImpl struct {
	configService configservice.Service
	forwards      map[string]*Listener
}

func (impl *serviceImpl) Start(ctx context.Context) error {
	forwards := impl.configService.GetForwards()
	for _, forward := range forwards {
		if forward.Src == "" || forward.Target == "" || forward.Network == "" {
			log.Error("forward定义错误", zap.Any("forward", forward))
			continue
		}
		listen := newListen(forward)
		impl.forwards[forward.Src] = listen
		err := listen.Start()
		if err != nil {
			log.Error("启动端口监听失败", zap.Error(err))
		}
	}
	return nil
}

func (impl *serviceImpl) Stop(ctx context.Context) error {
	return nil
}

func (impl *serviceImpl) onConfigChange() {
	forwards := impl.configService.GetForwards()
	closeListens := make([]string, 0)
	for key, _ := range impl.forwards {
		has := false
		for _, _forward := range forwards {
			if _forward.Src == key {
				has = true
			}
		}
		if !has {
			closeListens = append(closeListens, key)
		}
	}
	for _, key := range closeListens {
		impl.closeListen(key)
	}
	for _, forward := range forwards {
		if forward.Src == "" || forward.Target == "" || forward.Network == "" {
			log.Error("forward定义错误", zap.Any("forward", forward))
			continue
		}
		_, ok := impl.forwards[forward.Src]
		if ok {
			continue
		}
		listen := newListen(forward)
		impl.forwards[forward.Src] = listen
		err := listen.Start()
		if err != nil {
			log.Error("启动端口监听失败", zap.Error(err))
		}
	}
}

func (impl *serviceImpl) closeListen(key string) {
	listen := impl.forwards[key]
	listen.Close()
	delete(impl.forwards, key)
}
