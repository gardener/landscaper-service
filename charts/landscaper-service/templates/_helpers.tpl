{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "landscaper-service.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "landscaper-service.fullname" -}}
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

{{- define "landscaper-service.webhooks.fullname" -}}
{{- include "landscaper-service.fullname" . }}-webhooks
{{- end }}

{{- define "landscaper-service.agent.fullname" -}}
{{- include "landscaper-service.fullname" . }}-agent
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "landscaper-service.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "landscaper-service.labels" -}}
helm.sh/chart: {{ include "landscaper-service.chart" . }}
{{ include "landscaper-service.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "landscaper-service.selectorLabels" -}}
landscaper-service.gardener.cloud/component: controller
app.kubernetes.io/name: {{ include "landscaper-service.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "landscaper-service.webhooks.selectorLabels" -}}
landscaper-service.gardener.cloud/component: webhook-server
app.kubernetes.io/name: {{ include "landscaper-service.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "landscaper-service.controller.serviceAccountName" -}}
{{- default "landscaper-service" .Values.serviceAccount.controller.name }}
{{- end }}

{{- define "landscaper-service.webhooksServer.serviceAccountName" -}}
{{- default "landscaper-service-webhooks" .Values.serviceAccount.webhooksServer.name }}
{{- end }}

{{- define "landscaper-service-controller-image" -}}
{{- $tag := ( .Values.controller.image.tag | default .Values.image.tag | default .Chart.AppVersion )  -}}
{{- $image :=  dict "repository" .Values.controller.image.repository "tag" $tag  -}}
{{- include "utils-templates.image" $image }}
{{- end -}}

{{- define "landscaper-service-webhook-image" -}}
{{- $tag := ( .Values.webhooksServer.image.tag | default .Values.image.tag | default .Chart.AppVersion )  -}}
{{- $image :=  dict "repository" .Values.webhooksServer.image.repository "tag" $tag  -}}
{{- include "utils-templates.image" $image }}
{{- end -}}

{{- define "utils-templates.image" -}}
{{- if hasPrefix "sha256:" (required "$.tag is required" $.tag) -}}
{{ required "$.repository is required" $.repository }}@{{ required "$.tag is required" $.tag }}
{{- else -}}
{{ required "$.repository is required" $.repository }}:{{ required "$.tag is required" $.tag }}
{{- end -}}
{{- end -}}

{{- define "landscaper-service-config" -}}
apiVersion: config.landscaper-service.gardener.cloud/v1alpha1
kind: LandscaperServiceConfiguration

{{- if .Values.landscaperservice.metrics }}
metrics:
  port: {{ .Values.landscaperservice.metrics.port | default 8080 }}
{{- end }}

{{- if .Values.landscaperservice.crdManagement }}
crdManagement:
  deployCrd: {{ .Values.landscaperservice.crdManagement.deployCrd }}
  {{- if .Values.landscaperservice.crdManagement.forceUpdate }}
  forceUpdate: {{ .Values.landscaperservice.crdManagement.forceUpdate }}
  {{- end }}
{{- end }}

landscaperServiceComponent:
  name:  {{ .Values.landscaperservice.landscaperServiceComponent.name }}
  version: {{ .Values.landscaperservice.landscaperServiceComponent.version }}
  repositoryContext:
{{ toYaml .Values.landscaperservice.landscaperServiceComponent.repositoryContext | indent 4 }}
{{- if .Values.landscaperservice.landscaperServiceComponent.registryPullSecrets }}
  registryPullSecrets:
{{ toYaml .Values.landscaperservice.landscaperServiceComponent.registryPullSecrets | indent 4 }}
{{- end }}
{{- end }}
