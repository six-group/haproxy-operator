package v1alpha1_test

import (
	parser "github.com/haproxytech/config-parser/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	configv1alpha1 "github.com/six-group/haproxy-operator/apis/config/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var simpleFrontend = `
frontend foo
`

var withBackendRule = `
frontend foo
  use_backend %[base,map_reg(/usr/local/etc/haproxy/mymap.map)] if { base,map_reg(/usr/local/etc/haproxy/mymap.map) -m found }
`

var _ = Describe("Frontend", Label("type"), func() {
	Context("AddToParser", func() {
		var p parser.Parser
		BeforeEach(func() {
			var err error
			p, err = parser.New()
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should create frontend", func() {
			frontend := &configv1alpha1.Frontend{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
			}
			Ω(frontend.AddToParser(p)).ShouldNot(HaveOccurred())
			Ω(p.String()).Should(Equal(simpleFrontend))
		})

		It("should create map_reg", func() {
			backend := configv1alpha1.BackendReference{
				RegexMapping: &configv1alpha1.RegexBackendMapping{
					Name:      "mymap",
					Parameter: "base",
				},
			}

			frontend := &configv1alpha1.Frontend{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: configv1alpha1.FrontendSpec{
					BackendSwitching: []configv1alpha1.BackendSwitchingRule{
						{
							Rule: configv1alpha1.Rule{
								ConditionType: "if",
								Condition:     backend.RegexMapping.FoundCondition(),
							},
							Backend: backend,
						},
					},
				},
			}
			Ω(frontend.AddToParser(p)).ShouldNot(HaveOccurred())
			a := p.String()
			Ω(a).Should(Equal(withBackendRule))
		})
	})
})
