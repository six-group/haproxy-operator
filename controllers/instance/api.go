package instance

// InspectCluster will verify the availability of extra features available to the cluster, such as Prometheus and
// OpenShift Routes.
func InspectCluster() error {
	if err := VerifyPrometheusAPI(); err != nil {
		return err
	}

	return VerifyRouteAPI()
}
