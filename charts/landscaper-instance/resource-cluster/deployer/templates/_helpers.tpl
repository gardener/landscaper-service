{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "deployer.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.AppVersion | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "deployer.labels" -}}
helm.sh/chart: {{ include "deployer.chart" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}