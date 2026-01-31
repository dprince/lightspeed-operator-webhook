/*
Copyright 2026.

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

package v1alpha1

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	upstreamolsv1alpha1 "github.com/openshift/lightspeed-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	olsv1alpha1 "github.com/openstack-lightspeed/openstack-lightspeed-operator/api/v1alpha1"
	// TODO (user): Add any additional imports if needed
)

var _ = Describe("OLSConfig Webhook", func() {
	var (
		obj       *olsv1alpha1.OLSConfig
		oldObj    *olsv1alpha1.OLSConfig
		defaulter OLSConfigCustomDefaulter
	)

	BeforeEach(func() {
		obj = &olsv1alpha1.OLSConfig{}
		oldObj = &olsv1alpha1.OLSConfig{}
		defaulter = OLSConfigCustomDefaulter{}
		Expect(defaulter).NotTo(BeNil(), "Expected defaulter to be initialized")
		Expect(oldObj).NotTo(BeNil(), "Expected oldObj to be initialized")
		Expect(obj).NotTo(BeNil(), "Expected obj to be initialized")
		// TODO (user): Add any setup logic common to all tests
	})

	AfterEach(func() {
		// TODO (user): Add any teardown logic common to all tests
	})

	Context("When creating OLSConfig under Defaulting Webhook", func() {
		It("Should apply OpenStackLightSpeed defaults when annotation is set", func() {
			By("setting the OpenStackLightSpeed annotation")
			obj.ObjectMeta = metav1.ObjectMeta{
				Name: "test-config",
				Annotations: map[string]string{
					OpenStackLightSpeedAnnotation: "true",
				},
			}

			By("setting the RAG image environment variable")
			testRAGImage := "quay.io/test/openstack-lightspeed-rag:latest"
			os.Setenv(OpenStackLightSpeedRAGImageEnv, testRAGImage)
			defer os.Unsetenv(OpenStackLightSpeedRAGImageEnv)

			By("calling the Default method")
			err := defaulter.Default(context.Background(), obj)
			Expect(err).NotTo(HaveOccurred())

			By("checking that byokRAGOnly is set to true")
			Expect(obj.Spec.OLSConfig.ByokRAGOnly).To(BeTrue())

			By("checking that querySystemPrompt is set")
			Expect(obj.Spec.OLSConfig.QuerySystemPrompt).To(Equal(openstackLightspeedSystemPrompt))
			Expect(obj.Spec.OLSConfig.QuerySystemPrompt).NotTo(BeEmpty())

			By("checking that RAG image is set")
			Expect(obj.Spec.OLSConfig.RAG).To(HaveLen(1))
			Expect(obj.Spec.OLSConfig.RAG[0].Image).To(Equal(testRAGImage))
		})

		It("Should not apply defaults when annotation is not set", func() {
			By("creating an OLSConfig without the annotation")
			obj.ObjectMeta = metav1.ObjectMeta{
				Name: "test-config",
			}

			By("calling the Default method")
			err := defaulter.Default(context.Background(), obj)
			Expect(err).NotTo(HaveOccurred())

			By("checking that defaults were not applied")
			Expect(obj.Spec.OLSConfig.ByokRAGOnly).To(BeFalse())
			Expect(obj.Spec.OLSConfig.QuerySystemPrompt).To(BeEmpty())
			Expect(obj.Spec.OLSConfig.RAG).To(BeEmpty())
		})

		It("Should not apply defaults when annotation value is not 'true'", func() {
			By("setting the annotation to a non-true value")
			obj.ObjectMeta = metav1.ObjectMeta{
				Name: "test-config",
				Annotations: map[string]string{
					OpenStackLightSpeedAnnotation: "false",
				},
			}

			By("calling the Default method")
			err := defaulter.Default(context.Background(), obj)
			Expect(err).NotTo(HaveOccurred())

			By("checking that defaults were not applied")
			Expect(obj.Spec.OLSConfig.ByokRAGOnly).To(BeFalse())
			Expect(obj.Spec.OLSConfig.QuerySystemPrompt).To(BeEmpty())
			Expect(obj.Spec.OLSConfig.RAG).To(BeEmpty())
		})

		It("Should skip RAG image when RAG is already configured", func() {
			By("setting the OpenStackLightSpeed annotation")
			obj.ObjectMeta = metav1.ObjectMeta{
				Name: "test-config",
				Annotations: map[string]string{
					OpenStackLightSpeedAnnotation: "true",
				},
			}

			By("pre-configuring a RAG entry")
			existingRAGImage := "quay.io/existing/rag:v1"
			obj.Spec.OLSConfig.RAG = []upstreamolsv1alpha1.RAGSpec{
				{
					Image: existingRAGImage,
				},
			}

			By("setting the RAG image environment variable")
			testRAGImage := "quay.io/test/openstack-lightspeed-rag:latest"
			os.Setenv(OpenStackLightSpeedRAGImageEnv, testRAGImage)
			defer os.Unsetenv(OpenStackLightSpeedRAGImageEnv)

			By("calling the Default method")
			err := defaulter.Default(context.Background(), obj)
			Expect(err).NotTo(HaveOccurred())

			By("checking that byokRAGOnly and querySystemPrompt are still set")
			Expect(obj.Spec.OLSConfig.ByokRAGOnly).To(BeTrue())
			Expect(obj.Spec.OLSConfig.QuerySystemPrompt).To(Equal(openstackLightspeedSystemPrompt))

			By("checking that existing RAG configuration is preserved")
			Expect(obj.Spec.OLSConfig.RAG).To(HaveLen(1))
			Expect(obj.Spec.OLSConfig.RAG[0].Image).To(Equal(existingRAGImage))
		})

		It("Should not set RAG image when environment variable is not set", func() {
			By("setting the OpenStackLightSpeed annotation")
			obj.ObjectMeta = metav1.ObjectMeta{
				Name: "test-config",
				Annotations: map[string]string{
					OpenStackLightSpeedAnnotation: "true",
				},
			}

			By("ensuring the environment variable is not set")
			os.Unsetenv(OpenStackLightSpeedRAGImageEnv)

			By("calling the Default method")
			err := defaulter.Default(context.Background(), obj)
			Expect(err).NotTo(HaveOccurred())

			By("checking that byokRAGOnly and querySystemPrompt are still set")
			Expect(obj.Spec.OLSConfig.ByokRAGOnly).To(BeTrue())
			Expect(obj.Spec.OLSConfig.QuerySystemPrompt).To(Equal(openstackLightspeedSystemPrompt))

			By("checking that RAG is not configured")
			Expect(obj.Spec.OLSConfig.RAG).To(BeEmpty())
		})
	})

})
