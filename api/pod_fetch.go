package api

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/refer"
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func restartPodByName(envId uint64, namespace string, name string) error {
	k8s, _, err := base.K8S(envId)
	if nil != err {
		return err
	}
	return k8s.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}
