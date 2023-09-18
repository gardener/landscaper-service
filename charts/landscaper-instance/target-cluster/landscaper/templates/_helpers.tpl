{{/* SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Gardener contributors

 SPDX-License-Identifier: Apache-2.0
*/}}


{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "landscaper.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.AppVersion | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "landscaper.webhooks.name" -}}
{{- .Values.landscaper.name }}-webhooks
{{- end }}

{{- define "landscaper.main.name" -}}
{{- .Values.landscaper.name }}-main
{{- end }}

{{- define "landscaper.clusterrole.name" -}}
landscaper.gardener.cloud:{{- .Values.landscaper.name }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "landscaper.labels" -}}
helm.sh/chart: {{ include "landscaper.chart" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "landscaper.selectorLabels" -}}
landscaper.gardener.cloud/component: controller
app.kubernetes.io/name: {{ .Values.landscaper.name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "landscaper.main.selectorLabels" -}}
landscaper.gardener.cloud/component: controller-main
app.kubernetes.io/name: {{ .Values.landscaper.name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "landscaper.webhooks.selectorLabels" -}}
landscaper.gardener.cloud/component: webhook-server
app.kubernetes.io/name: {{ .Values.landscaper.name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "utils-templates.image" -}}
{{- if hasPrefix "sha256:" (required "$.tag is required" $.tag) -}}
{{ required "$.repository is required" $.repository }}@{{ required "$.tag is required" $.tag }}
{{- else -}}
{{ required "$.repository is required" $.repository }}:{{ required "$.tag is required" $.tag }}
{{- end -}}
{{- end -}}

{{- define "landscaper-image" -}}
{{- $tag := ( .Values.controller.image.tag | default .Chart.AppVersion )  -}}
{{- $image :=  dict "repository" .Values.controller.image.repository "tag" $tag  -}}
{{- include "utils-templates.image" $image }}
{{- end -}}

{{- define "landscaper-webhook-image" -}}
{{- $tag := ( .Values.webhooksServer.image.tag | default .Chart.AppVersion )  -}}
{{- $image :=  dict "repository" .Values.webhooksServer.image.repository "tag" $tag  -}}
{{- include "utils-templates.image" $image }}
{{- end -}}

{{- define "landscaper-config" -}}
apiVersion: config.landscaper.gardener.cloud/v1alpha1
kind: LandscaperConfiguration

registry:
  oci:
    allowPlainHttp: {{ .Values.landscaper.registryConfig.allowPlainHttpRegistries }}
    insecureSkipVerify: {{ .Values.landscaper.registryConfig.insecureSkipVerify }}
    {{- if .Values.landscaper.registryConfig.secrets }}
    configFiles:
    {{- range $key, $value := .Values.landscaper.registryConfig.secrets }}
    - /app/ls/registry/secrets/{{ $key }}
    {{- end }}
    {{- end }}
    cache:
      path: /app/ls/oci-cache/
      useInMemoryOverlay: {{ .Values.landscaper.registryConfig.cache.useInMemoryOverlay | default false }}


crdManagement:
  deployCrd: {{ .Values.landscaper.crdManagement.deployCrd }}
  {{- if .Values.landscaper.crdManagement.forceUpdate }}
  forceUpdate: {{ .Values.landscaper.crdManagement.forceUpdate }}
  {{- end }}

lsDeployments:
  lsController: "{{- .Values.landscaper.name }}"
  webHook: "{{- include "landscaper.webhooks.name" . }}"
  deploymentsNamespace: "{{ .Release.Namespace }}"
  lsHealthCheckName: "{{- .Values.landscaper.healthCheck.name }}"

{{- if .Values.landscaper.healthCheck.additionalDeployments }}
  additionalDeployments:
{{ toYaml .Values.landscaper.healthCheck.additionalDeployments | indent 4 }}
{{- end }}

{{- if .Values.controller.main.hpa }}
hpaMain:
{{ .Values.controller.main.hpa | toYaml | indent 2 }}
{{- end }}

{{- end }}