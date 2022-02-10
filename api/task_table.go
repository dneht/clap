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

package api

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/log"
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"encoding/json"
	"github.com/robfig/cron/v3"
	"time"
	"xorm.io/xorm"
)

func InitTask() {
	exec := cron.New()
	tasks, err := findReadyTimetableForTask()
	if nil != err {
		log.Warnf("find ready task failed: %v", err)
		return
	}

	for _, task := range *tasks {
		if "" == task.TaskCron {
			continue
		}
		_, err = exec.AddFunc(task.TaskCron, func() {
			err = startTask(&task)
			if nil != err {
				log.Warnf("start task failed: %v, skip execute", err)
			} else {
				var failed error
				var result interface{}
				switch task.TaskType {
				case refer.AcmeAliyunTaskType, refer.AcmeDnspodTaskType:
					result, failed = execTaskAcmeDomains(&task)
					break
				default:
					log.Warnf("task type not support: %v | %v", task.Id, task.TaskType)
				}
				finishTask(&task, failed, result)
			}
		})
		if nil != err {
			log.Warnf("parse cron failed: %v, skip execute", err)
		} else {
			log.Infof("ready add acme task: %v | %v", task.Id, task.TaskType)
		}
	}
	exec.Start()
}

func startTask(task *model.Timetable) error {
	_, err := base.Engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		return updateTimetableById(session, task.Id, refer.TaskStatusExecuting)
	})
	return err
}

func finishTask(task *model.Timetable, failed error, result interface{}) {
	if nil != result {
		status := 0
		if nil != failed {
			status = 1
		}
		taskResult := "{}"
		resultBytes, err := json.Marshal(result)
		if nil == err {
			taskResult = string(resultBytes)
		}
		_, err = base.Engine.Transaction(func(session *xorm.Session) (interface{}, error) {
			_, _ = insertTimetableResult(session, &model.TimetableResult{
				TaskId:     task.Id,
				LastStatus: status,
				LastResult: taskResult,
				CreatedAt:  time.Now(),
			})
			return updateTimetableById(session, task.Id, refer.TaskStatusWaiting)
		})
		if nil != err {
			log.Errorf("finish task with result failed: %v", err)
		}
	} else {
		_, err := base.Engine.Transaction(func(session *xorm.Session) (interface{}, error) {
			return updateTimetableById(session, task.Id, refer.TaskStatusWaiting)
		})
		if nil != err {
			log.Errorf("finish task without result failed: %v", err)
		}
	}
}
