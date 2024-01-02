package v1alpha1_test

import (
	parser "github.com/haproxytech/config-parser/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"time"
)

var simpleResolver = `
resolvers foo
  timeout resolve 1000
  timeout retry 1000
  resolve_retries 3
`

var _ = Describe("Resolver", Label("type"), func() {
	Context("AddToParser", func() {
		var p parser.Parser
		BeforeEach(func() {
			var err error
			p, err = parser.New()
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("should create resolver", func() {
			resolver := &configv1alpha1.Resolver{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
			}
			Ω(resolver.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(Equal(simpleResolver))
		})
		It("should set accepted payload size", func() {
			resolver := &configv1alpha1.Resolver{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ResolverSpec{
					AcceptedPayloadSize: ptr.To(int64(1024)),
				},
			}
			Ω(resolver.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("accepted_payload_size 1024\n"))
		})
		It("should set nameservers", func() {
			resolver := &configv1alpha1.Resolver{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ResolverSpec{
					Nameservers: []configv1alpha1.Nameserver{
						{Name: "ns1", Address: "ns1.com", Port: 53},
						{Name: "ns2", Address: "ns2.com", Port: 5553},
					},
				},
			}
			Ω(resolver.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("nameserver ns1 ns1.com:53\n"))
			Ω(p.String()).Should(ContainSubstring("nameserver ns2 ns2.com:5553\n"))
		})
		It("should set hold time periods", func() {
			resolver := &configv1alpha1.Resolver{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ResolverSpec{
					Hold: &configv1alpha1.Hold{
						Nx:       ptr.To(metav1.Duration{Duration: 1 * time.Second}),
						Obsolete: ptr.To(metav1.Duration{Duration: 2 * time.Second}),
						Other:    ptr.To(metav1.Duration{Duration: 3 * time.Second}),
						Refused:  ptr.To(metav1.Duration{Duration: 4 * time.Second}),
						Timeout:  ptr.To(metav1.Duration{Duration: 5 * time.Second}),
						Valid:    ptr.To(metav1.Duration{Duration: 6 * time.Second}),
					},
				},
			}
			Ω(resolver.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("hold nx 1000\n"))
			Ω(p.String()).Should(ContainSubstring("hold obsolete 2000\n"))
			Ω(p.String()).Should(ContainSubstring("hold other 3000\n"))
			Ω(p.String()).Should(ContainSubstring("hold refused 4000\n"))
			Ω(p.String()).Should(ContainSubstring("hold timeout 5000\n"))
			Ω(p.String()).Should(ContainSubstring("hold valid 6000\n"))
		})
		It("should overwrite default retires and timouts", func() {
			resolver := &configv1alpha1.Resolver{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.ResolverSpec{
					ResolveRetries: ptr.To(int64(10)),
					Timeouts: &configv1alpha1.Timeouts{
						Resolve: ptr.To(metav1.Duration{Duration: 2 * time.Second}),
						Retry:   ptr.To(metav1.Duration{Duration: 5 * time.Second}),
					},
				},
			}
			Ω(resolver.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(ContainSubstring("timeout resolve 2000\n"))
			Ω(p.String()).Should(ContainSubstring("timeout retry 5000\n"))
			Ω(p.String()).Should(ContainSubstring("resolve_retries 10\n"))
		})
	})
})
