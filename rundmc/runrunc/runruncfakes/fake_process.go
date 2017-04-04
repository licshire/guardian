// This file was generated by counterfeiter
package runruncfakes

import (
	"sync"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/guardian/rundmc/runrunc"
)

type FakeProcess struct {
	IDStub        func() string
	iDMutex       sync.RWMutex
	iDArgsForCall []struct{}
	iDReturns     struct {
		result1 string
	}
	iDReturnsOnCall map[int]struct {
		result1 string
	}
	WaitStub        func() (int, error)
	waitMutex       sync.RWMutex
	waitArgsForCall []struct{}
	waitReturns     struct {
		result1 int
		result2 error
	}
	waitReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	ExitStatusStub        func() chan garden.ProcessStatus
	exitStatusMutex       sync.RWMutex
	exitStatusArgsForCall []struct{}
	exitStatusReturns     struct {
		result1 chan garden.ProcessStatus
	}
	exitStatusReturnsOnCall map[int]struct {
		result1 chan garden.ProcessStatus
	}
	SetTTYStub        func(garden.TTYSpec) error
	setTTYMutex       sync.RWMutex
	setTTYArgsForCall []struct {
		arg1 garden.TTYSpec
	}
	setTTYReturns struct {
		result1 error
	}
	setTTYReturnsOnCall map[int]struct {
		result1 error
	}
	SignalStub        func(garden.Signal) error
	signalMutex       sync.RWMutex
	signalArgsForCall []struct {
		arg1 garden.Signal
	}
	signalReturns struct {
		result1 error
	}
	signalReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeProcess) ID() string {
	fake.iDMutex.Lock()
	ret, specificReturn := fake.iDReturnsOnCall[len(fake.iDArgsForCall)]
	fake.iDArgsForCall = append(fake.iDArgsForCall, struct{}{})
	fake.recordInvocation("ID", []interface{}{})
	fake.iDMutex.Unlock()
	if fake.IDStub != nil {
		return fake.IDStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.iDReturns.result1
}

func (fake *FakeProcess) IDCallCount() int {
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	return len(fake.iDArgsForCall)
}

func (fake *FakeProcess) IDReturns(result1 string) {
	fake.IDStub = nil
	fake.iDReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeProcess) IDReturnsOnCall(i int, result1 string) {
	fake.IDStub = nil
	if fake.iDReturnsOnCall == nil {
		fake.iDReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.iDReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeProcess) Wait() (int, error) {
	fake.waitMutex.Lock()
	ret, specificReturn := fake.waitReturnsOnCall[len(fake.waitArgsForCall)]
	fake.waitArgsForCall = append(fake.waitArgsForCall, struct{}{})
	fake.recordInvocation("Wait", []interface{}{})
	fake.waitMutex.Unlock()
	if fake.WaitStub != nil {
		return fake.WaitStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.waitReturns.result1, fake.waitReturns.result2
}

func (fake *FakeProcess) WaitCallCount() int {
	fake.waitMutex.RLock()
	defer fake.waitMutex.RUnlock()
	return len(fake.waitArgsForCall)
}

func (fake *FakeProcess) WaitReturns(result1 int, result2 error) {
	fake.WaitStub = nil
	fake.waitReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeProcess) WaitReturnsOnCall(i int, result1 int, result2 error) {
	fake.WaitStub = nil
	if fake.waitReturnsOnCall == nil {
		fake.waitReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.waitReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeProcess) ExitStatus() chan garden.ProcessStatus {
	fake.exitStatusMutex.Lock()
	ret, specificReturn := fake.exitStatusReturnsOnCall[len(fake.exitStatusArgsForCall)]
	fake.exitStatusArgsForCall = append(fake.exitStatusArgsForCall, struct{}{})
	fake.recordInvocation("ExitStatus", []interface{}{})
	fake.exitStatusMutex.Unlock()
	if fake.ExitStatusStub != nil {
		return fake.ExitStatusStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.exitStatusReturns.result1
}

func (fake *FakeProcess) ExitStatusCallCount() int {
	fake.exitStatusMutex.RLock()
	defer fake.exitStatusMutex.RUnlock()
	return len(fake.exitStatusArgsForCall)
}

func (fake *FakeProcess) ExitStatusReturns(result1 chan garden.ProcessStatus) {
	fake.ExitStatusStub = nil
	fake.exitStatusReturns = struct {
		result1 chan garden.ProcessStatus
	}{result1}
}

func (fake *FakeProcess) ExitStatusReturnsOnCall(i int, result1 chan garden.ProcessStatus) {
	fake.ExitStatusStub = nil
	if fake.exitStatusReturnsOnCall == nil {
		fake.exitStatusReturnsOnCall = make(map[int]struct {
			result1 chan garden.ProcessStatus
		})
	}
	fake.exitStatusReturnsOnCall[i] = struct {
		result1 chan garden.ProcessStatus
	}{result1}
}

func (fake *FakeProcess) SetTTY(arg1 garden.TTYSpec) error {
	fake.setTTYMutex.Lock()
	ret, specificReturn := fake.setTTYReturnsOnCall[len(fake.setTTYArgsForCall)]
	fake.setTTYArgsForCall = append(fake.setTTYArgsForCall, struct {
		arg1 garden.TTYSpec
	}{arg1})
	fake.recordInvocation("SetTTY", []interface{}{arg1})
	fake.setTTYMutex.Unlock()
	if fake.SetTTYStub != nil {
		return fake.SetTTYStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setTTYReturns.result1
}

func (fake *FakeProcess) SetTTYCallCount() int {
	fake.setTTYMutex.RLock()
	defer fake.setTTYMutex.RUnlock()
	return len(fake.setTTYArgsForCall)
}

func (fake *FakeProcess) SetTTYArgsForCall(i int) garden.TTYSpec {
	fake.setTTYMutex.RLock()
	defer fake.setTTYMutex.RUnlock()
	return fake.setTTYArgsForCall[i].arg1
}

func (fake *FakeProcess) SetTTYReturns(result1 error) {
	fake.SetTTYStub = nil
	fake.setTTYReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeProcess) SetTTYReturnsOnCall(i int, result1 error) {
	fake.SetTTYStub = nil
	if fake.setTTYReturnsOnCall == nil {
		fake.setTTYReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setTTYReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeProcess) Signal(arg1 garden.Signal) error {
	fake.signalMutex.Lock()
	ret, specificReturn := fake.signalReturnsOnCall[len(fake.signalArgsForCall)]
	fake.signalArgsForCall = append(fake.signalArgsForCall, struct {
		arg1 garden.Signal
	}{arg1})
	fake.recordInvocation("Signal", []interface{}{arg1})
	fake.signalMutex.Unlock()
	if fake.SignalStub != nil {
		return fake.SignalStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.signalReturns.result1
}

func (fake *FakeProcess) SignalCallCount() int {
	fake.signalMutex.RLock()
	defer fake.signalMutex.RUnlock()
	return len(fake.signalArgsForCall)
}

func (fake *FakeProcess) SignalArgsForCall(i int) garden.Signal {
	fake.signalMutex.RLock()
	defer fake.signalMutex.RUnlock()
	return fake.signalArgsForCall[i].arg1
}

func (fake *FakeProcess) SignalReturns(result1 error) {
	fake.SignalStub = nil
	fake.signalReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeProcess) SignalReturnsOnCall(i int, result1 error) {
	fake.SignalStub = nil
	if fake.signalReturnsOnCall == nil {
		fake.signalReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.signalReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeProcess) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	fake.waitMutex.RLock()
	defer fake.waitMutex.RUnlock()
	fake.exitStatusMutex.RLock()
	defer fake.exitStatusMutex.RUnlock()
	fake.setTTYMutex.RLock()
	defer fake.setTTYMutex.RUnlock()
	fake.signalMutex.RLock()
	defer fake.signalMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeProcess) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ runrunc.Process = new(FakeProcess)
