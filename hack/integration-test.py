#!/usr/bin/env python3

# SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
#
# SPDX-License-Identifier: Apache-2.0

import os
import sys
import utils
import yaml
import json
import model.container_registry
import oci.auth as oa

from util import ctx
from subprocess import run

test_purpose = os.environ["TEST_PURPOSE"]
project_root = os.environ["PROJECT_ROOT"]
test_cluster = os.environ["TEST_CLUSTER"]
hosting_cluster = os.environ["HOSTING_CLUSTER"]
gardener_cluster = os.environ["GARDENER_CLUSTER"]
laas_version = os.environ["LAAS_VERSION"]
laas_repository = os.environ["LAAS_REPOSITORY"]

factory = ctx().cfg_factory()
print(f"Getting kubeconfig for {gardener_cluster}")
gardener_cluster_kubeconfig = factory.kubernetes(gardener_cluster)
print(f"Getting kubeconfig for {test_cluster}")
test_cluster_kubeconfig = utils.get_shoot_adminkubeconfig(test_cluster, gardener_cluster, "garden-laas")
print(f"Getting kubeconfig for {hosting_cluster}")
hosting_cluster_kubeconfig = utils.get_shoot_adminkubeconfig(hosting_cluster, gardener_cluster, "garden-laas")

print(f"Getting credentials for {repo_ctx_base_url}")
cr_conf = model.container_registry.find_config(repo_ctx_base_url, oa.Privileges.READONLY)

with (
    utils.TempFileAuto(prefix="test_cluster_kubeconfig_") as test_cluster_kubeconfig_temp_file,
    utils.TempFileAuto(prefix="hosting_cluster_kubeconfig_") as hosting_cluster_kubeconfig_temp_file,
    utils.TempFileAuto(prefix="gardener_cluster_kubeconfig_") as gardener_cluster_kubeconfig_temp_file,
    utils.TempFileAuto(prefix="registry_auth_", suffix=".json") as registry_temp_file
):
    test_cluster_kubeconfig_temp_file.write(test_cluster_kubeconfig)
    test_cluster_kubeconfig_path = test_cluster_kubeconfig_temp_file.switch()

    hosting_cluster_kubeconfig_temp_file.write(hosting_cluster_kubeconfig)
    hosting_cluster_kubeconfig_path = hosting_cluster_kubeconfig_temp_file.switch()

    gardener_cluster_kubeconfig_temp_file.write(yaml.safe_dump(gardener_cluster_kubeconfig.kubeconfig()))
    gardener_cluster_kubeconfig_path = gardener_cluster_kubeconfig_temp_file.switch()

    registry_temp_file.write(json.dumps(auths))
    registry_secrets_path = registry_temp_file.switch()

    command = ["go", "run", "./pkg/main.go",
                "--test-purpose", test_purpose,
                "--kubeconfig", test_cluster_kubeconfig_path,
                "--hosting-kubeconfig", hosting_cluster_kubeconfig_path,
                "--gardener-service-account-kubeconfig", gardener_cluster_kubeconfig_path,
                "--laas-version", laas_version,
                "--laas-repository", laas_repository]

    print(f"Running integration test with command: {' '.join(command)}")

    mod_path = os.path.join(project_root, "integration-test")
    run = run(command, cwd=mod_path)

    if run.returncode != 0:
        raise EnvironmentError("Integration test exited with errors")
