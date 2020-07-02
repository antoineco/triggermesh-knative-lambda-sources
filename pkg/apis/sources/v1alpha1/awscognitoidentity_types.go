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

	pkgapis "knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"

	"github.com/triggermesh/aws-event-sources/pkg/apis"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSCognitoIdentitySource is the Schema for the event source.
type AWSCognitoIdentitySource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSCognitoIdentitySourceSpec `json:"spec,omitempty"`
	Status AWSEventSourceStatus         `json:"status,omitempty"`
}

// Check the interfaces the event source should be implementing.
var (
	_ runtime.Object      = (*AWSCognitoIdentitySource)(nil)
	_ kmeta.OwnerRefable  = (*AWSCognitoIdentitySource)(nil)
	_ pkgapis.Validatable = (*AWSCognitoIdentitySource)(nil)
	_ pkgapis.Defaultable = (*AWSCognitoIdentitySource)(nil)
	_ pkgapis.HasSpec     = (*AWSCognitoIdentitySource)(nil)
	_ AWSEventSource      = (*AWSCognitoIdentitySource)(nil)
)

// AWSCognitoIdentitySourceSpec defines the desired state of the event source.
type AWSCognitoIdentitySourceSpec struct {
	duckv1.SourceSpec `json:",inline"`

	// Identity Pool ARN
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/list_amazoncognitoidentity.html#amazoncognitoidentity-resources-for-iam-policies
	ARN apis.ARN `json:"arn"`

	// Credentials to interact with the AWS Cognito API.
	Credentials AWSSecurityCredentials `json:"credentials"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSCognitoIdentitySourceList contains a list of event sources.
type AWSCognitoIdentitySourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSCognitoIdentitySource `json:"items"`
}
