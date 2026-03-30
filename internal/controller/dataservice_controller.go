/*
Copyright 2026.

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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	appv1 "github.com/marius0711/k8s-operator-lab/api/v1"
)

type DataServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=app.schwarz.io,resources=dataservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.schwarz.io,resources=dataservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=app.schwarz.io,resources=dataservices/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

func (r *DataServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// 1. CR laden
	var ds appv1.DataService
	if err := r.Get(ctx, req.NamespacedName, &ds); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	log.Info("Reconciling DataService", "name", ds.Name, "replicas", ds.Spec.Replicas)

	deploymentName := fmt.Sprintf("%s-deployment", ds.Name)
	labels := map[string]string{
		"app":        ds.Name,
		"controller": "dataservice",
	}

	desired := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: ds.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &ds.Spec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "dataservice",
							Image: ds.Spec.Image,
						},
					},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(&ds, desired, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	var existing appsv1.Deployment
	err := r.Get(ctx, types.NamespacedName{Name: deploymentName, Namespace: ds.Namespace}, &existing)
	if errors.IsNotFound(err) {
		log.Info("Creating Deployment", "name", deploymentName)
		if err := r.Create(ctx, desired); err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		existing.Spec.Replicas = &ds.Spec.Replicas
		existing.Spec.Template.Spec.Containers[0].Image = ds.Spec.Image
		if err := r.Update(ctx, &existing); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Fix: nur updaten wenn sich ReadyReplicas wirklich geändert hat
	currentReady := existing.Status.ReadyReplicas
	if ds.Status.ReadyReplicas != currentReady {
		ds.Status.ReadyReplicas = currentReady
		ds.Status.Conditions = []metav1.Condition{{
			Type:               "Available",
			Status:             metav1.ConditionTrue,
			Reason:             "DeploymentReconciled",
			Message:            fmt.Sprintf("Deployment %s reconciled", deploymentName),
			LastTransitionTime: metav1.Now(),
		}}
		if err := r.Status().Update(ctx, &ds); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *DataServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.DataService{}).
		Owns(&appsv1.Deployment{}).
		Named("dataservice").
		Complete(r)
}
