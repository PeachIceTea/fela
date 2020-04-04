INSERT INTO audio_file (name, hash, metadata)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE id = id
