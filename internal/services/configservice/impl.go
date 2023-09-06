package configservice

import (
	"context"
	"github.com/coffeehc/base/log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
)

type Service interface {
	RegisterOnConfigChange(name string, f func())
	GetForwards() []*Forward
}

func newService(ctx context.Context) Service {
	impl := &serviceImpl{
		onChanges: make(map[string]func()),
	}
	impl.loadConfig()
	viper.OnConfigChange(func(_ fsnotify.Event) {
		impl.loadConfig()
	})
	return impl
}

type serviceImpl struct {
	config       *Config
	onChanges    map[string]func()
	onChangeLock sync.Mutex
}

func (impl *serviceImpl) GetForwards() []*Forward {
	return impl.config.Forwards
}

func (impl *serviceImpl) RegisterOnConfigChange(name string, f func()) {
	impl.onChanges[name] = f
}

func (impl *serviceImpl) loadConfig() {
	impl.onChangeLock.Lock()
	defer impl.onChangeLock.Unlock()
	config := &Config{}
	err := viper.Unmarshal(config)
	if err == nil {
		impl.config = config
	}
	for _, handler := range impl.onChanges {
		func() {
			defer func() {
				if e := recover(); e != nil {
					log.DPanic("不可处理的异常", zap.Any("error", e))
				}
			}()
			handler()
		}()
	}
}

func (impl *serviceImpl) Start(ctx context.Context) error {
	return nil
}

func (impl *serviceImpl) Stop(ctx context.Context) error {
	return nil
}
