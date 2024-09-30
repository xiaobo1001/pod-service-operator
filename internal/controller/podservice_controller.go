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
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "wuleiwl/pod-service-operator/api/v1"
)

// PodServiceReconciler reconciles a PodService object
type PodServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=app.entrytask.com,resources=podservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.entrytask.com,resources=podservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=app.entrytask.com,resources=podservices/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodService object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *PodServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	// 获取 PodService 实例
	var podService appsv1.PodService
	if err := r.Get(ctx, req.NamespacedName, &podService); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 创建 Deployment 以管理 Pod 副本
	dep := &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podService.Name,
			Namespace: podService.Namespace,
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &podService.Spec.Replicas,
			Strategy: appv1.DeploymentStrategy{
				Type: appv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appv1.RollingUpdateDeployment{
					MaxSurge:       &intstr.IntOrString{IntVal: podService.Spec.MaxSurge},
					MaxUnavailable: &intstr.IntOrString{IntVal: podService.Spec.MaxUnavailable},
				},
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": podService.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": podService.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  podService.Spec.Name,
							Image: podService.Spec.Image,
							Ports: []corev1.ContainerPort{{ContainerPort: podService.Spec.Port}},
						},
					},
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(&podService, dep, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	found := &appv1.Deployment{}
	err := r.Get(ctx, req.NamespacedName, found)
	if err != nil {
		// 如果未找到，创建新的 Deployment
		if errors.IsNotFound(err) {
			if err := r.Create(ctx, dep); err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	if found != nil {
		if err := r.Update(ctx, dep); err != nil {
			return ctrl.Result{}, err
		}
	}

	// 创建 Service
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podService.Name,
			Namespace: podService.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": podService.Name},
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       podService.Spec.ServicePort,
				TargetPort: intstr.FromInt32(dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort),
			}},
		},
	}

	if err := controllerutil.SetControllerReference(&podService, svc, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	foundSvc := &corev1.Service{}
	svcErr := r.Get(ctx, req.NamespacedName, foundSvc)
	if svcErr != nil {
		// 如果未找到，创建新的 Service
		if errors.IsNotFound(err) {
			if err := r.Create(ctx, svc); err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	if foundSvc != nil {
		if err := r.Update(ctx, svc); err != nil {
			return ctrl.Result{}, err
		}
	}

	// 更新 PodService 状态
	podService.Status.AvailableReplicas = found.Status.AvailableReplicas
	podService.Status.UpdatedReplicas = found.Status.UpdatedReplicas
	if err := r.Status().Update(ctx, &podService); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.PodService{}).
		Complete(r)
}
