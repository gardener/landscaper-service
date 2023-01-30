{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "ls-service-target-shoot-sidecar.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "ls-service-target-shoot-sidecar.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "ls-service-target-shoot-sidecar.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "ls-service-target-shoot-sidecar.controller.containerName" -}}
{{- if .Values.controller.containerName -}}
{{- .Values.controller.containerName | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- print "ls-service-target-shoot-sidecar-controller" }}
{{- end }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "ls-service-target-shoot-sidecar.labels" -}}
helm.sh/chart: {{ include "ls-service-target-shoot-sidecar.chart" . }}
{{ include "ls-service-target-shoot-sidecar.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "ls-service-target-shoot-sidecar.selectorLabels" -}}
ls-service-target-shoot-sidecar.gardener.cloud/component: controller
app.kubernetes.io/name: {{ include "ls-service-target-shoot-sidecar.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "ls-service-target-shoot-sidecar.controller.serviceAccountName" -}}
{{- default "ls-service-target-shoot-sidecar" .Values.serviceAccount.controller.name }}
{{- end }}

{{- define "ls-service-target-shoot-sidecar-controller-image" -}}
{{- $tag := ( .Values.controller.image.tag | default .Chart.AppVersion )  -}}
{{- $image :=  dict "repository" .Values.controller.image.repository "tag" $tag  -}}
{{- include "utils-templates.image" $image }}
{{- end -}}

{{- define "utils-templates.image" -}}
{{- if hasPrefix "sha256:" (required "$.tag is required" $.tag) -}}
{{ required "$.repository is required" $.repository }}@{{ required "$.tag is required" $.tag }}
{{- else -}}
{{ required "$.repository is required" $.repository }}:{{ required "$.tag is required" $.tag }}
{{- end -}}
{{- end -}}

{{- define "ls-service-target-shoot-sidecar-config" -}}
apiVersion: config.landscaper-service.gardener.cloud/v1alpha1
kind: TargetShootSidecarConfiguration

{{- if .Values.lsServiceTargetShootSidecar.metrics }}
metrics:
  port: {{ .Values.lsServiceTargetShootSidecar.metrics.port | default 8080 }}
{{- end }}

{{- if .Values.lsServiceTargetShootSidecar.crdManagement }}
crdManagement:
  deployCrd: {{ .Values.lsServiceTargetShootSidecar.crdManagement.deployCrd }}
  {{- if .Values.lsServiceTargetShootSidecar.crdManagement.forceUpdate }}
  forceUpdate: {{ .Values.lsServiceTargetShootSidecar.crdManagement.forceUpdate }}
  {{- end }}
{{- end }}

{{- end }}
