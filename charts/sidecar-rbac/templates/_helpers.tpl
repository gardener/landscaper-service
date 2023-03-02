{{/* SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors

 SPDX-License-Identifier: Apache-2.0
*/}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "sidecar.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "sidecar.labels" -}}
helm.sh/chart: {{ include "sidecar.chart" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "sidecar.serviceAccountName" -}}
{{- default "sidecar" .Values.serviceAccount.name }}
{{- end }}

