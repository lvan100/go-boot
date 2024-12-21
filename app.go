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
	"os"
	"os/signal"
	"syscall"
)

// AppContext is the context of the application.
type AppContext any

// App is a web application with some fixed steps.
type App[T AppContext] struct {
	Bootstrap    BootstrapInterface
	appCtx       T
	exitChan     chan struct{}
	InitConf     func(appCtx T)
	InitLoggers  func(appCtx T)
	InitClients  func(appCtx T)
	StartTasks   func(appCtx T)
	StartServers func(appCtx T)
	StopServers  func(appCtx T)
	StopTasks    func(appCtx T)
	CloseClients func(appCtx T)
	CloseLoggers func(appCtx T)
}

// NewApp creates a new App.
func NewApp[T AppContext](appCtx T) *App[T] {
	return &App[T]{
		appCtx:   appCtx,
		exitChan: make(chan struct{}),
	}
}

// Msg prints a message to the console.
func (app *App[T]) Msg(msg string) {
	log.Println("application: " + msg)
}

// Run runs the application.
func (app *App[T]) Run() {
	defer func() {
		app.Msg("program is exited")
	}()

	// bootstrap
	if app.Bootstrap != nil {
		app.Msg("bootstrap before app run")
		app.Bootstrap.doBeforeAppRun()
		defer func() {
			app.Msg("bootstrap after app exit")
			app.Bootstrap.doAfterAppExit()
		}()
	}

	// config
	if app.InitConf != nil {
		app.Msg("init config")
		app.InitConf(app.appCtx)
	}

	// logger
	{
		defer func() {
			if app.CloseLoggers != nil {
				app.Msg("close loggers")
				app.CloseLoggers(app.appCtx)
			}
		}()
		if app.InitLoggers != nil {
			app.Msg("init loggers")
			app.InitLoggers(app.appCtx)
		}
	}

	// client
	{
		defer func() {
			if app.CloseClients != nil {
				app.Msg("close clients")
				app.CloseClients(app.appCtx)
			}
		}()
		if app.InitClients != nil {
			app.Msg("init clients")
			app.InitClients(app.appCtx)
		}
	}

	// task
	{
		defer func() {
			if app.StopTasks != nil {
				app.Msg("stop tasks")
				app.StopTasks(app.appCtx)
			}
		}()
		if app.StartTasks != nil {
			app.Msg("start tasks")
			app.StartTasks(app.appCtx)
		}
	}

	// server
	{
		defer func() {
			if app.StopServers != nil {
				app.Msg("stop servers")
				app.StopServers(app.appCtx)
			}
		}()
		if app.StartServers != nil {
			app.Msg("start servers")
			app.StartServers(app.appCtx)
		}
	}

	app.Msg("program is running")

	// signal
	{
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
			sig := <-ch
			app.ShutDown("received signal " + sig.String())
		}()
		<-app.exitChan
	}
}

// ShutDown shuts down the application.
func (app *App[T]) ShutDown(msg string) {
	app.Msg("program is exiting, " + msg)
	select {
	case <-app.exitChan:
	default:
		close(app.exitChan)
	}
}
