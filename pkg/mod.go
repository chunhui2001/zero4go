package pkg

import (
	"os"
	"os/signal"
	"syscall"

	. "github.com/chunhui2001/zero4go/pkg/logs"
)

var shutdownSignals chan os.Signal
var funcChannel chan func()
var isStarted bool = false
var funcArray []func()
var done chan struct{}

func AddShutDownHook(f func()) {
	if !isStarted {
		shutdownSignals = make(chan os.Signal, 1)
		funcChannel = make(chan func())
		done = make(chan struct{})
		signal.Notify(shutdownSignals, syscall.SIGINT, syscall.SIGTERM)

		start()

		isStarted = true
	}

	funcChannel <- f
}

func WaitShutDown() {
	<-done
}

func executeHooks() {
	for _, f := range funcArray {
		f()
	}
}

func start() {
	go func() {
		shutdown := false

		for !shutdown {
			select {
			case <-shutdownSignals:
				for signal := range shutdownSignals {
					if signal == syscall.SIGTERM || signal == syscall.SIGQUIT {
						Log.Infof("kill -15 退出进程: signal=%v", signal)
						break
					} else if signal == syscall.SIGILL {
						Log.Infof("kill -4 退出进程: signal=%v", signal)
						break
					} else {
						Log.Infof("kill -? 退出进程: signal=%v", signal)
						break
					}
				}

				executeHooks()

				shutdown = true
			case f := <-funcChannel:
				funcArray = append(funcArray, f)
			}
		}

		done <- struct{}{}
	}()
}
