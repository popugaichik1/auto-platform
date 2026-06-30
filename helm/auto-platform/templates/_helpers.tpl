{{/*
Имя приложения-сервиса (auth/listing/user) -> "<release>-<key>"
Использование: {{ include "auto-platform.appName" (dict "root" $ "key" $key) }}
*/}}
{{- define "auto-platform.appName" -}}
{{ .root.Release.Name }}-{{ .key }}
{{- end -}}

{{/*
Имя StatefulSet/Service для Postgres конкретного сервиса -> "<release>-<key>-postgres"
*/}}
{{- define "auto-platform.postgresName" -}}
{{ .root.Release.Name }}-{{ .key }}-postgres
{{- end -}}

{{/*
Имя Deployment/Service для Redis -> "<release>-redis"
*/}}
{{- define "auto-platform.redisName" -}}
{{ .root.Release.Name }}-redis
{{- end -}}

{{/*
Полный путь к образу приложения с учётом global.imageRegistry.
Если imageRegistry пуст (локальный образ без registry, например для minikube) —
слэш-префикс не добавляется, иначе получилась бы невалидная ссылка "/repo:tag".
*/}}
{{- define "auto-platform.image" -}}
{{- if .root.Values.global.imageRegistry -}}
{{ .root.Values.global.imageRegistry }}/{{ .svc.image.repository }}:{{ .svc.image.tag }}
{{- else -}}
{{ .svc.image.repository }}:{{ .svc.image.tag }}
{{- end -}}
{{- end -}}

{{/*
Общие labels для ресурсов конкретного сервиса
*/}}
{{- define "auto-platform.labels" -}}
app.kubernetes.io/name: {{ .key }}
app.kubernetes.io/instance: {{ .root.Release.Name }}
app.kubernetes.io/part-of: auto-platform
{{- end -}}
