package internal

import "sync"

type SingleInstModuleCore struct {
	sync.RWMutex
}
