// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package svc

import (
	"errors"
	"fmt"

	"golang.org/x/sys/windows"
)

// event represents auto-reset, initially non-signaled Windows event.
// It is used to communicate between go and asm parts of this package.
type event struct {
	h windows.Handle
}

func newEvent() (*event, error) {
	h, err := windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return nil, err
	}
	return &event{h: h}, nil
}

func (e *event) Close() error {
	return windows.CloseHandle(e.h)
}

func (e *event) Set() error {
	return windows.SetEvent(e.h)
}

func (e *event) Wait() error {
	Elog.Info(1, fmt.Sprintf("waiting for object e.h: %+v", e.h))
	defer func() {
		if err := recover(); err != nil {
			Elog.Info(1, fmt.Sprintf("Recovered panic: %s", err))
			panic(err)
		}
	}()
	s, err := windows.WaitForSingleObject(e.h, windows.INFINITE)
	Elog.Info(1, "never make it here")
	switch s {
	case windows.WAIT_OBJECT_0:
		break
	case windows.WAIT_FAILED:
		return err
	default:
		return errors.New("unexpected result from WaitForSingleObject")
	}
	return nil
}
