package api

import (
	"cana.io/clap/util"
	"github.com/gofiber/fiber/v2"
)

const EnvApiPre = "/api/env"
const SpaceApiPre = "/api/space"

func GetEnv(c *fiber.Ctx) error {
	id, err := util.CheckIdInput(c, "id")
	if nil != err {
		return util.ErrorInput(c, "env id must be set")
	}
	info, err := getEnvById(id)
	//TODO add auth
	if nil == err {
		info.DeployInfo = ""
	}
	return util.ResultParam(c, err, true, info)
}

func SimpleEnv(c *fiber.Ctx) error {
	size, list, err := findAllEnvSimple(c)
	if nil != err {
		return util.ErrorInternal(c, err)
	}
	return util.ResultList(c, err, size, list)
}

func ListEnv(c *fiber.Ctx) error {
	param, err := util.CheckMainInput(c)
	if nil != err {
		return util.ErrorInputErrorMessage(c, err, "main input error")
	}
	return util.ResultPageOrList(c, param,
		func(input *util.MainInput) (int64, error) {
			return countEnvWithPage(c, input)
		}, func(input *util.MainInput) (int, interface{}, error) {
			return findEnvWithPage(c, input)
		})
}

func GetSpace(c *fiber.Ctx) error {
	envId, err := util.CheckIdInput(c, "eid")
	if nil != err {
		return util.ErrorInput(c, "env id must be set")
	}
	id, err := util.CheckIdInput(c, "id")
	if nil != err {
		return util.ErrorInput(c, "space id must be set")
	}
	info, err := getSpaceById(id)
	if nil == err {
		if info.EnvId != envId {
			return util.ErrorInput(c, "env id not match")
		}
	}
	return util.ResultParam(c, err, true, info)
}

func SimpleSpace(c *fiber.Ctx) error {
	envId, err := util.CheckIdQuery(c, "eid")
	if nil != err {
		return util.ErrorInput(c, "env id must be set")
	}
	size, list, err := findAllSpaceSimple(c, envId)
	if nil != err {
		return util.ErrorInternal(c, err)
	}
	return util.ResultList(c, err, size, list)
}

func ListSpace(c *fiber.Ctx) error {
	param, err := util.CheckReferInput(c)
	if nil != err {
		return err
	}
	return util.ResultPageOrList(c, param,
		func(input *util.MainInput) (int64, error) {
			return countSpaceWithPage(c, input)
		}, func(input *util.MainInput) (int, interface{}, error) {
			return findSpaceWithPage(c, input)
		})
}
