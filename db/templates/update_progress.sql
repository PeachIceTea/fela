INSERT into progress (user, audiobook, file, progress)
VALUES (:user, :audiobook, :file, :progress)
ON DUPLICATE KEY UPDATE 
    file = :file,
    progress = :progress