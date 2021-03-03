/*


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

// SkydiveSuiteSpec defines the desired state of SkydiveSuite
type SkydiveSuiteSpec struct {
	Namespace string `json:"namespace,omitempty"`

	Enable EnableSpec `json:"enable"`

	Logging LoggingSpec `json:"logging"`
}

type EnableSpec struct {
	// +optional
	// +kubebuilder:default=true
	Analyzer bool `json:"analyzer,omitempty"`

	// +optional
	// +kubebuilder:default=true
	Agents bool `json:"agents,omitempty"`

	// +optional
	// +kubebuilder:default=false
	Route bool `json:"route,omitempty"`
}

type LoggingSpec struct {
	// +optional
	// +kubebuilder:default="INFO"
	Level string `json:"level,omitempty"`
}

// SkydiveSuiteStatus defines the observed state of SkydiveSuite
type SkydiveSuiteStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// SkydiveSuite is the Schema for the skydivesuites API
type SkydiveSuite struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SkydiveSuiteSpec   `json:"spec,omitempty"`
	Status SkydiveSuiteStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SkydiveSuiteList contains a list of SkydiveSuite
type SkydiveSuiteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SkydiveSuite `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SkydiveSuite{}, &SkydiveSuiteList{})
}
