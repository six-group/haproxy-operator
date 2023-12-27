package utils

import (
	"fmt"

	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	proxyv1alpha1 "github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
)

func GetConfigSecretName(instance *proxyv1alpha1.Instance) string {
	return fmt.Sprintf("%s-haproxy-config", instance.Name)
}

func GetServiceName(instance *proxyv1alpha1.Instance) string {
	return fmt.Sprintf("%s-haproxy", instance.Name)
}

func GetRouteName(frontend *configv1alpha1.Frontend, bind configv1alpha1.Bind) string {
	if bind.Name != "" {
		return fmt.Sprintf("%s-%s-haproxy", frontend.Name, bind.Name)
	}

	return fmt.Sprintf("%s-haproxy", frontend.Name)
}
