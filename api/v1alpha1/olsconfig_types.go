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
	olsv1alpha1 "github.com/openshift/lightspeed-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Re-export types from the upstream OpenShift Lightspeed Operator CRD
// This acts as a shim to satisfy local API requirements

// OLSConfigSpec is an alias to the upstream OLSConfig spec
type OLSConfigSpec = olsv1alpha1.OLSConfigSpec

// OLSConfigStatus is an alias to the upstream OLSConfig status
type OLSConfigStatus = olsv1alpha1.OLSConfigStatus

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// OLSConfig is the Schema for the olsconfigs API.
type OLSConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OLSConfigSpec   `json:"spec,omitempty"`
	Status OLSConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OLSConfigList contains a list of OLSConfig.
type OLSConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OLSConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OLSConfig{}, &OLSConfigList{})
}
