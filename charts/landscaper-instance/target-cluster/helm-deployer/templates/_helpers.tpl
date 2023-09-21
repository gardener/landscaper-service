
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
{{ include "deployer.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "deployer.selectorLabels" -}}
app.kubernetes.io/name: {{ .Values.deployer.name}}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the Helm deployer config file which will be encapsulated in a secret.
*/}}
{{- define "deployer-config" -}}
apiVersion: helm.deployer.landscaper.gardener.cloud/v1alpha1
kind: Configuration
identity: {{ .Values.deployer.identity }}
{{- if .Values.deployer.registryConfig }}
oci:
  allowPlainHttp: {{ .Values.deployer.registryConfig.allowPlainHttp }}
  insecureSkipVerify: {{ .Values.deployer.registryConfig.insecureSkipVerify }}
  {{- if .Values.deployer.registryConfig.secrets }}
  configFiles:
  {{- range $key, $value := .Values.deployer.registryConfig.secrets }}
  - /app/ls/registry/secrets/{{ $key }}
  {{- end }}
  {{- end }}
{{- end }}
{{- if .Values.deployer.hpa }}
hpa:
{{ .Values.deployer.hpa | toYaml | indent 2 }}
{{- end }}
{{- if .Values.deployer.controller }}
controller:
{{ .Values.deployer.controller | toYaml | indent 2 }}
{{- end }}
{{- end }}

{{- define "deployer-image" -}}
{{- $tag := ( .Values.deployer.image.tag | default .Chart.AppVersion )  -}}
{{- $image :=  dict "repository" .Values.deployer.image.repository "tag" $tag  -}}
{{- include "utils-templates.image" $image }}
{{- end -}}

{{- define "utils-templates.image" -}}
{{- if hasPrefix "sha256:" (required "$.tag is required" $.tag) -}}
{{ required "$.repository is required" $.repository }}@{{ required "$.tag is required" $.tag }}
{{- else -}}
{{ required "$.repository is required" $.repository }}:{{ required "$.tag is required" $.tag }}
{{- end -}}
{{- end -}}