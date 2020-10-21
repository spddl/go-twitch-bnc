package main

import (
	"bytes"
	"sync"
)

type MutexBacklog struct { // mutual-exclusion lock
	Mutex      sync.RWMutex
	msgBacklog map[string][]byte
}

func (v *MutexBacklog) len(channel string) int {
	v.Mutex.RLock()
	defer v.Mutex.RUnlock()
	return len(v.msgBacklog[channel])
}

func (v *MutexBacklog) count(channel string) int {
	v.Mutex.RLock()
	defer v.Mutex.RUnlock()
	return bytes.Count(v.msgBacklog[channel], []byte{13, 10})
}

func (v *MutexBacklog) getChannel(channel string) []byte {
	v.Mutex.RLock()
	defer v.Mutex.RUnlock()
	return v.msgBacklog[channel]
}

func (v *MutexBacklog) append(channel string, line []byte) {
	v.Mutex.Lock()
	if len(v.msgBacklog[channel]) == 0 {
		v.msgBacklog[channel] = append(line, []byte{13, 10}...)
	} else {
		v.msgBacklog[channel] = append(v.msgBacklog[channel], append(line, []byte{13, 10}...)...)
	}
	v.Mutex.Unlock()
}

func (v *MutexBacklog) remove(channel string, limit int) {
	v.Mutex.Lock()
	v.msgBacklog[channel] = v.msgBacklog[channel][:limit]
	v.Mutex.Unlock()
}

func (v *MutexBacklog) reset(channel string) {
	v.Mutex.Lock()
	v.msgBacklog[channel] = []byte{}
	v.Mutex.Unlock()
}
