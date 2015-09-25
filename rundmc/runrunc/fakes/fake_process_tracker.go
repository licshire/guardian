// This file was generated by counterfeiter
package fakes

import (
	"os/exec"
	"sync"

	"github.com/cloudfoundry-incubator/garden"
	"github.com/cloudfoundry-incubator/guardian/rundmc/process_tracker"
	"github.com/cloudfoundry-incubator/guardian/rundmc/runrunc"
)

type FakeProcessTracker struct {
	RunStub        func(id uint32, cmd *exec.Cmd, io garden.ProcessIO, tty *garden.TTYSpec, signaller process_tracker.Signaller) (garden.Process, error)
	runMutex       sync.RWMutex
	runArgsForCall []struct {
		id        uint32
		cmd       *exec.Cmd
		io        garden.ProcessIO
		tty       *garden.TTYSpec
		signaller process_tracker.Signaller
	}
	runReturns struct {
		result1 garden.Process
		result2 error
	}
}

func (fake *FakeProcessTracker) Run(id uint32, cmd *exec.Cmd, io garden.ProcessIO, tty *garden.TTYSpec, signaller process_tracker.Signaller) (garden.Process, error) {
	fake.runMutex.Lock()
	fake.runArgsForCall = append(fake.runArgsForCall, struct {
		id        uint32
		cmd       *exec.Cmd
		io        garden.ProcessIO
		tty       *garden.TTYSpec
		signaller process_tracker.Signaller
	}{id, cmd, io, tty, signaller})
	fake.runMutex.Unlock()
	if fake.RunStub != nil {
		return fake.RunStub(id, cmd, io, tty, signaller)
	} else {
		return fake.runReturns.result1, fake.runReturns.result2
	}
}

func (fake *FakeProcessTracker) RunCallCount() int {
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	return len(fake.runArgsForCall)
}

func (fake *FakeProcessTracker) RunArgsForCall(i int) (uint32, *exec.Cmd, garden.ProcessIO, *garden.TTYSpec, process_tracker.Signaller) {
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	return fake.runArgsForCall[i].id, fake.runArgsForCall[i].cmd, fake.runArgsForCall[i].io, fake.runArgsForCall[i].tty, fake.runArgsForCall[i].signaller
}

func (fake *FakeProcessTracker) RunReturns(result1 garden.Process, result2 error) {
	fake.RunStub = nil
	fake.runReturns = struct {
		result1 garden.Process
		result2 error
	}{result1, result2}
}

var _ runrunc.ProcessTracker = new(FakeProcessTracker)