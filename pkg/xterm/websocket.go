package xterm

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/util"
	"context"
	"errors"
	"github.com/gofiber/websocket/v2"
	"log"
	"strconv"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

const pingPeriod = 5 * time.Second

var ingress = []string{"/bin/sh", "-c", "TERM=xterm-256color; export TERM; [ -x /bin/bash ] && ([ -x /usr/bin/script ] && /usr/bin/script -q -c \"/bin/bash\" /dev/null || exec /bin/bash) || exec /bin/sh"}

type ExecInput struct {
	EnvId         uint64   `json:"env"`
	Namespace     string   `json:"namespace"`
	PodName       string   `json:"pod"`
	ContainerName string   `json:"container"`
	SinceSecond   int64    `json:"since"`
	TailLines     int64    `json:"tail"`
	Resource      string   `json:"resource"`
	Timeout       int      `json:"timeout"`
	Ingress       []string `json:"ingress"`
}

func ExecSelectPod(ws *websocket.Conn, f func(ws *websocket.Conn) (*ExecInput, error)) error {
	in, err := f(ws)
	if err != nil {
		return err
	}

	t := &TerminalSession{
		id:       strconv.FormatUint(util.UniqueId(), 10),
		sizeChan: make(chan remotecommand.TerminalSize),
		ws:       ws,
		lock:     &sync.Mutex{},
		input:    in,
	}

	defer t.Toast("warning", "connection close")
	if in.Resource != "exec" && in.Resource != "attach" {
		t.Toast("error", "resource error")
		return errors.New("resource error")
	}

	clientset, config, err := base.K8S(in.EnvId)
	if err != nil {
		t.Toast("error", err.Error())
		return err
	}

	// 如果没有指定容器名的，使用第一个容器
	if len(in.PodName) != 0 && len(in.ContainerName) == 0 {
		pod, err := clientset.CoreV1().Pods(in.Namespace).Get(context.TODO(), in.PodName, metav1.GetOptions{})
		if err == nil {
			in.ContainerName = pod.Spec.Containers[0].Name
		}
	}
	t.Toast("info", "connected")

	if t.input.Resource == "attach" {
		err = Logs(t, clientset)
		if err != nil {
			t.Toast("error", err.Error())
			return err
		}
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(t.input.PodName).
		Namespace(t.input.Namespace).
		SubResource(t.input.Resource)

	tty := false
	cmd := ingress
	if nil != in.Ingress && len(in.Ingress) > 0 {
		cmd = in.Ingress
	}
	if t.input.Resource == "exec" {
		tty = true
		req.VersionedParams(&corev1.PodExecOptions{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
			Command:   cmd,
			Container: t.input.ContainerName,
		}, scheme.ParameterCodec)
	} else {
		req.VersionedParams(&corev1.PodAttachOptions{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			Container: t.input.ContainerName,
		}, scheme.ParameterCodec)
	}

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		t.Toast("error", err.Error())
		return err
	}

	done := make(chan struct{})
	go ping(t, done)
	defer close(done)

	ptyHandler := PtyHandler(t)
	if t.input.Resource == "exec" {
		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:             ptyHandler,
			Stdout:            ptyHandler,
			Stderr:            ptyHandler,
			TerminalSizeQueue: ptyHandler,
			Tty:               tty,
		})
	} else {
		err = exec.Stream(remotecommand.StreamOptions{
			// Stdin:  ptyHandler,
			Stdout:            ptyHandler,
			Stderr:            ptyHandler,
			TerminalSizeQueue: ptyHandler,
			Tty:               tty,
		})
	}

	if err != nil {
		t.Toast("error", err.Error())
		return err
	}
	return nil
}

func ping(t *TerminalSession, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// 保持与客户端连接不断开
			if err := t.ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println(err)
				return
			}
			// 保持与 kubernetes 连接不断开
			if t.input.Resource == "exec" {
				t.sizeChan <- remotecommand.TerminalSize{}
			}
		case <-done:
			return
		}
	}
}
