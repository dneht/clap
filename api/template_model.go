package api

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/model"
	"cana.io/clap/util"
	"errors"
	"github.com/gofiber/fiber/v2"
	"xorm.io/xorm"
)

func getTemplateById(id uint64) (*model.Template, error) {
	value, ok := templateMap[id]
	if ok {
		return value, nil
	}

	var info model.Template
	result, err := base.Engine.ID(id).Get(&info)
	if nil == err {
		templateMap[id] = &info
	}
	if !result {
		return nil, errors.New("can not found template")
	}
	return &info, err
}

func invalidTemplateById(id uint64) {
	delete(templateMap, id)
}

func findAllTemplateSimple() (int, []model.Environment, error) {
	var list []model.Environment
	err := base.Engine.Cols(model.IdInTemplate, model.TemplateNameInTemplate).Find(&list)
	return len(list), list, err
}

func countTemplateWithPage(c *fiber.Ctx, input *util.MainInput) (int64, error) {
	var info model.Template
	sql := base.Engine.Cols()
	err := SelectAuth(c, model.TemplateTable, sql)
	if nil != err {
		return 0, err
	}
	return input.Apply(sql).Count(&info)
}

func findTemplateWithPage(c *fiber.Ctx, input *util.MainInput) (int, []model.Template, error) {
	var list []model.Template
	sql := base.Engine.Omit(model.CreatedAt, model.UpdatedAt)
	err := SelectAuth(c, model.TemplateTable, sql)
	if nil != err {
		return 0, nil, err
	}
	err = input.Apply(sql).Find(&list)
	return len(list), list, err
}

func updateTemplateById(session *xorm.Session, info *model.Template) (int64, error) {
	if nil == info || info.Id <= 0 {
		return -1, errors.New("input model error, id is empty")
	}
	invalidTemplateById(info.Id)
	return session.Omit(model.IdInTemplate).Update(info)
}

func insertTemplate(session *xorm.Session, info *model.Template) (int64, error) {
	return session.InsertOne(info)
}
