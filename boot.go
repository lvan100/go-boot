// Copyright 2024 github.com/lvan100
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package boot

import (
	"log"
)

// BootstrapInterface is a bootstrap interface.
type BootstrapInterface interface {
	doBeforeAppRun()
	doAfterAppExit()
}

// BootstrapContext is a bootstrap context.
type BootstrapContext any

// Bootstrap is a bootstrap for app with some fixed steps.
type Bootstrap[T BootstrapContext] struct {
	bootCtx      T
	Bootstrap    func(bootCtx T)
	InitConf     func(bootCtx T)
	InitLoggers  func(bootCtx T)
	InitClients  func(bootCtx T)
	CloseClients func(bootCtx T)
	CloseLoggers func(bootCtx T)
}

var _ BootstrapInterface = (*Bootstrap[any])(nil)

// NewBootstrap creates a new bootstrap.
func NewBootstrap[T BootstrapContext](bootCtx T) *Bootstrap[T] {
	return &Bootstrap[T]{
		bootCtx: bootCtx,
	}
}

// Msg prints a message to the console.
func (bootstrap *Bootstrap[T]) Msg(msg string) {
	log.Println("bootstrap: " + msg)
}

// doBeforeAppRun bootstraps before app run.
func (bootstrap *Bootstrap[T]) doBeforeAppRun() {

	// config
	if bootstrap.InitConf != nil {
		bootstrap.Msg("init config")
		bootstrap.InitConf(bootstrap.bootCtx)
	}

	// logger
	if bootstrap.InitLoggers != nil {
		bootstrap.Msg("init loggers")
		bootstrap.InitLoggers(bootstrap.bootCtx)
	}

	// client
	if bootstrap.InitClients != nil {
		bootstrap.Msg("init clients")
		bootstrap.InitClients(bootstrap.bootCtx)
	}

	bootstrap.Bootstrap(bootstrap.bootCtx)
}

// doAfterAppExit bootstraps after app exit.
func (bootstrap *Bootstrap[T]) doAfterAppExit() {

	// client
	if bootstrap.CloseClients != nil {
		bootstrap.Msg("close clients")
		bootstrap.CloseClients(bootstrap.bootCtx)
	}

	// logger
	if bootstrap.CloseLoggers != nil {
		bootstrap.Msg("close loggers")
		bootstrap.CloseLoggers(bootstrap.bootCtx)
	}
}
