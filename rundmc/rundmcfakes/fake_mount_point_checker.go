// Code generated by counterfeiter. DO NOT EDIT.
package rundmcfakes

import (
	"sync"

	"code.cloudfoundry.org/guardian/rundmc"
)

type FakeMountPointChecker struct {
	Stub        func(path string) (bool, error)
	mutex       sync.RWMutex
	argsForCall []struct {
		path string
	}
	returns struct {
		result1 bool
		result2 error
	}
	returnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeMountPointChecker) Spy(path string) (bool, error) {
	fake.mutex.Lock()
	ret, specificReturn := fake.returnsOnCall[len(fake.argsForCall)]
	fake.argsForCall = append(fake.argsForCall, struct {
		path string
	}{path})
	fake.recordInvocation("MountPointChecker", []interface{}{path})
	fake.mutex.Unlock()
	if fake.Stub != nil {
		return fake.Stub(path)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.returns.result1, fake.returns.result2
}

func (fake *FakeMountPointChecker) CallCount() int {
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	return len(fake.argsForCall)
}

func (fake *FakeMountPointChecker) ArgsForCall(i int) string {
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	return fake.argsForCall[i].path
}

func (fake *FakeMountPointChecker) Returns(result1 bool, result2 error) {
	fake.Stub = nil
	fake.returns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeMountPointChecker) ReturnsOnCall(i int, result1 bool, result2 error) {
	fake.Stub = nil
	if fake.returnsOnCall == nil {
		fake.returnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.returnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeMountPointChecker) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeMountPointChecker) recordInvocation(key string, args []interface{}) {
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

var _ rundmc.MountPointChecker = new(FakeMountPointChecker).Spy
