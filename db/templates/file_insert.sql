INSERT INTO file (name, hash, metadata, kind)
VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE id = id
