{{- define "cart-api.name" -}}
cart-api
{{- end }}

{{- define "cart-api.fullname" -}}
{{ include "cart-api.name" . }}
{{- end }}
