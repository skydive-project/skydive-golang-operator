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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SkydiveSpec defines the desired state of Skydive
type SkydiveSpec struct {
	Namespace string `json:"namespace,omitempty"`

	Enable EnableSpec `json:"enable"`

	Agents AgentsSpec `json:"agents"`

	Analyzer AnalyzerSpec `json:"analyzer"`
}

type AgentsSpec struct {
	DaemonSet AgentsDaemonSetSpec `json:"daemonSet"`
}

type AnalyzerSpec struct {
	Deployment AnalyzerDeploymentSpec `json:"deployment"`
}

type AgentsDaemonSetSpec struct {

	// List of environment variables to set in the container.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Env []v1.EnvVar `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,7,rep,name=env"`
}

type AnalyzerDeploymentSpec struct {

	// List of environment variables to set in the container.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Env []v1.EnvVar `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,7,rep,name=env"`
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

// SkydiveStatus defines the observed state of Skydive
type SkydiveStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// Skydive is the Schema for the skydives API
type Skydive struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SkydiveSpec   `json:"spec,omitempty"`
	Status SkydiveStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SkydiveList contains a list of Skydive
type SkydiveList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Skydive `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Skydive{}, &SkydiveList{})
}
