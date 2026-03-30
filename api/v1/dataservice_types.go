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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DataServiceSpec defines the desired state of DataService
type DataServiceSpec struct {
	// replicas is the number of desired pod replicas
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`

	// image is the container image to run
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// configValue is an optional config string passed to the container
	// +optional
	ConfigValue *string `json:"configValue,omitempty"`
}

// DataServiceStatus defines the observed state of DataService.
type DataServiceStatus struct {
	// readyReplicas is the number of ready pods
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas"
// +kubebuilder:printcolumn:name="Ready",type="integer",JSONPath=".status.readyReplicas"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// DataService is the Schema for the dataservices API
type DataService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataServiceSpec   `json:"spec,omitempty"`
	Status DataServiceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DataServiceList contains a list of DataService
type DataServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DataService{}, &DataServiceList{})
}
