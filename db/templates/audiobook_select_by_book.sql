SELECT a.*, SUM(f.duration) duration, f.codec
FROM audiobook a
JOIN file f ON f.audiobook = a.id
WHERE a.book = ?
GROUP BY a.id
