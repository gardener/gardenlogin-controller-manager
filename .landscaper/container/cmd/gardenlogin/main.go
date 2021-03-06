// Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"runtime"

	"github.com/gardener/gardenlogin-controller-manager/.landscaper/container/cmd/gardenlogin/app"
	"github.com/gardener/gardenlogin-controller-manager/.landscaper/container/internal/util"
)

func main() {
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	f := &util.FactoryImpl{}
	ctx := f.Context()

	if err := app.NewCommandGardenlogin(f).ExecuteContext(ctx); err != nil {
		panic(err)
	}
}
