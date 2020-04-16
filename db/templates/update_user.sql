UPDATE user
SET
{{- if .Name }} name = :name {{- if or .Password .Role -}} , {{- end -}} {{- end -}}
{{- if .Password }} password = :password {{- if .Role -}} , {{- end -}} {{- end -}}
{{- if .Role }} role = :role {{- end }}
WHERE id = :id
