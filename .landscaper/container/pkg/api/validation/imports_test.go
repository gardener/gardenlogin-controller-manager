// Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package validation_test

import (
	"encoding/json"

	"github.com/gardener/gardenlogin-controller-manager/.landscaper/container/pkg/api"
	. "github.com/gardener/gardenlogin-controller-manager/.landscaper/container/pkg/api/validation"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var _ = Describe("Imports", func() {
	Describe("#ValidateImports", func() {
		var (
			obj *api.Imports
		)

		BeforeEach(func() {
			obj = &api.Imports{
				RuntimeClusterTarget: lsv1alpha1.Target{
					Spec: lsv1alpha1.TargetSpec{
						Configuration: lsv1alpha1.AnyJSON{
							RawMessage: json.RawMessage(`{"config":{"kubeconfig":"foo1"}}`),
						},
					},
				},
				ApplicationClusterTarget: lsv1alpha1.Target{
					Spec: lsv1alpha1.TargetSpec{
						Configuration: lsv1alpha1.AnyJSON{
							RawMessage: json.RawMessage(`{"config":{"kubeconfig":"foo2"}}`),
						},
					},
				},
			}
		})

		It("should pass for a valid configuration", func() {
			Expect(ValidateImports(obj)).To(BeEmpty())
		})

		Context("hosting cluster", func() {
			It("should fail for an invalid configuration", func() {
				obj.RuntimeClusterTarget = lsv1alpha1.Target{}
				obj.ApplicationClusterTarget = lsv1alpha1.Target{}

				Expect(ValidateImports(obj)).To(ConsistOf(
					PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeRequired),
						"Field": Equal("runtimeClusterTarget"),
					})),
					PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeRequired),
						"Field": Equal("applicationClusterTarget"),
					})),
				))
			})
		})
	})
})
