INSERT INTO file (name, hash, codec, duration, metadata)
VALUES (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE id = id
