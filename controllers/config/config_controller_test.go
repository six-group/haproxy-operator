package config_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	proxyv1alpha1 "github.com/six-group/haproxy-operator/apis/proxy/v1alpha1"
	"github.com/six-group/haproxy-operator/controllers/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Reconcile", Label("controller"), func() {
	Context("Reconcile", func() {
		var scheme *runtime.Scheme

		BeforeEach(func() {
			scheme = runtime.NewScheme()
			Ω(clientgoscheme.AddToScheme(scheme)).ShouldNot(HaveOccurred())
			Ω(configv1alpha1.AddToScheme(scheme)).ShouldNot(HaveOccurred())
			Ω(proxyv1alpha1.AddToScheme(scheme)).ShouldNot(HaveOccurred())
		})
		It("should set owner reference for empty label selector", func() {
			proxy := &proxyv1alpha1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar-foo",
					Namespace: "foo",
					UID:       uuid.NewUUID(),
				},
				Spec: proxyv1alpha1.InstanceSpec{
					Configuration: proxyv1alpha1.Configuration{
						LabelSelector: metav1.LabelSelector{},
					},
				},
			}

			frontend := &configv1alpha1.Frontend{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "foo",
					Labels: map[string]string{
						"key2": "value2",
					},
				},
			}

			cli := fake.NewClientBuilder().WithScheme(scheme).WithObjects(proxy, frontend).Build()
			r := config.Reconciler{
				Client: cli,
				Scheme: scheme,
				Object: &configv1alpha1.Frontend{},
			}
			result, err := r.Reconcile(context.TODO(), ctrl.Request{NamespacedName: types.NamespacedName{Name: frontend.Name, Namespace: frontend.Namespace}})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).ShouldNot(BeNil())

			Ω(cli.Get(context.TODO(), client.ObjectKeyFromObject(frontend), frontend)).ShouldNot(HaveOccurred())
			Ω(frontend.OwnerReferences).ShouldNot(BeEmpty())
			Ω(frontend.OwnerReferences[0].UID).Should(Equal(proxy.UID))
			Ω(frontend.OwnerReferences[0].Name).Should(Equal(proxy.Name))
		})
		It("should set owner reference for non-empty label selector", func() {
			proxy := &proxyv1alpha1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar-foo",
					Namespace: "foo",
					UID:       uuid.NewUUID(),
				},
				Spec: proxyv1alpha1.InstanceSpec{
					Configuration: proxyv1alpha1.Configuration{
						LabelSelector: metav1.LabelSelector{
							MatchLabels: map[string]string{
								"key1": "value1",
							},
						},
					},
				},
			}

			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "foo",
					Labels: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
			}

			cli := fake.NewClientBuilder().WithScheme(scheme).WithObjects(proxy, listen).Build()
			r := config.Reconciler{
				Client: cli,
				Scheme: scheme,
				Object: &configv1alpha1.Listen{},
			}
			result, err := r.Reconcile(context.TODO(), ctrl.Request{NamespacedName: types.NamespacedName{Name: listen.Name, Namespace: listen.Namespace}})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).ShouldNot(BeNil())

			Ω(cli.Get(context.TODO(), client.ObjectKeyFromObject(listen), listen)).ShouldNot(HaveOccurred())
			Ω(listen.OwnerReferences).ShouldNot(BeEmpty())
			Ω(listen.OwnerReferences[0].UID).Should(Equal(proxy.UID))
			Ω(listen.OwnerReferences[0].Name).Should(Equal(proxy.Name))
		})
		It("should update error status if no instance", func() {
			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "foo",
				},
			}

			cli := fake.NewClientBuilder().WithScheme(scheme).WithObjects(listen).Build()
			r := config.Reconciler{
				Client: cli,
				Scheme: scheme,
				Object: &configv1alpha1.Listen{},
			}
			result, err := r.Reconcile(context.TODO(), ctrl.Request{NamespacedName: types.NamespacedName{Name: listen.Name, Namespace: listen.Namespace}})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).ShouldNot(BeNil())

			Ω(cli.Get(context.TODO(), client.ObjectKeyFromObject(listen), listen)).ShouldNot(HaveOccurred())
			Ω(listen.Status.Error).ShouldNot(BeNil())
			Ω(listen.Status.Phase).Should(Equal(configv1alpha1.StatusPhaseInternalError))
		})
		It("should update error status if instances do not match", func() {
			proxy := &proxyv1alpha1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar-foo",
					Namespace: "foo",
					UID:       uuid.NewUUID(),
				},
				Spec: proxyv1alpha1.InstanceSpec{
					Configuration: proxyv1alpha1.Configuration{
						LabelSelector: metav1.LabelSelector{
							MatchLabels: map[string]string{
								"key1": "value1",
							},
						},
					},
				},
			}

			listen := &configv1alpha1.Listen{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "foo",
					Labels: map[string]string{
						"key2": "value2",
					},
				},
			}

			cli := fake.NewClientBuilder().WithScheme(scheme).WithObjects(proxy, listen).Build()
			r := config.Reconciler{
				Client: cli,
				Scheme: scheme,
				Object: &configv1alpha1.Listen{},
			}
			result, err := r.Reconcile(context.TODO(), ctrl.Request{NamespacedName: types.NamespacedName{Name: listen.Name, Namespace: listen.Namespace}})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).ShouldNot(BeNil())

			Ω(cli.Get(context.TODO(), client.ObjectKeyFromObject(listen), listen)).ShouldNot(HaveOccurred())
			Ω(listen.Status.Error).ShouldNot(BeNil())
			Ω(listen.Status.Phase).Should(Equal(configv1alpha1.StatusPhaseInternalError))
		})
		It("should not update owner reference", func() {
			reference := metav1.OwnerReference{
				APIVersion: proxyv1alpha1.GroupVersion.String(),
				Kind:       "Instance",
				Name:       "foo",
				UID:        uuid.NewUUID(),
			}

			backend := &configv1alpha1.Backend{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "foo",
					Labels: map[string]string{
						"key2": "value2",
					},
					OwnerReferences: []metav1.OwnerReference{*reference.DeepCopy()},
				},
			}
			cli := fake.NewClientBuilder().WithScheme(scheme).WithObjects(backend).Build()
			r := config.Reconciler{
				Client: cli,
				Scheme: scheme,
				Object: &configv1alpha1.Backend{},
			}

			result, err := r.Reconcile(context.TODO(), ctrl.Request{NamespacedName: types.NamespacedName{Name: backend.Name, Namespace: backend.Namespace}})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).ShouldNot(BeNil())

			Ω(cli.Get(context.TODO(), client.ObjectKeyFromObject(backend), backend)).ShouldNot(HaveOccurred())
			Ω(backend.OwnerReferences).Should(ConsistOf(reference))
		})
	})
})
