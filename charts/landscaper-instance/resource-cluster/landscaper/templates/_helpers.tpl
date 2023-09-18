{{/* SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Gardener contributors

 SPDX-License-Identifier: Apache-2.0
*/}}


{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "landscaper.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.AppVersion | replace "+" "_" | trunc 63 | trimSuffix "-" }}
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
