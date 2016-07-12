// This file was generated by counterfeiter
package runruncfakes

import (
	"sync"

	"code.cloudfoundry.org/lager"
	"github.com/cloudfoundry-incubator/guardian/rundmc/runrunc"
)

type FakeRuncCmdRunner struct {
	RunAndLogStub        func(log lager.Logger, cmd runrunc.LoggingCmd) error
	runAndLogMutex       sync.RWMutex
	runAndLogArgsForCall []struct {
		log lager.Logger
		cmd runrunc.LoggingCmd
	}
	runAndLogReturns struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeRuncCmdRunner) RunAndLog(log lager.Logger, cmd runrunc.LoggingCmd) error {
	fake.runAndLogMutex.Lock()
	fake.runAndLogArgsForCall = append(fake.runAndLogArgsForCall, struct {
		log lager.Logger
		cmd runrunc.LoggingCmd
	}{log, cmd})
	fake.recordInvocation("RunAndLog", []interface{}{log, cmd})
	fake.runAndLogMutex.Unlock()
	if fake.RunAndLogStub != nil {
		return fake.RunAndLogStub(log, cmd)
	} else {
		return fake.runAndLogReturns.result1
	}
}

func (fake *FakeRuncCmdRunner) RunAndLogCallCount() int {
	fake.runAndLogMutex.RLock()
	defer fake.runAndLogMutex.RUnlock()
	return len(fake.runAndLogArgsForCall)
}

func (fake *FakeRuncCmdRunner) RunAndLogArgsForCall(i int) (lager.Logger, runrunc.LoggingCmd) {
	fake.runAndLogMutex.RLock()
	defer fake.runAndLogMutex.RUnlock()
	return fake.runAndLogArgsForCall[i].log, fake.runAndLogArgsForCall[i].cmd
}

func (fake *FakeRuncCmdRunner) RunAndLogReturns(result1 error) {
	fake.RunAndLogStub = nil
	fake.runAndLogReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRuncCmdRunner) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.runAndLogMutex.RLock()
	defer fake.runAndLogMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeRuncCmdRunner) recordInvocation(key string, args []interface{}) {
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

var _ runrunc.RuncCmdRunner = new(FakeRuncCmdRunner)