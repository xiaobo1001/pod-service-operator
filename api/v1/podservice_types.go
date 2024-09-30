/*
Copyright 2024.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PodServiceSpec defines the desired state of PodService
type PodServiceSpec struct {
	Name           string `json:"name"`
	Replicas       int32  `json:"replicas"`
	Image          string `json:"image"`
	MaxSurge       int32  `json:"maxSurge"`
	MaxUnavailable int32  `json:"maxUnavailable"`
	Port           int32  `json:"port,omitempty"`
	ServicePort    int32  `json:"servicePort,omitempty"`
}

// PodServiceStatus defines the observed state of PodService
type PodServiceStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
	UpdatedReplicas   int32 `json:"updatedReplicas"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// PodService is the Schema for the podservices API
type PodService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodServiceSpec   `json:"spec,omitempty"`
	Status PodServiceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PodServiceList contains a list of PodService
type PodServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodService{}, &PodServiceList{})
}
