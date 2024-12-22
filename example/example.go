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
	"fmt"
	"time"

	"github.com/lvan100/go-boot"
)

type BootCtx struct {
	Value int
}

type AppCtx struct {
	Value time.Duration
}

func main() {

	bootstrap := boot.NewBootstrap(&BootCtx{
		Value: 1,
	})

	bootstrap.InitConf = func(bootCtx *BootCtx) {
		p, err := bootstrap.Bootstrapper.Refresh()
		if err != nil {
			panic(err)
		}
		fmt.Println(p.Data())
	}
	bootstrap.InitLoggers = func(bootCtx *BootCtx) {}
	bootstrap.CloseLoggers = func(bootCtx *BootCtx) {}
	bootstrap.InitClients = func(bootCtx *BootCtx) {}
	bootstrap.CloseClients = func(bootCtx *BootCtx) {}

	bootstrap.Bootstrap = func(bootCtx *BootCtx) {
		bootstrap.Msg("bootstrap is run")
	}

	app := boot.NewApp(&AppCtx{
		Value: time.Second,
	})

	app.InitConf = func(appCtx *AppCtx) {
		p, err := app.Configuration.Refresh()
		if err != nil {
			panic(err)
		}
		fmt.Println(p.Data())
	}
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

	app.Bootstrap = bootstrap
	app.Run()
}
