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
	"cana.io/clap/pkg/model"
	"errors"
	"xorm.io/xorm"
)

func updateTimetableById(session *xorm.Session, id uint64, status int) (interface{}, error) {
	var info model.Timetable
	get, err := session.Cols(model.TaskStatusInTimetable).
		ForUpdate().ID(id).Get(&info)
	if nil != err {
		return -1, err
	}
	if !get {
		return 0, errors.New("deploy not exist")
	}
	if info.TaskStatus == status {
		return 0, nil
	}
	info.TaskStatus = status
	var result int64
	result, err = session.Cols(model.TaskStatusInTimetable).
		Where(model.IdInTimetable+" = ?", id).Update(info)
	if nil != err {
		return result, err
	}
	return 1, nil
}

func getLatestTimetableResult(taskId uint64) (*model.TimetableResult, error) {
	var result []model.TimetableResult
	err := base.Engine.Cols(model.LastStatusInTimetableResult,
		model.LastResultInTimetableResult, model.CreatedAt).Where(model.TaskIdInTimetableResult+" = ?", taskId).
		Desc(model.IdInTimetableResult).Limit(1).Find(&result)
	if nil == err {
		if len(result) > 0 {
			return &result[0], nil
		} else {
			return nil, nil
		}
	} else {
		return nil, err
	}
}

func insertTimetableResult(session *xorm.Session, info *model.TimetableResult) (int64, error) {
	return session.InsertOne(info)
}

func findReadyTimetableForTask() ([]model.Timetable, error) {
	var list []model.Timetable
	err := base.Engine.Cols(model.IdInTimetable, model.TaskCronInTimetable,
		model.TaskTypeInTimetable, model.TaskInfoInTimetable).
		Where(model.TaskStatusInTimetable + " = 0").Where(model.IsDisableInTimetable + " = 0").
		Find(&list)
	return list, err
}

func findAllSpaceForTask() ([]model.EnvironmentSpace, error) {
	var list []model.EnvironmentSpace
	err := base.Engine.Cols(model.IdInEnvironmentSpace, model.EnvIdInEnvironmentSpace,
		model.SpaceKeepInEnvironmentSpace, model.SpaceInfoInEnvironmentSpace).
		Where(model.IsDisableInEnvironmentSpace + " = 0").
		Find(&list)
	return list, err
}
