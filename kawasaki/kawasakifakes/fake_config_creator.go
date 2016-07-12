// This file was generated by counterfeiter
package kawasakifakes

import (
	"net"
	"sync"

	"code.cloudfoundry.org/lager"
	"github.com/cloudfoundry-incubator/guardian/kawasaki"
)

type FakeConfigCreator struct {
	CreateStub        func(log lager.Logger, handle string, subnet *net.IPNet, ip net.IP) (kawasaki.NetworkConfig, error)
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		log    lager.Logger
		handle string
		subnet *net.IPNet
		ip     net.IP
	}
	createReturns struct {
		result1 kawasaki.NetworkConfig
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeConfigCreator) Create(log lager.Logger, handle string, subnet *net.IPNet, ip net.IP) (kawasaki.NetworkConfig, error) {
	fake.createMutex.Lock()
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		log    lager.Logger
		handle string
		subnet *net.IPNet
		ip     net.IP
	}{log, handle, subnet, ip})
	fake.recordInvocation("Create", []interface{}{log, handle, subnet, ip})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(log, handle, subnet, ip)
	} else {
		return fake.createReturns.result1, fake.createReturns.result2
	}
}

func (fake *FakeConfigCreator) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeConfigCreator) CreateArgsForCall(i int) (lager.Logger, string, *net.IPNet, net.IP) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.createArgsForCall[i].log, fake.createArgsForCall[i].handle, fake.createArgsForCall[i].subnet, fake.createArgsForCall[i].ip
}

func (fake *FakeConfigCreator) CreateReturns(result1 kawasaki.NetworkConfig, result2 error) {
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 kawasaki.NetworkConfig
		result2 error
	}{result1, result2}
}

func (fake *FakeConfigCreator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeConfigCreator) recordInvocation(key string, args []interface{}) {
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

var _ kawasaki.ConfigCreator = new(FakeConfigCreator)