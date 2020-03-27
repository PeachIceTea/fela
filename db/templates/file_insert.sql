INSERT INTO file (name, metadata, kind)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE id = id
