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

package awscloudwatchlogssource

import (
	"context"

	"github.com/kelseyhightower/envconfig"

	"knative.dev/eventing/pkg/reconciler/source"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"

	"github.com/triggermesh/aws-event-sources/pkg/apis/sources/v1alpha1"
	informerv1alpha1 "github.com/triggermesh/aws-event-sources/pkg/client/generated/injection/informers/sources/v1alpha1/awscloudwatchlogssource"
	reconcilerv1alpha1 "github.com/triggermesh/aws-event-sources/pkg/client/generated/injection/reconciler/sources/v1alpha1/awscloudwatchlogssource"
	"github.com/triggermesh/aws-event-sources/pkg/reconciler/common"
)

// NewController creates a Reconciler for the event source and returns the result of NewImpl.
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {

	typ := (*v1alpha1.AWSCloudWatchLogsSource)(nil)
	app := common.AdapterName(typ)

	// Calling envconfig.Process() with a prefix appends that prefix
	// (uppercased) to the Go field name, e.g. MYSOURCE_IMAGE.
	adapterCfg := &adapterConfig{
		configs: source.WatchConfigurations(ctx, app, cmw, source.WithLogging, source.WithMetrics),
	}
	envconfig.MustProcess(app, adapterCfg)

	r := &Reconciler{
		adapterCfg: adapterCfg,
	}
	impl := reconcilerv1alpha1.NewImpl(ctx, r)

	r.base = common.NewGenericDeploymentReconciler(
		ctx,
		typ.GetGroupVersionKind(),
		impl.EnqueueKey,
		impl.EnqueueControllerOf,
	)

	informerv1alpha1.Get(ctx).Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	return impl
}
