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
	"k8s.io/apimachinery/pkg/runtime/schema"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	"github.com/triggermesh/aws-event-sources/pkg/apis"
)

// GetGroupVersionKind implements kmeta.OwnerRefable.
func (s *AWSSQSSource) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("AWSSQSSource")
}

// GetUntypedSpec implements apis.HasSpec.
func (s *AWSSQSSource) GetUntypedSpec() interface{} {
	return s.Spec
}

// GetSink implements AWSEventSource.
func (s *AWSSQSSource) GetSink() *duckv1.Destination {
	return &s.Spec.Sink
}

// GetARN implements AWSEventSource.
func (s *AWSSQSSource) GetARN() apis.ARN {
	return s.Spec.ARN
}

// GetStatus implements AWSEventSource.
func (s *AWSSQSSource) GetStatus() *AWSEventSourceStatus {
	return &s.Status
}

// Supported event types
const (
	AWSSQSGenericEventType = "message"
)

// GetEventTypes implements AWSEventSource.
func (s *AWSSQSSource) GetEventTypes() []string {
	return []string{
		AWSSQSGenericEventType,
	}
}
