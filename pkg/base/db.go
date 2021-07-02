/*
Copyright 2020 Dasheng.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package base

import (
	"cana.io/clap/pkg/model"
	"fmt"
	"xorm.io/xorm/log"

	"xorm.io/core"
	"xorm.io/xorm"

	_ "github.com/go-sql-driver/mysql"
)

var Engine *xorm.Engine

func DbInit() {
	config := DatabaseConf()
	size := len(config.Addresses)

	if size == 0 {
		panic("database list is empty")
	} else if size == 1 {
		engine, err := xorm.NewEngine("mysql", fmt.Sprintf(
			"%v:%v@tcp(%v)/%v?tls=false&charset=utf8mb4&parseTime=False",
			config.User, config.Password, config.Addresses[0], config.Database))
		if nil != err {
			panic(err)
		}
		Engine = engine
	} else {
		urls := make([]string, size)
		for i := 0; i < size; i++ {
			urls[i] = fmt.Sprintf(
				"%v:%v@tcp(%v)/%v?tls=false&charset=utf8mb4&parseTime=False",
				config.User, config.Password, config.Addresses[i], config.Database)
		}
		engine, err := xorm.NewEngineGroup("mysql", urls)
		if nil != err {
			panic(err)
		}
		Engine = engine.Engine
	}

	Engine.ShowSQL(true)
	Engine.Logger().SetLevel(log.LOG_DEBUG)
	Engine.SetMapper(core.SnakeMapper{})
	Engine.SetMaxIdleConns(config.MaxIdleCon)
	Engine.SetMaxOpenConns(config.MaxOpenCon)
	BuildProperty()
}

func dangListFullProp() *[]model.Property {
	var list []model.Property
	err := Engine.Where(model.EnvInProperty+"=?", *envFlag).
		Where(model.IsDisableInProperty+"=?", 0).Find(&list)
	if nil != err {
		panic(err)
	}
	return &list
}
