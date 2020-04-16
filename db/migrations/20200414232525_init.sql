-- migrate:up
CREATE TABLE user (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,

	name VARCHAR(255) NOT NULL UNIQUE,
	password BINARY(60) NOT NULL,

	role ENUM('user', 'uploader', 'admin') NOT NULL DEFAULT 'user',

	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE audiobook (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,

	title VARCHAR(255),
	author VARCHAR(255),

	uploader BIGINT NOT NULL,

	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME ON UPDATE CURRENT_TIMESTAMP,

	FOREIGN KEY (uploader) REFERENCES user(id)
);

CREATE TABLE file (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,

	name VARCHAR(255) NOT NULL,
	codec VARCHAR(50) NOT NULL,
	duration DOUBLE NOT NULL,
	metadata JSON NOT NULL,

	path VARCHAR(255) NOT NULL,

	audiobook BIGINT NOT NULL,

	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME ON UPDATE CURRENT_TIMESTAMP,

	FOREIGN KEY (audiobook) REFERENCES audiobook(id)
);

CREATE VIEW audiobook_with_duration AS
SELECT a.*, SUM(f.duration) duration
FROM audiobook a
JOIN file f ON a.id = f.audiobook;


-- migrate:down
DROP VIEW audiobook_with_duration;
DROP TABLE file;
DROP TABLE audiobook;
DROP TABLE user;
