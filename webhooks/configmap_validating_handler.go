/*
SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors

SPDX-License-Identifier: Apache-2.0
*/

package webhooks

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gardener/garden-login-controller-manager/internal/util"

	"github.com/go-logr/logr"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// ConfigmapValidator handles ConfigMap
type ConfigmapValidator struct {
	client      client.Client
	Log         logr.Logger
	Config      *util.ControllerManagerConfiguration
	configMutex sync.RWMutex

	// Decoder decodes objects
	decoder *admission.Decoder
}

func (h *ConfigmapValidator) getConfig() *util.ControllerManagerConfiguration {
	h.configMutex.RLock()
	defer h.configMutex.RUnlock()

	return h.Config
}

//// Mainly used for tests to inject config
//func (h *ConfigmapValidator) injectConfig(config *util.ControllerManagerConfiguration) {
//	h.configMutex.Lock()
//	defer h.configMutex.Unlock()
//
//	h.Config = config
//}

func (h *ConfigmapValidator) validatingKubeconfigConfigMapFn(ctx context.Context, t *corev1.ConfigMap, oldT *corev1.ConfigMap, admissionReq admissionv1.AdmissionRequest) (bool, string, error) {
	fldValidations := getFieldValidations(t)
	if err := validateRequiredFields(fldValidations); err != nil {
		return false, err.Error(), nil
	}
	// TODO

	return true, "allowed to be admitted", nil
}

type fldValidation struct {
	value   *string
	fldPath *field.Path
}

func getFieldValidations(t *corev1.ConfigMap) *[]fldValidation {
	kubeconfig := t.Data["kubeconfig"]
	fldValidations := &[]fldValidation{
		{
			value:   &kubeconfig,
			fldPath: field.NewPath("data", "kubeconfig"),
		},
	}

	return fldValidations
}

func validateRequiredFields(fldValidations *[]fldValidation) error {
	for _, fldValidation := range *fldValidations {
		if err := validateRequiredField(fldValidation.value, fldValidation.fldPath); err != nil {
			return err
		}
	}

	return nil
}

func validateRequiredField(val *string, fldPath *field.Path) error {
	if val == nil || len(*val) == 0 {
		return field.Required(fldPath, "field is required")
	}

	return nil
}

var _ admission.Handler = &ConfigmapValidator{}

// Handle handles admission requests.
func (h *ConfigmapValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	obj := &corev1.ConfigMap{}
	oldObj := &corev1.ConfigMap{}

	maxObjSize := h.getConfig().Webhooks.ConfigMapValidation.MaxObjectSize
	objSize := len(req.Object.Raw)

	if objSize > maxObjSize {
		err := fmt.Errorf("resource must not have more than %d bytes", maxObjSize)
		h.Log.Error(err, "maxObjectSize exceeded", "objSize", objSize, "maxObjSize", maxObjSize)

		return admission.Errored(http.StatusBadRequest, err)
	}

	err := h.decoder.Decode(req, obj)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if req.AdmissionRequest.Operation != admissionv1.Create {
		err = h.decoder.DecodeRaw(req.AdmissionRequest.OldObject, oldObj)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
	}

	allowed, reason, err := h.validatingKubeconfigConfigMapFn(ctx, obj, oldObj, req.AdmissionRequest)
	if err != nil {
		h.Log.Error(err, reason)
		return admission.Errored(http.StatusInternalServerError, err)
	}

	if !allowed {
		h.Log.Info("admission request denied", "reason", reason)
	}

	return admission.ValidationResponse(allowed, reason)
}

var _ inject.Client = &ConfigmapValidator{}

// A client will be automatically injected.

// InjectClient injects the client.
func (h *ConfigmapValidator) InjectClient(c client.Client) error {
	h.client = c
	return nil
}

// ConfigmapValidator implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (h *ConfigmapValidator) InjectDecoder(d *admission.Decoder) error {
	h.decoder = d
	return nil
}