SELECT {{range $i, $s := .}}{{if $i}},{{end}}{{$s}}{{else}}*{{end}}
FROM file
WHERE id = ?
