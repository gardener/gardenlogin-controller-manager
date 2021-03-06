// Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package api

import "path/filepath"

// Contents defines the structure for the used content data of the landscaper component.
type Contents struct {
	// DefaultPath holds the path of the "default" folder in which the kube-rbac-proxy image is defined
	DefaultPath string
	// ManagerPath holds the path of the "manager" folder in which the gardenlogin-controller-manager image is defined
	ManagerPath string

	// GardenloginTLSPath holds the path of the "tls" folder in which the tls certificate files are placed for kustomize's secretGenerator to pick them up
	GardenloginTLSPath string
	// GardenloginTLSPemFile holds the file path of the gardenlogin-controller-manager-tls.pem file for the webhook server
	GardenloginTLSPemFile string
	// GardenloginTLSKeyPemFile holds the file path of the gardenlogin-controller-manager-tls-key.pem file for the webhook server
	GardenloginTLSKeyPemFile string

	// RuntimeManagerPath is the path to the manager directory of the runtime overlay
	RuntimeManagerPath string
	// GardenloginKubeconfigPath holds the file path of the kubeconfig for the gardenlogin-controller-manager
	GardenloginKubeconfigPath string

	// Kustomize Overlay Paths

	// VirtualGardenOverlayPath holds the path of the virtual garden kustomize overlay
	VirtualGardenOverlayPath string
	// RuntimeOverlayPath holds the path of the runtime-cluster kustomize overlay
	RuntimeOverlayPath string
	// SingleClusterPath holds the path of the single-cluster kustomize overlay
	SingleClusterPath string

	// ManagerConfigurationRuntimePath holds the path of the ControllerManagerConfiguration under the runtime overlay
	ManagerConfigurationRuntimePath string
	// ManagerConfigurationSingleClusterPath holds the path of the ControllerManagerConfiguration under the single-cluster overlay
	ManagerConfigurationSingleClusterPath string
}

// NewContentsFromPath returns Contents struct for the given contentPath
func NewContentsFromPath(contentPath string) *Contents {
	contents := &Contents{
		DefaultPath: filepath.Join(contentPath, "config", "default"),
		ManagerPath: filepath.Join(contentPath, "config", "manager"),

		GardenloginTLSPath:       filepath.Join(contentPath, "config", "secret", "tls"),
		GardenloginTLSPemFile:    filepath.Join(contentPath, "config", "secret", "tls", "gardenlogin-controller-manager-tls.pem"),
		GardenloginTLSKeyPemFile: filepath.Join(contentPath, "config", "secret", "tls", "gardenlogin-controller-manager-tls-key.pem"),

		RuntimeManagerPath:        filepath.Join(contentPath, "config", "overlay", "multi-cluster", "runtime", "manager"),
		GardenloginKubeconfigPath: filepath.Join(contentPath, "config", "overlay", "multi-cluster", "runtime", "manager", "kubeconfig.yaml"),

		VirtualGardenOverlayPath: filepath.Join(contentPath, "config", "overlay", "multi-cluster", "virtual-garden"),
		RuntimeOverlayPath:       filepath.Join(contentPath, "config", "overlay", "multi-cluster", "runtime"),
		SingleClusterPath:        filepath.Join(contentPath, "config", "overlay", "single-cluster"),

		ManagerConfigurationRuntimePath:       filepath.Join(contentPath, "config", "overlay", "multi-cluster", "runtime", "manager", "config.yaml"),
		ManagerConfigurationSingleClusterPath: filepath.Join(contentPath, "config", "overlay", "single-cluster", "manager", "config.yaml"),
	}

	return contents
}
