DELETE a, f
FROM audiobook a
JOIN file f ON a.id = f.audiobook
WHERE a.id = ?
