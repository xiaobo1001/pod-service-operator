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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var podservicelog = logf.Log.WithName("podservice-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *PodService) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-app-entrytask-com-v1-podservice,mutating=true,failurePolicy=fail,sideEffects=None,groups=app.entrytask.com,resources=podservices,verbs=create;update,versions=v1,name=mpodservice.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &PodService{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *PodService) Default() {
	podservicelog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
	if r.Spec.ServicePort == 0 {
		r.Spec.ServicePort = 8086
		podservicelog.Info("default service port", "servicePort", r.Spec.ServicePort)
	} else {
		podservicelog.Info("service port exists", "servicePort", r.Spec.ServicePort)
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-app-entrytask-com-v1-podservice,mutating=false,failurePolicy=fail,sideEffects=None,groups=app.entrytask.com,resources=podservices,verbs=create;update,versions=v1,name=vpodservice.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &PodService{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *PodService) ValidateCreate() (admission.Warnings, error) {
	podservicelog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil, r.validatePort()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *PodService) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	podservicelog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil, r.validatePort()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *PodService) ValidateDelete() (admission.Warnings, error) {
	podservicelog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

func (r *PodService) validatePort() error {
	if r.Spec.ServicePort < 8000 || r.Spec.ServicePort > 9999 {
		podservicelog.Info("validate service port failed", "servicePort", r.Spec.ServicePort)
		return field.Invalid(field.NewPath("spec").Child("servicePort"), r.Spec.ServicePort, "port must be in range 8000-9999")
	}
	return nil
}
