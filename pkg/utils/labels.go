package utils

import (
	"github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func GetAppSelectorLabels(instance *v1alpha1.Instance) map[string]string {
	return map[string]string{
		corev1.LabelMetadataName: GetServiceAndStatefulsetName(instance),
	}
}

func GetPodLabels(instance *v1alpha1.Instance) map[string]string {
	r := GetAppSelectorLabels(instance)
	for k, v := range instance.Spec.Labels {
		r[k] = v
	}
	return r
}
