package utils

import (
	"fmt"

	"github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
)

func GetAppSelectorLabels(instance *v1alpha1.Instance) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name": fmt.Sprintf("%s-haproxy", instance.Name),
	}
}

func GetPodLabels(instance *v1alpha1.Instance) map[string]string {
	r := GetAppSelectorLabels(instance)
	for k, v := range instance.Spec.Labels {
		r[k] = v
	}
	return r
}
