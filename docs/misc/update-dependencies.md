<!--
SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"

SPDX-License-Identifier: Apache-2.0
-->
# Using Renovate Bot for Dependency Updates

To automate dependency updates, we use [Renovate Bot](https://docs.renovatebot.com/#why-use-renovate). Renovate regularly scans our project for outdated dependencies and automatically creates pull requests to update them and contributors will get notifications about. This helps keep your project secure and up-to-date with minimal manual effort.

Renovate can be configured via a `renovate.json` file in the repository, for all the different places where we reference or make usage of Free & Open Source Software.

Renovate provides a [Dependency Dashboard](https://github.com/gardener/landscaper-service/issues/328) where you can find details of configuration and open pull requests.
