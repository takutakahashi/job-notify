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

package controller

import (
	"context"

	v1 "github.com/takutakahashi/job-notify/api/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// JobReconciler reconciles a Job object
type JobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=jobs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Job object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *JobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	nowJob := &batchv1.Job{}
	if err := r.Get(ctx, req.NamespacedName, nowJob); err != nil {
		logger.Error(err, "unable to fetch Job")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	jobNotifyList := &v1.JobNotifyList{}
	if err := r.List(ctx, jobNotifyList, nil); err != nil {
		logger.Error(err, "unable to fetch JobList")
		return ctrl.Result{}, err
	}
	for _, jobnotify := range jobNotifyList.Items {
		if IsMatchedJob(nowJob, jobnotify.GetLabels()) {
			logger.Info("matched job", "job", nowJob.Name)
			NotifyJob(ctx, r.Client, &jobnotify, nowJob)
		}
	}

	return ctrl.Result{}, nil
}

func NotifyJob(ctx context.Context, client client.Client, jobnotify *v1.JobNotify, job *batchv1.Job) {
	logger := log.FromContext(ctx)
	logger.Info("notify job", "job", job.Name)
	if jobnotify.Spec.Slack != nil {
	}
}

func IsMatchedJob(job *batchv1.Job, selector map[string]string) bool {
	if selector == nil {
		return true
	}
	jobLabels := job.GetLabels()
	if jobLabels == nil {
		return false
	}
	for key, value := range selector {
		if jobLabels[key] != value {
			return false
		}
	}
	return true
}

// SetupWithManager sets up the controller with the Manager.
func (r *JobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.Job{}).
		Complete(r)
}
