/*
Copyright (c) 2020 TriggerMesh Inc.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	duckv1 "knative.dev/pkg/apis/duck/v1"

	"github.com/triggermesh/aws-event-sources/pkg/apis"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSCloudWatchSource is the Schema for the event source.
type AWSCloudWatchLogSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSCloudWatchLogSourceSpec `json:"spec,omitempty"`
	Status EventSourceStatus          `json:"status,omitempty"`
}

// Check the interfaces the event source should be implementing.
var (
	_ runtime.Object = (*AWSCloudWatchLogSource)(nil)
	_ EventSource    = (*AWSCloudWatchLogSource)(nil)
)

// AWSCloudWatchSourceSpec defines the desired state of the event source.
type AWSCloudWatchLogSourceSpec struct {
	duckv1.SourceSpec `json:",inline"`

	// ARN for Log Group
	// https://docs.aws.amazon.com/service-authorization/latest/reference/list_amazoncloudwatchlogs.html#amazoncloudwatchlogs-resources-for-iam-policies
	ARN apis.ARN `json:"arn"`
	// PollingFrequency in a duration format for how often to pull metrics data from. Default is 5m
	// +optional
	PollingFrequency *string `json:"pollingFrequency,omitempty"`

	// Credentials to interact with the AWS CloudWatch Logs API.
	Credentials AWSSecurityCredentials `json:"credentials"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSCloudWatchLogSourceList contains a list of event sources.
type AWSCloudWatchLogSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSCloudWatchLogSource `json:"items"`
}
