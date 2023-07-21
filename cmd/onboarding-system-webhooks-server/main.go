// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gardener/landscaper-service/cmd/onboarding-system-webhooks-server/app"
)

func main() {
	ctx := context.Background()
	defer ctx.Done()
	cmd := app.NewOnboardingSystemWebhooksCommand(ctx)

	if err := cmd.Execute(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
