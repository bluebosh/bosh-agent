// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/cloudfoundry/bosh-agent/agent/script"
	boshdrain "github.com/cloudfoundry/bosh-agent/agent/script/drain"
)

type FakeJobScriptProvider struct {
	NewScriptStub        func(jobName string, scriptName string) script.Script
	newScriptMutex       sync.RWMutex
	newScriptArgsForCall []struct {
		jobName    string
		scriptName string
	}
	newScriptReturns struct {
		result1 script.Script
	}
	NewDrainScriptStub        func(jobName string, params boshdrain.ScriptParams) script.Script
	newDrainScriptMutex       sync.RWMutex
	newDrainScriptArgsForCall []struct {
		jobName string
		params  boshdrain.ScriptParams
	}
	newDrainScriptReturns struct {
		result1 script.Script
	}
	NewParallelScriptStub        func(scriptName string, scripts []script.Script) script.Script
	newParallelScriptMutex       sync.RWMutex
	newParallelScriptArgsForCall []struct {
		scriptName string
		scripts    []script.Script
	}
	newParallelScriptReturns struct {
		result1 script.Script
	}
}

func (fake *FakeJobScriptProvider) NewScript(jobName string, scriptName string) script.Script {
	fake.newScriptMutex.Lock()
	fake.newScriptArgsForCall = append(fake.newScriptArgsForCall, struct {
		jobName    string
		scriptName string
	}{jobName, scriptName})
	fake.newScriptMutex.Unlock()
	if fake.NewScriptStub != nil {
		return fake.NewScriptStub(jobName, scriptName)
	} else {
		return fake.newScriptReturns.result1
	}
}

func (fake *FakeJobScriptProvider) NewScriptCallCount() int {
	fake.newScriptMutex.RLock()
	defer fake.newScriptMutex.RUnlock()
	return len(fake.newScriptArgsForCall)
}

func (fake *FakeJobScriptProvider) NewScriptArgsForCall(i int) (string, string) {
	fake.newScriptMutex.RLock()
	defer fake.newScriptMutex.RUnlock()
	return fake.newScriptArgsForCall[i].jobName, fake.newScriptArgsForCall[i].scriptName
}

func (fake *FakeJobScriptProvider) NewScriptReturns(result1 script.Script) {
	fake.NewScriptStub = nil
	fake.newScriptReturns = struct {
		result1 script.Script
	}{result1}
}

func (fake *FakeJobScriptProvider) NewDrainScript(jobName string, params boshdrain.ScriptParams) script.Script {
	fake.newDrainScriptMutex.Lock()
	fake.newDrainScriptArgsForCall = append(fake.newDrainScriptArgsForCall, struct {
		jobName string
		params  boshdrain.ScriptParams
	}{jobName, params})
	fake.newDrainScriptMutex.Unlock()
	if fake.NewDrainScriptStub != nil {
		return fake.NewDrainScriptStub(jobName, params)
	} else {
		return fake.newDrainScriptReturns.result1
	}
}

func (fake *FakeJobScriptProvider) NewDrainScriptCallCount() int {
	fake.newDrainScriptMutex.RLock()
	defer fake.newDrainScriptMutex.RUnlock()
	return len(fake.newDrainScriptArgsForCall)
}

func (fake *FakeJobScriptProvider) NewDrainScriptArgsForCall(i int) (string, boshdrain.ScriptParams) {
	fake.newDrainScriptMutex.RLock()
	defer fake.newDrainScriptMutex.RUnlock()
	return fake.newDrainScriptArgsForCall[i].jobName, fake.newDrainScriptArgsForCall[i].params
}

func (fake *FakeJobScriptProvider) NewDrainScriptReturns(result1 script.Script) {
	fake.NewDrainScriptStub = nil
	fake.newDrainScriptReturns = struct {
		result1 script.Script
	}{result1}
}

func (fake *FakeJobScriptProvider) NewParallelScript(scriptName string, scripts []script.Script) script.Script {
	fake.newParallelScriptMutex.Lock()
	fake.newParallelScriptArgsForCall = append(fake.newParallelScriptArgsForCall, struct {
		scriptName string
		scripts    []script.Script
	}{scriptName, scripts})
	fake.newParallelScriptMutex.Unlock()
	if fake.NewParallelScriptStub != nil {
		return fake.NewParallelScriptStub(scriptName, scripts)
	} else {
		return fake.newParallelScriptReturns.result1
	}
}

func (fake *FakeJobScriptProvider) NewParallelScriptCallCount() int {
	fake.newParallelScriptMutex.RLock()
	defer fake.newParallelScriptMutex.RUnlock()
	return len(fake.newParallelScriptArgsForCall)
}

func (fake *FakeJobScriptProvider) NewParallelScriptArgsForCall(i int) (string, []script.Script) {
	fake.newParallelScriptMutex.RLock()
	defer fake.newParallelScriptMutex.RUnlock()
	return fake.newParallelScriptArgsForCall[i].scriptName, fake.newParallelScriptArgsForCall[i].scripts
}

func (fake *FakeJobScriptProvider) NewParallelScriptReturns(result1 script.Script) {
	fake.NewParallelScriptStub = nil
	fake.newParallelScriptReturns = struct {
		result1 script.Script
	}{result1}
}

var _ script.JobScriptProvider = new(FakeJobScriptProvider)
