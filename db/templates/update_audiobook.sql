UPDATE audiobook
SET
{{- if .Title }} title = :title {{- if .Author -}} , {{- end -}} {{- end -}}
{{- if .Author }} author = :author {{ end -}}
WHERE id = :id
