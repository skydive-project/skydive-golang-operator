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

// SkydiveFlowExporterSpec defines the desired state of SkydiveFlowExporter
type SkydiveFlowExporterSpec struct {
	Namespace string `json:"namespace,omitempty"`

	Deployment FlowExporterDeploymentSpec `json:"deployment"`

	// +optional
	// +kubebuilder:default=false
	DeployDevEnv bool `json:"deployDevEnv"`
}

type FlowExporterDeploymentSpec struct {
	// List of environment variables to set in the container.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Env []v1.EnvVar `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,7,rep,name=env"`
}

// SkydiveFlowExporterStatus defines the observed state of SkydiveFlowExporter
type SkydiveFlowExporterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// SkydiveFlowExporter is the Schema for the skydiveflowexporters API
type SkydiveFlowExporter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SkydiveFlowExporterSpec   `json:"spec,omitempty"`
	Status SkydiveFlowExporterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SkydiveFlowExporterList contains a list of SkydiveFlowExporter
type SkydiveFlowExporterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SkydiveFlowExporter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SkydiveFlowExporter{}, &SkydiveFlowExporterList{})
}
