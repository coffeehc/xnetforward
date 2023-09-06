package forwardservice

import (
	"errors"
	"github.com/coffeehc/base/log"
	"github.com/coffeehc/commons/asyncservice"
	"github.com/coffeehc/xnetforward/internal/services/configservice"
	"go.uber.org/zap"
	"io"
	"net"
	"time"
)

func newListen(forward *configservice.Forward) *Listener {
	impl := &Listener{
		forward: forward,
	}
	return impl
}

type Listener struct {
	forward *configservice.Forward
	lister  net.Listener
	isClose bool
}

func (impl *Listener) Start() error {
	lister, err := net.Listen(impl.forward.Network, impl.forward.Src)
	if err != nil {
		return err
	}
	impl.lister = lister
	asyncService := asyncservice.GetService()
	asyncService.Submit(func() {
		log.Debug("开始监听端口", zap.String("src", impl.forward.Src), zap.String("network", impl.forward.Network))
		for !impl.isClose {
			srcConn, err1 := lister.Accept()
			if err1 != nil {
				if errors.Is(err1, net.ErrClosed) {
					return
				}
				continue
			}
			asyncService.Submit(func() {
				targetConn, err2 := net.DialTimeout(impl.forward.Network, impl.forward.Target, time.Second*5)
				if err2 != nil {
					srcConn.Close()
					return
				}
				asyncService.Submit(func() {
					io.Copy(targetConn, srcConn)
				})
				asyncService.Submit(func() {
					io.Copy(srcConn, targetConn)
				})
			})

		}
	})
	return nil
}

func (impl *Listener) Close() {
	impl.isClose = true
	impl.lister.Close()
	log.Debug("关闭监听端口", zap.String("src", impl.forward.Src), zap.String("network", impl.forward.Network))
}
