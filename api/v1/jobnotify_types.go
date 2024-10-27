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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JobNotifySpec defines the desired state of JobNotify
type JobNotifySpec struct {
	Slack       NotifySlack          `json:"slack,omitempty"`
	JobSelector metav1.LabelSelector `json:"jobSelector,omitempty"`
}

type NotifySlack struct {
	WebhookURL corev1.EnvVarSource `json:"webhookURL,omitempty"`
	Channel    corev1.EnvVarSource `json:"channel,omitempty"`
}

// JobNotifyStatus defines the observed state of JobNotify
type JobNotifyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// JobNotify is the Schema for the jobnotifies API
type JobNotify struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JobNotifySpec   `json:"spec,omitempty"`
	Status JobNotifyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JobNotifyList contains a list of JobNotify
type JobNotifyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JobNotify `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JobNotify{}, &JobNotifyList{})
}
