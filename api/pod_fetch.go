package api

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/refer"
	"context"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"os"
	"strings"
)

func listAllPod(envId uint64, namespace string) (*v1.PodList, error) {
	k8s, _, err := base.K8S(envId)
	if nil != err {
		return nil, err
	}
	return k8s.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
}

func listPodByLabel(envId uint64, namespace string, labels *[]string) (*v1.PodList, error) {
	k8s, _, err := base.K8S(envId)
	if nil != err {
		return nil, err
	}
	return k8s.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: refer.LabelAppManaged + "=Clap, " +
			strings.Join(*labels, ", "),
	})
}

func restartPodByName(envId uint64, namespace string, pod string) error {
	k8s, _, err := base.K8S(envId)
	if nil != err {
		return err
	}
	return k8s.CoreV1().Pods(namespace).Delete(context.TODO(), pod, metav1.DeleteOptions{})
}

func downloadPodByName(envId uint64, namespace, pod, container, filename string) (*io.PipeReader, error) {
	k8s, config, err := base.K8S(envId)
	if nil != err {
		return nil, err
	}
	reader, outStream := io.Pipe()
	req := k8s.CoreV1().RESTClient().Get().
		Resource("pods").
		Name(pod).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: container,
			Command:   []string{"tar", "cf", "-", filename},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return nil, err
	}
	go func() {
		defer outStream.Close()
		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  os.Stdin,
			Stdout: outStream,
			Stderr: os.Stderr,
			Tty:    false,
		})
	}()
	return reader, err
}
