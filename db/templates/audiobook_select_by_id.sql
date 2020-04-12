SELECT a.*, SUM(f.duration) duration
FROM audiobook a
JOIN file f ON f.audiobook = a.id
WHERE a.id = ?
GROUP BY a.id
