package mnnet

import (
	"context"
	"mn-utils/logger"
	"net"
	"syscall"
)

func Listen(network, address string) (net.Listener, error) {
	var lc net.ListenConfig

	lc.Control = func(network, address string, c syscall.RawConn) error {
		return c.Control(func(fd uintptr) {
			//上面都是默认写法，这里是设置地址复用
			err := syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, ^syscall.SO_REUSEADDR, 1)
			if err != nil {
				logger.Errorln(err)
			}
		})
	}

	return lc.Listen(context.Background(), network, address)
}
