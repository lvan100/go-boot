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

package main

import (
	"time"

	"github.com/lvan100/go-boot"
)

type AppCtx struct {
	Value time.Duration
}

func main() {
	app := boot.NewApp(&AppCtx{
		Value: time.Second,
	})

	app.InitConf = func(appCtx *AppCtx) {}
	app.InitLoggers = func(appCtx *AppCtx) {}
	app.CloseLoggers = func(appCtx *AppCtx) {}
	app.InitClients = func(appCtx *AppCtx) {}
	app.CloseClients = func(appCtx *AppCtx) {}
	app.StartTasks = func(appCtx *AppCtx) {}
	app.StopTasks = func(appCtx *AppCtx) {}
	app.StartServers = func(appCtx *AppCtx) {
		go func() {
			app.Msg("program is sleeping")
			time.Sleep(appCtx.Value)
			app.ShutDown("timeout")
		}()
	}
	app.StopServers = func(appCtx *AppCtx) {}

	app.Run()
}
