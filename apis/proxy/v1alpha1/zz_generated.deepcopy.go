//go:build !ignore_autogenerated

/*
Copyright 2023 SIX Group Services Ltd., Switzerland

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	routev1 "github.com/openshift/api/route/v1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	timex "time"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Configuration) DeepCopyInto(out *Configuration) {
	*out = *in
	in.Global.DeepCopyInto(&out.Global)
	in.Defaults.DeepCopyInto(&out.Defaults)
	in.LabelSelector.DeepCopyInto(&out.LabelSelector)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Configuration.
func (in *Configuration) DeepCopy() *Configuration {
	if in == nil {
		return nil
	}
	out := new(Configuration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DefaultsConfiguration) DeepCopyInto(out *DefaultsConfiguration) {
	*out = *in
	if in.ErrorFiles != nil {
		in, out := &in.ErrorFiles, &out.ErrorFiles
		*out = make([]*configv1alpha1.ErrorFile, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(configv1alpha1.ErrorFile)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.Timeouts != nil {
		in, out := &in.Timeouts, &out.Timeouts
		*out = make(map[string]metav1.Duration, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Logging != nil {
		in, out := &in.Logging, &out.Logging
		*out = new(DefaultsLoggingConfiguration)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DefaultsConfiguration.
func (in *DefaultsConfiguration) DeepCopy() *DefaultsConfiguration {
	if in == nil {
		return nil
	}
	out := new(DefaultsConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DefaultsLoggingConfiguration) DeepCopyInto(out *DefaultsLoggingConfiguration) {
	*out = *in
	if in.HTTPLog != nil {
		in, out := &in.HTTPLog, &out.HTTPLog
		*out = new(bool)
		**out = **in
	}
	if in.TCPLog != nil {
		in, out := &in.TCPLog, &out.TCPLog
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DefaultsLoggingConfiguration.
func (in *DefaultsLoggingConfiguration) DeepCopy() *DefaultsLoggingConfiguration {
	if in == nil {
		return nil
	}
	out := new(DefaultsLoggingConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalConfiguration) DeepCopyInto(out *GlobalConfiguration) {
	*out = *in
	if in.StatsTimeout != nil {
		in, out := &in.StatsTimeout, &out.StatsTimeout
		*out = new(metav1.Duration)
		**out = **in
	}
	if in.Logging != nil {
		in, out := &in.Logging, &out.Logging
		*out = new(GlobalLoggingConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.AdditionalCertificates != nil {
		in, out := &in.AdditionalCertificates, &out.AdditionalCertificates
		*out = make([]configv1alpha1.SSLCertificate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Maxconn != nil {
		in, out := &in.Maxconn, &out.Maxconn
		*out = new(int64)
		**out = **in
	}
	if in.Nbthread != nil {
		in, out := &in.Nbthread, &out.Nbthread
		*out = new(int64)
		**out = **in
	}
	if in.TuneOptions != nil {
		in, out := &in.TuneOptions, &out.TuneOptions
		*out = new(GlobalTuneOptions)
		(*in).DeepCopyInto(*out)
	}
	if in.SSL != nil {
		in, out := &in.SSL, &out.SSL
		*out = new(GlobalSSL)
		(*in).DeepCopyInto(*out)
	}
	if in.HardStopAfter != nil {
		in, out := &in.HardStopAfter, &out.HardStopAfter
		*out = new(timex.Duration)
		**out = **in
	}
	if in.Ocsp != nil {
		in, out := &in.Ocsp, &out.Ocsp
		*out = new(GlobalOCSPConfiguration)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalConfiguration.
func (in *GlobalConfiguration) DeepCopy() *GlobalConfiguration {
	if in == nil {
		return nil
	}
	out := new(GlobalConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalLoggingConfiguration) DeepCopyInto(out *GlobalLoggingConfiguration) {
	*out = *in
	if in.SendHostname != nil {
		in, out := &in.SendHostname, &out.SendHostname
		*out = new(bool)
		**out = **in
	}
	if in.Hostname != nil {
		in, out := &in.Hostname, &out.Hostname
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalLoggingConfiguration.
func (in *GlobalLoggingConfiguration) DeepCopy() *GlobalLoggingConfiguration {
	if in == nil {
		return nil
	}
	out := new(GlobalLoggingConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalOCSPConfiguration) DeepCopyInto(out *GlobalOCSPConfiguration) {
	*out = *in
	if in.MaxDelay != nil {
		in, out := &in.MaxDelay, &out.MaxDelay
		*out = new(int64)
		**out = **in
	}
	if in.MinDelay != nil {
		in, out := &in.MinDelay, &out.MinDelay
		*out = new(int64)
		**out = **in
	}
	if in.HTTPProxy != nil {
		in, out := &in.HTTPProxy, &out.HTTPProxy
		*out = new(OcspUpdateOptionsHttpproxy)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalOCSPConfiguration.
func (in *GlobalOCSPConfiguration) DeepCopy() *GlobalOCSPConfiguration {
	if in == nil {
		return nil
	}
	out := new(GlobalOCSPConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalSSL) DeepCopyInto(out *GlobalSSL) {
	*out = *in
	if in.DefaultBindCiphers != nil {
		in, out := &in.DefaultBindCiphers, &out.DefaultBindCiphers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.DefaultBindCipherSuites != nil {
		in, out := &in.DefaultBindCipherSuites, &out.DefaultBindCipherSuites
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.DefaultBindOptions != nil {
		in, out := &in.DefaultBindOptions, &out.DefaultBindOptions
		*out = new(GlobalSSLDefaultBindOptions)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalSSL.
func (in *GlobalSSL) DeepCopy() *GlobalSSL {
	if in == nil {
		return nil
	}
	out := new(GlobalSSL)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalSSLDefaultBindOptions) DeepCopyInto(out *GlobalSSLDefaultBindOptions) {
	*out = *in
	if in.MinVersion != nil {
		in, out := &in.MinVersion, &out.MinVersion
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalSSLDefaultBindOptions.
func (in *GlobalSSLDefaultBindOptions) DeepCopy() *GlobalSSLDefaultBindOptions {
	if in == nil {
		return nil
	}
	out := new(GlobalSSLDefaultBindOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalSSLTuneOptions) DeepCopyInto(out *GlobalSSLTuneOptions) {
	*out = *in
	if in.CacheSize != nil {
		in, out := &in.CacheSize, &out.CacheSize
		*out = new(int64)
		**out = **in
	}
	if in.Lifetime != nil {
		in, out := &in.Lifetime, &out.Lifetime
		*out = new(metav1.Duration)
		**out = **in
	}
	if in.MaxRecord != nil {
		in, out := &in.MaxRecord, &out.MaxRecord
		*out = new(int64)
		**out = **in
	}
	if in.CaptureBufferSize != nil {
		in, out := &in.CaptureBufferSize, &out.CaptureBufferSize
		*out = new(int64)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalSSLTuneOptions.
func (in *GlobalSSLTuneOptions) DeepCopy() *GlobalSSLTuneOptions {
	if in == nil {
		return nil
	}
	out := new(GlobalSSLTuneOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalTuneOptions) DeepCopyInto(out *GlobalTuneOptions) {
	*out = *in
	if in.Maxrewrite != nil {
		in, out := &in.Maxrewrite, &out.Maxrewrite
		*out = new(int64)
		**out = **in
	}
	if in.BuffersLimit != nil {
		in, out := &in.BuffersLimit, &out.BuffersLimit
		*out = new(int64)
		**out = **in
	}
	if in.Bufsize != nil {
		in, out := &in.Bufsize, &out.Bufsize
		*out = new(int64)
		**out = **in
	}
	if in.SSL != nil {
		in, out := &in.SSL, &out.SSL
		*out = new(GlobalSSLTuneOptions)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalTuneOptions.
func (in *GlobalTuneOptions) DeepCopy() *GlobalTuneOptions {
	if in == nil {
		return nil
	}
	out := new(GlobalTuneOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Instance) DeepCopyInto(out *Instance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Instance.
func (in *Instance) DeepCopy() *Instance {
	if in == nil {
		return nil
	}
	out := new(Instance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Instance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceList) DeepCopyInto(out *InstanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Instance, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceList.
func (in *InstanceList) DeepCopy() *InstanceList {
	if in == nil {
		return nil
	}
	out := new(InstanceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceSpec) DeepCopyInto(out *InstanceSpec) {
	*out = *in
	in.Network.DeepCopyInto(&out.Network)
	in.Configuration.DeepCopyInto(&out.Configuration)
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.Sidecars != nil {
		in, out := &in.Sidecars, &out.Sidecars
		*out = make([]v1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ImagePullSecrets != nil {
		in, out := &in.ImagePullSecrets, &out.ImagePullSecrets
		*out = make([]v1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
	if in.AllowPrivilegedPorts != nil {
		in, out := &in.AllowPrivilegedPorts, &out.AllowPrivilegedPorts
		*out = new(bool)
		**out = **in
	}
	if in.Placement != nil {
		in, out := &in.Placement, &out.Placement
		*out = new(Placement)
		(*in).DeepCopyInto(*out)
	}
	if in.Metrics != nil {
		in, out := &in.Metrics, &out.Metrics
		*out = new(Metrics)
		(*in).DeepCopyInto(*out)
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ReadinessProbe != nil {
		in, out := &in.ReadinessProbe, &out.ReadinessProbe
		*out = new(v1.Probe)
		(*in).DeepCopyInto(*out)
	}
	if in.LivenessProbe != nil {
		in, out := &in.LivenessProbe, &out.LivenessProbe
		*out = new(v1.Probe)
		(*in).DeepCopyInto(*out)
	}
	in.PodDisruptionBudget.DeepCopyInto(&out.PodDisruptionBudget)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceSpec.
func (in *InstanceSpec) DeepCopy() *InstanceSpec {
	if in == nil {
		return nil
	}
	out := new(InstanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceStatus) DeepCopyInto(out *InstanceStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceStatus.
func (in *InstanceStatus) DeepCopy() *InstanceStatus {
	if in == nil {
		return nil
	}
	out := new(InstanceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Metrics) DeepCopyInto(out *Metrics) {
	*out = *in
	if in.Address != nil {
		in, out := &in.Address, &out.Address
		*out = new(string)
		**out = **in
	}
	if in.RelabelConfigs != nil {
		in, out := &in.RelabelConfigs, &out.RelabelConfigs
		*out = make([]monitoringv1.RelabelConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Metrics.
func (in *Metrics) DeepCopy() *Metrics {
	if in == nil {
		return nil
	}
	out := new(Metrics)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Network) DeepCopyInto(out *Network) {
	*out = *in
	if in.HostIPs != nil {
		in, out := &in.HostIPs, &out.HostIPs
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Route.DeepCopyInto(&out.Route)
	in.Service.DeepCopyInto(&out.Service)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Network.
func (in *Network) DeepCopy() *Network {
	if in == nil {
		return nil
	}
	out := new(Network)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OcspUpdateOptionsHttpproxy) DeepCopyInto(out *OcspUpdateOptionsHttpproxy) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int64)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OcspUpdateOptionsHttpproxy.
func (in *OcspUpdateOptionsHttpproxy) DeepCopy() *OcspUpdateOptionsHttpproxy {
	if in == nil {
		return nil
	}
	out := new(OcspUpdateOptionsHttpproxy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Placement) DeepCopyInto(out *Placement) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.TopologySpreadConstraints != nil {
		in, out := &in.TopologySpreadConstraints, &out.TopologySpreadConstraints
		*out = make([]v1.TopologySpreadConstraint, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Placement.
func (in *Placement) DeepCopy() *Placement {
	if in == nil {
		return nil
	}
	out := new(Placement)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodDisruptionBudget) DeepCopyInto(out *PodDisruptionBudget) {
	*out = *in
	if in.MinAvailable != nil {
		in, out := &in.MinAvailable, &out.MinAvailable
		*out = new(intstr.IntOrString)
		**out = **in
	}
	if in.MaxUnavailable != nil {
		in, out := &in.MaxUnavailable, &out.MaxUnavailable
		*out = new(intstr.IntOrString)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodDisruptionBudget.
func (in *PodDisruptionBudget) DeepCopy() *PodDisruptionBudget {
	if in == nil {
		return nil
	}
	out := new(PodDisruptionBudget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RouteSpec) DeepCopyInto(out *RouteSpec) {
	*out = *in
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(routev1.TLSConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RouteSpec.
func (in *RouteSpec) DeepCopy() *RouteSpec {
	if in == nil {
		return nil
	}
	out := new(RouteSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceSpec) DeepCopyInto(out *ServiceSpec) {
	*out = *in
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(v1.ServiceType)
		**out = **in
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceSpec.
func (in *ServiceSpec) DeepCopy() *ServiceSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceSpec)
	in.DeepCopyInto(out)
	return out
}
