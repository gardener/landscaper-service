# SPDX-FileCopyrightText: 2022 "SAP SE or an SAP affiliate company and Gardener contributors"
#
# SPDX-License-Identifier: Apache-2.0

import os
import tempfile
import base64
import yaml
import json
import util
import subprocess

class TempFileAuto(object):
    def __init__(self, prefix=None, mode='w+', suffix=".yaml"):
        self.file_obj = tempfile.NamedTemporaryFile(mode=mode, prefix=prefix, suffix=suffix, delete=False)
        self.name = self.file_obj.name
    def __enter__(self):
        return self
    def write(self, b):
        self.file_obj.write(b)
    def writelines(self, lines):
        self.file_obj.writelines(lines)
    def switch(self):
        self.file_obj.close()
        return self.file_obj.name
    def __exit__(self, type, value, traceback):
        if not self.file_obj.closed:
            self.file_obj.close()
        os.remove(self.file_obj.name)
        return False

def base64_encode_to_string(value):
    value = value.rstrip()
    if isinstance(value, str):
        value = value.encode('utf-8')
    return base64.b64encode(value).decode('utf-8')

def write_data(path: str, data: str):
    with open(path, "w") as file:
        file.write(data)

def run_command(args: any, cwd: str = None, silent: bool = False):
    if isinstance(args, str):
        command = args.split(' ')
    else:
        command = args
    print(f'running command "{args}"')
    if silent:
        result = subprocess.run(command, cwd=cwd, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    else:
        result = subprocess.run(command, cwd=cwd)
    if result.returncode != 0:
        raise RuntimeError(f'Could not run command "{args}"')

# see https://github.com/gardener/gardener/blob/master/docs/usage/shoot_access.md#shootsadminkubeconfig-subresource
# expirationSeconds: 86400 is 1 day
def get_shoot_adminkubeconfig(shoot_name: str, service_account_name: str, namespace: str, expiration_seconds=86400):
    with tempfile.TemporaryDirectory() as tmpdir:
        print(f'Getting kubeconfig for service account {service_account_name}, expiration_seconds={expiration_seconds}')
        factory = util.ctx().cfg_factory()
        service_account = factory.kubernetes(service_account_name)
        service_account_kubeconfig_path = os.path.join(tmpdir, 'service_account_kubeconfig')
        print(f'DEBUG laas_admin_core_kubeconfig_path={service_account_kubeconfig_path}')
        write_data(
            service_account_kubeconfig_path,
            yaml.safe_dump(service_account.kubeconfig()),
        )

        admin_kube_config_request = f'{{"apiVersion": "authentication.gardener.cloud/v1alpha1", "kind": "AdminKubeconfigRequest", "spec": {{"expirationSeconds": {expiration_seconds}}}}}'

        write_data(
            'AdminKubeconfigRequest.json',
            admin_kube_config_request,
        )

        print(f'Getting shoots/adminkubeconfig subresource for {shoot_name} in namespace {namespace}')
        command = f'kubectl --kubeconfig={service_account_kubeconfig_path} create --raw /apis/core.gardener.cloud/v1beta1/namespaces/{namespace}/shoots/{shoot_name}/adminkubeconfig -f AdminKubeconfigRequest.json'

        rc = run_command(command)
        rc_json = json.loads(rc)
        kubeconfig_bytes = base64.b64decode(rc_json["status"]["kubeconfig"])
        kubeconfig = kubeconfig_bytes.decode('utf-8')

        return kubeconfig
