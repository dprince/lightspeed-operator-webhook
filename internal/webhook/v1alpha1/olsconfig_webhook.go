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
	_ "embed"
	"fmt"
	"os"

	upstreamolsv1alpha1 "github.com/openshift/lightspeed-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	olsv1alpha1 "github.com/openstack-lightspeed/openstack-lightspeed-operator/api/v1alpha1"
)

const (
	// OpenStackLightSpeedAnnotation is the annotation that identifies this CR as an OpenStackLightSpeed instance
	OpenStackLightSpeedAnnotation = "openstack-lightspeed.openshift.io/enabled"
	// OpenStackLightSpeedRAGImageEnv is the environment variable containing the default RAG image
	OpenStackLightSpeedRAGImageEnv = "RELATED_IMAGE_OPENSTACK_LIGHTSPEED_RAG_IMAGE_URL_DEFAULT"
)

//go:embed system_prompt.txt
var openstackLightspeedSystemPrompt string

// nolint:unused
// log is for logging in this package.
var olsconfiglog = logf.Log.WithName("olsconfig-resource")

// SetupOLSConfigWebhookWithManager registers the webhook for OLSConfig in the manager.
func SetupOLSConfigWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&olsv1alpha1.OLSConfig{}).
		WithDefaulter(&OLSConfigCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-ols-openshift-io-v1alpha1-olsconfig,mutating=true,failurePolicy=fail,sideEffects=None,groups=ols.openshift.io,resources=olsconfigs,verbs=create;update,versions=v1alpha1,name=molsconfig-v1alpha1.kb.io,admissionReviewVersions=v1

// OLSConfigCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind OLSConfig when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type OLSConfigCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &OLSConfigCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind OLSConfig.
func (d *OLSConfigCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	olsconfig, ok := obj.(*olsv1alpha1.OLSConfig)

	if !ok {
		return fmt.Errorf("expected an OLSConfig object but got %T", obj)
	}
	olsconfiglog.Info("Defaulting for OLSConfig", "name", olsconfig.GetName())

	// Check if this is an OpenStackLightSpeed instance
	if enabled, exists := olsconfig.Annotations[OpenStackLightSpeedAnnotation]; exists && enabled == "true" {
		olsconfiglog.Info("Applying OpenStackLightSpeed defaults", "name", olsconfig.GetName())

		// Set byokRAGOnly to true
		olsconfig.Spec.OLSConfig.ByokRAGOnly = true

		// Set querySystemPrompt from embedded file
		olsconfig.Spec.OLSConfig.QuerySystemPrompt = openstackLightspeedSystemPrompt

		// Set ragImage only if RAG is not already configured
		if len(olsconfig.Spec.OLSConfig.RAG) == 0 {
			ragImage := os.Getenv(OpenStackLightSpeedRAGImageEnv)
			if ragImage != "" {
				olsconfig.Spec.OLSConfig.RAG = []upstreamolsv1alpha1.RAGSpec{
					{
						Image: ragImage,
					},
				}
				olsconfiglog.Info("Set OpenStackLightSpeed RAG image", "image", ragImage)
			} else {
				olsconfiglog.Info("OpenStackLightSpeed RAG image environment variable not set",
					"envVar", OpenStackLightSpeedRAGImageEnv)
			}
		} else {
			olsconfiglog.Info("RAG already configured, skipping RAG image default")
		}
	}

	return nil
}
