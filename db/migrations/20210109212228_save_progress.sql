-- migrate:up
CREATE TABLE progress (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    user BIGINT NOT NULL,
    audiobook BIGINT NOT NULL,
    file BIGINT NOT NULL,

    progress DOUBLE NOT NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE(user, audiobook),

    FOREIGN KEY (user) REFERENCES user(id),
    FOREIGN KEY (audiobook) REFERENCES audiobook(id),
    FOREIGN KEY (file) REFERENCES file(id)
);

-- migrate:down
DROP TABLE progress;