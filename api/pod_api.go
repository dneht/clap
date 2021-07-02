package api

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"cana.io/clap/pkg/xterm"
	"cana.io/clap/util"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	v1 "k8s.io/api/core/v1"
	"log"
	"strconv"
	"strings"
)

func ListSpacePod(c *fiber.Ctx) error {
	spaceId, err := util.CheckIdInput(c, "id")
	if nil != err {
		return util.ErrorInput(c, "space id must be set")
	}
	spaceBase, err := getSpaceById(spaceId)
	if nil != err {
		return util.ErrorInternal(c, err)
	}
	envBase, err := base.Env(spaceBase.EnvId)
	if nil != err {
		return util.ErrorInternal(c, err)
	}
	return listPodByType(c, spaceBase.EnvId, &[]string{refer.LabelAppEnv + "=" + envBase.Env, refer.LabelAppSpace + "=" + spaceBase.SpaceName})
}

func ListDeployPod(c *fiber.Ctx) error {
	deployId, err := util.CheckIdInput(c, "id")
	if nil != err {
		return util.ErrorInput(c, "env id must be set")
	}
	envBase, spaceBase, appBase, deployBase, err := getMoreModels(deployId)
	if nil != err {
		return util.ErrorInternal(c, err)
	}
	return listPodByType(c, deployBase.EnvId, selectDeployLabelList(envBase, spaceBase, appBase))
}

func RestartDeployPod(c *fiber.Ctx) error {
	deployId, err := util.CheckIdInput(c, "id")
	if nil != err {
		return util.ErrorInput(c, "env id must be set")
	}
	spaceId, err := util.CheckIdQuery(c, "sid")
	if nil != err {
		return util.ErrorInput(c, "env id must be set")
	}
	appId, err := util.CheckIdQuery(c, "aid")
	if nil != err {
		return util.ErrorInput(c, "env id must be set")
	}
	podName := c.Query("pod")
	_, spaceBase, appBase, deployBase, err := getMoreModels(deployId)
	if nil != err {
		return util.ErrorInternal(c, err)
	}
	if spaceBase.Id != spaceId || appBase.Id != appId || strings.Index(podName, appBase.AppName) < 0 {
		return util.ErrorInput(c, "input check error")
	}
	return restartPodByName(deployBase.EnvId, spaceBase.SpaceKeep, podName)
}

func ExecSelect(ws *websocket.Conn) {
	err := xterm.ExecSelectPod(ws, func(ws *websocket.Conn) (*xterm.ExecInput, error) {
		resourceType := ws.Params("type")
		if resourceType == refer.SelectPodInnerType {
			return execInnerSelect(ws)
		} else {
			return execOtherSelect(ws, resourceType)
		}
	})
	if nil != err {
		log.Println("[error] select exec: ", err)
	}
}

func execInnerSelect(ws *websocket.Conn) (*xterm.ExecInput, error) {
	envId, err := util.CheckWsIdQuery(ws, "env")
	if nil != err {
		return nil, errors.New("env must be set")
	}
	podName := ws.Params("pod")
	if "" == podName {
		return nil, errors.New("pod name must be set")
	}

	containerName := ws.Query("container")
	tailLines, err := util.CheckWsIntQueryInput(ws, "tail", 100)
	if nil != err {
		return nil, err
	}
	return &xterm.ExecInput{
		EnvId:         envId,
		Namespace:     base.Now().Namespace,
		PodName:       podName,
		ContainerName: containerName,
		Resource:      refer.SelectPodAttachType,
		TailLines:     tailLines,
		Timeout:       600,
	}, nil
}

func execOtherSelect(ws *websocket.Conn, resourceType string) (*xterm.ExecInput, error) {
	envStr := ws.Query("env")
	spaceStr := ws.Query("space")
	projectStr := ws.Query("project")
	deployStr := ws.Query("deploy")

	var envBase *model.Environment
	var spaceBase *model.EnvironmentSpace
	var appBase *model.Project
	var deploy bool
	var err error
	if "" == deployStr {
		deploy = false
		if "" == envStr && "" == spaceStr && "" == projectStr {
			return nil, errors.New("deploy id must be set")
		}
		var envId uint64
		envId, err = strconv.ParseUint(envStr, 10, 64)
		if nil != err {
			return nil, errors.New("env id is error")
		}
		var spaceId uint64
		spaceId, err = strconv.ParseUint(spaceStr, 10, 64)
		if nil != err {
			return nil, errors.New("space id is error")
		}
		var projectId uint64
		projectId, err = strconv.ParseUint(projectStr, 10, 64)
		if nil != err {
			return nil, errors.New("project id is error")
		}
		envBase, err = base.Env(envId)
		appBase, err = getAppById(projectId)
		spaceBase, err = getSpaceById(spaceId)
	} else {
		deploy = true
		var deployId uint64
		deployId, err = strconv.ParseUint(deployStr, 10, 64)
		if nil != err {
			return nil, errors.New("deploy id is error")
		}
		envBase, spaceBase, appBase, _, err = getMoreModels(deployId)
	}
	if nil != err {
		return nil, err
	}
	if resourceType == refer.SelectPodExecType && appBase.IsIngress == 0 {
		return nil, errors.New("shell is disabled")
	}
	podName := ws.Params("pod")
	if "" == podName {
		return nil, errors.New("pod name must be set")
	}
	namespace := spaceBase.SpaceKeep
	var podList *v1.PodList
	if deploy {
		podList, err = listPodByLabel(envBase.Id, namespace, selectDeployLabelList(envBase, spaceBase, appBase))
	} else {
		podList, err = listPodByLabel(envBase.Id, namespace, selectAppLabelList(envBase, appBase))
	}
	if nil != err {
		return nil, err
	}
	checkIn := false
	for _, pod := range podList.Items {
		if pod.Name == podName {
			checkIn = true
			break
		}
	}
	if !checkIn {
		return nil, errors.New("pod name not exist")
	}

	containerName := ws.Query("container")
	tailLines, err := util.CheckWsIntQueryInput(ws, "tail", 100)
	if nil != err {
		return nil, err
	}
	return &xterm.ExecInput{
		EnvId:         envBase.Id,
		Namespace:     namespace,
		PodName:       podName,
		ContainerName: containerName,
		Resource:      resourceType,
		TailLines:     tailLines,
		Timeout:       600,
	}, nil
}

func listPodByType(c *fiber.Ctx, envId uint64, labels *[]string) error {
	namespace := c.Query("namespace")
	//sinceTime := c.QueryParam("since")
	list, err := listPodByLabel(envId, namespace, labels)
	var pods []*refer.PodInfo
	if nil == err {
		pods = refer.BuildListFromPod(list)
		//filter
	}
	return util.ResultList(c, err, len(pods), pods)
}
