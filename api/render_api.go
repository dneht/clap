package api

import (
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"cana.io/clap/util"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-jsonnet"
	"strconv"
)

const RenderApiPre = "/api/render"

func GetTemplate(c *fiber.Ctx) error {
	id, err := util.CheckIdInput(c, "id")
	if nil != err {
		return err
	}
	info, err := getTemplateById(id)
	return util.ResultParam(c, err, true, info)
}

func SimpleTemplate(c *fiber.Ctx) error {
	size, list, err := findAllTemplateSimple()
	if nil != err {
		return util.ErrorInternal(c, err)
	}
	return util.ResultList(c, err, size, list)
}

func ListTemplate(c *fiber.Ctx) error {
	param, err := util.CheckMainInput(c)
	if nil != err {
		return err
	}
	return util.ResultPageOrList(c, param,
		func(input *util.MainInput) (int64, error) {
			return countTemplateWithPage(c, input)
		}, func(input *util.MainInput) (int, interface{}, error) {
			return findTemplateWithPage(c, input)
		})
}

func ExecRender(c *fiber.Ctx) error {
	deployId, err := util.CheckIdInput(c, "deploy")
	if nil != err {
		return err
	}
	//err = DeploymentAuth(c, deployId, AllowThisPackageDeploy)
	//if nil != err {
	//	return err
	//}
	envBase, spaceBase, appBase, deployBase, err := getMoreModels(deployId)
	if nil != err {
		return util.ErrorInternal(c, err)
	}
	selectType := c.Params("type")
	if "jsonnet" == selectType {
		imageUrl := c.Query("image")
		if "" == imageUrl {
			return errors.New("must set image url")
		}

		_, jsonStr, err := renderJsonnet(envBase, spaceBase, appBase, deployBase, imageUrl)
		if nil != err {
			return err
		}
		return c.SendString(jsonStr)
	}

	return errors.New("select type is not support")
}

func renderJsonnet(envBase *model.Environment, spaceBase *model.EnvironmentSpace, appBase *model.Project, deployBase *model.Deployment, repoImage string) (string, string, error) {
	template, appJson, err := buildTemplate(envBase, spaceBase, appBase, deployBase, repoImage)
	if nil != err {
		return "", "", err
	}
	if "" == template.TemplateContent {
		return "", "", errors.New("template is empty: " + strconv.FormatUint(template.Id, 10))
	}

	vm := jsonnet.MakeVM()
	vm.TLACode("app", appJson)
	jsonStr, err := vm.EvaluateSnippet("", template.TemplateContent)
	if nil != err {
		return "", "", err
	}
	return template.TemplateKind, jsonStr, nil
}

func buildTemplate(envBase *model.Environment, spaceBase *model.EnvironmentSpace, appBase *model.Project, deployBase *model.Deployment, repoImage string) (*model.Template, string, error) {
	_, appInfo, err := getMoreInfos(envBase.Id, appBase.Id, deployBase.Id)
	if nil != err {
		return nil, "", err
	}
	templateBase, err := getTemplateById(appInfo.Template)
	if nil != err {
		return nil, "", err
	}
	err = generateRenderProps(appBase.Id, envBase.Id, spaceBase.Id, deployBase.Id, appInfo)
	if nil != err {
		return nil, "", err
	}

	fillMap := make(map[string]interface{}, 40)
	fillMap["id"] = appBase.Id
	fillMap["key"] = appBase.AppKey
	fillMap["name"] = refer.GetAppName(appBase, spaceBase)
	fillMap["type"] = refer.ConvertAppType(appBase.AppType)
	fillMap["kind"] = templateBase.TemplateKind
	fillMap["env"] = envBase.Env
	fillMap["image"] = repoImage
	fillMap["space"] = spaceBase.SpaceName
	fillMap["namespace"] = spaceBase.SpaceKeep
	fillMap["component"] = appInfo.Component
	formatStr := envBase.FormatInfo
	var formatInfo refer.FormatInfo
	if err = json.Unmarshal([]byte(formatStr), &formatInfo); err != nil {
		return nil, "", err
	}
	fillMap["specs"] = formatInfo.Spec
	for key, appValue := range appInfo.Param {
		fillMap[key] = appValue
	}
	fillMap["labels"] = configAppLabelMap(envBase, spaceBase, appBase, appInfo.Component)
	fillMap["selector"] = selectAppLabelMap(envBase, spaceBase, appBase, appInfo.Component)
	fillMap["constant"] = map[string]string{
		"type":      refer.LabelAppType,
		"name":      refer.LabelAppName,
		"env":       refer.LabelAppEnv,
		"space":     refer.LabelAppSpace,
		"managed":   refer.LabelAppManaged,
		"component": refer.LabelAppComponent,
	}
	fillMap["contour"] = map[string]string{
		"apiVersion": refer.ContourApiVersion,
	}
	bytes, err := json.Marshal(fillMap)
	if nil != err {
		return nil, "", err
	}
	return templateBase, string(bytes), nil
}
