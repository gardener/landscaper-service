// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gardener/landscaper-service/cmd/landscaper-service-webhooks-server/app"
)

func main() {
	ctx := context.Background()
	defer ctx.Done()
	cmd := app.NewLandscaperServiceWebhooksCommand(ctx)

	if err := cmd.Execute(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
