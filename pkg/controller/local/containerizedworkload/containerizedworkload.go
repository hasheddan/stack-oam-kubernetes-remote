/*
Copyright 2020 The Crossplane Authors.

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

package containerizedworkload

import (
	"reflect"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	oamv1alpha2 "github.com/crossplane/crossplane/apis/oam/v1alpha2"

	client "github.com/crossplane/addon-oam-kubernetes-remote/pkg/client/containerizedworkload"
	"github.com/crossplane/addon-oam-kubernetes-remote/pkg/reconciler/workload"
)

const (
	reconcileTimeout = 1 * time.Minute
	shortWait        = 30 * time.Second
	longWait         = 1 * time.Minute
)

// Reconcile error strings.
const (
	errNotContainerizedWorkload = "object is not a containerized workload"
)

const labelKey = "containerizedworkload.oam.crossplane.io"

var (
	deploymentKind             = reflect.TypeOf(appsv1.Deployment{}).Name()
	deploymentAPIVersion       = appsv1.SchemeGroupVersion.String()
	deploymentGroupVersionKind = appsv1.SchemeGroupVersion.WithKind(deploymentKind)
)

// SetupContainerizedWorkload adds a controller that reconciles ContainerizedWorkloads locally.
func SetupContainerizedWorkload(mgr ctrl.Manager, l logging.Logger) error {
	name := "local.oam/" + strings.ToLower(oamv1alpha2.ContainerizedWorkloadGroupKind)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&oamv1alpha2.ContainerizedWorkload{}).
		Complete(workload.NewReconciler(mgr,
			workload.Kind(oamv1alpha2.ContainerizedWorkloadGroupVersionKind),
			workload.WithLogger(l.WithValues("controller", name)),
			workload.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
			workload.WithPacker(workload.PackageFn(client.Packager)),
		))
}
