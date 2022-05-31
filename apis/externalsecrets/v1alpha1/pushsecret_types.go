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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PushSecretStoreRef struct {
	// Name of the SecretStore resource
	Name string `json:"name"`

	// Kind of the SecretStore resource (SecretStore or ClusterSecretStore)
	// Defaults to `SecretStore`
	// +optional
	Kind string `json:"kind,omitempty"`
}

// PushSecretSpec configures the behavior of the PushSecret.
type PushSecretSpec struct {
	SecretStoreRefs []PushSecretStoreRef `json:"secretStoreRefs"`
	Selector        PushSecretSelector   `json:"selector"`
	Data            []PushSecretData     `json:"data,omitempty"`
}

type PushSecretSecret struct {
	Name string `json:"name"`
}

type PushSecretSelector struct {
	Secret PushSecretSecret `json:"secret"`
}

type PushSecretRemoteRefs struct {
	RemoteKey string `json:"remoteKey"`
}

func (r PushSecretRemoteRefs) GetRemoteKey() string {
	return r.RemoteKey
}

type PushSecretMatch struct {
	SecretKey  string                 `json:"secretKey"`
	RemoteRefs []PushSecretRemoteRefs `json:"remoteRefs"`
}

type PushSecretData struct {
	Match []PushSecretMatch `json:"match"`
}

// PushSecretConditionType indicates the condition of the PushSecret.
type PushSecretConditionType string

const (
	PushSecretReady PushSecretConditionType = "Ready"
)

// PushSecretStatusCondition indicates the status of the PushSecret.
type PushSecretStatusCondition struct {
	Type   PushSecretConditionType `json:"type"`
	Status corev1.ConditionStatus  `json:"status"`

	// +optional
	Reason string `json:"reason,omitempty"`

	// +optional
	Message string `json:"message,omitempty"`

	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

// PushSecretStatus indicates the history of the status of PushSecret.
type PushSecretStatus struct {
	// +nullable
	// refreshTime is the time and date the external secret was fetched and
	// the target secret updated
	RefreshTime metav1.Time `json:"refreshTime,omitempty"`

	// SyncedResourceVersion keeps track of the last synced version.
	SyncedResourceVersion string `json:"syncedResourceVersion,omitempty"`

	// +optional
	Conditions []PushSecretStatusCondition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// PushSecrets is the Schema for the PushSecrets API.
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].reason`
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={PushSecrets}

type PushSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PushSecretSpec   `json:"spec,omitempty"`
	Status PushSecretStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].reason`
// PushSecretList contains a list of PushSecret resources.
type PushSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PushSecret `json:"items"`
}