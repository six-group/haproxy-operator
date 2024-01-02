package instance

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	proxyv1alpha1 "github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/uuid"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Reconcile", Label("controller"), func() {
	Context("Reconcile", func() {
		var (
			scheme   *runtime.Scheme
			ctx      context.Context
			proxy    *proxyv1alpha1.Instance
			initObjs []client.Object
		)

		BeforeEach(func() {
			scheme = runtime.NewScheme()
			Ω(clientgoscheme.AddToScheme(scheme)).ShouldNot(HaveOccurred())
			Ω(configv1alpha1.AddToScheme(scheme)).ShouldNot(HaveOccurred())
			Ω(proxyv1alpha1.AddToScheme(scheme)).ShouldNot(HaveOccurred())

			ctx = context.Background()

			labels := map[string]string{
				"label-test": "ok",
			}

			proxy = &proxyv1alpha1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar-foo",
					Namespace: "foo",
					UID:       uuid.NewUUID(),
				},
				Spec: proxyv1alpha1.InstanceSpec{
					Configuration: proxyv1alpha1.Configuration{
						Global: proxyv1alpha1.GlobalConfiguration{
							Logging: &proxyv1alpha1.GlobalLoggingConfiguration{
								Enabled:      true,
								Address:      "/var/lib/rsyslog/rsyslog.sock",
								Facility:     "local0",
								SendHostname: pointer.Bool(true),
							},
						},
					},
					Network: proxyv1alpha1.Network{
						HostIPs: map[string]string{
							"host1": "10.158.182.27",
						},
					},
					Labels: labels,
				},
			}

			initObjs = []client.Object{proxy}
		})

		It("create statefulset", func() {
			cli := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjs...).WithStatusSubresource(initObjs...).Build()
			r := Reconciler{
				Client: cli,
				Scheme: scheme,
			}
			err := r.reconcileStatefulSet(ctx, proxy)
			Ω(err).ShouldNot(HaveOccurred())

			statefulSet := &appsv1.StatefulSet{}
			Ω(cli.Get(ctx, client.ObjectKey{Namespace: proxy.Namespace, Name: "bar-foo-haproxy"}, statefulSet)).ShouldNot(HaveOccurred())
			Ω(statefulSet.Spec.Template.ObjectMeta.Labels["app.kubernetes.io/name"]).Should(Equal(proxy.Name + "-haproxy"))
			Ω(statefulSet.Spec.Template.ObjectMeta.Labels["label-test"]).Should(Equal("ok"))
			Ω(statefulSet.Spec.Template.Spec.InitContainers).Should(HaveLen(1))
			Ω(statefulSet.Spec.Template.Spec.InitContainers[0].Args[0]).Should(ContainSubstring("10.158.182.27"))
			Ω(statefulSet.Spec.Template.Spec.InitContainers[0].Args[0]).Should(Equal("\nif [ \"$HOSTNAME\" = \"host1\" ]\n" +
				"then\n  i=0\n  while [ $(ip a show to '10.158.182.27' | wc -l) -eq 0 ]\n  do\n    ((i=i+1))\n    if [ \"$i\" -gt \"20\" ]\n" +
				"      then echo 'timeout waiting for IP 10.158.182.27, aborting'\n      exit 1\n    fi\n    echo 'waiting for IP 10.158.182.27 to be assigned...'\n" +
				"    sleep 5\n  done\n\n  echo 'IP 10.158.182.27 assignment verified, waiting 5 seconds before continuing...'\n\n" +
				"  sleep 5\n\n  echo -n \"BIND_ADDRESS=10.158.182.27\" > /var/lib/haproxy/run/env\n  cat /var/lib/haproxy/run/env\n  exit 0\nfi\n\nexit 1\n"))
		})
	})
})
