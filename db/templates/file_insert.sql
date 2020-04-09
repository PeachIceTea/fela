INSERT INTO file (name, hash, codec)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE id = id
