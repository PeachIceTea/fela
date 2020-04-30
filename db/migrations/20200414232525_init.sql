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

INSERT INTO user (name, password, role)
VALUES ('unknown', '', 'user')

INSERT INTO user (name, password, role)
VALUES ('admin', '$2y$10$ddqHCunFxTAZxznsuk.O7.RrOc3Hx1s.LIbCo8d4.XapAOUIVRU3O', 'admin')

-- migrate:down
DROP VIEW audiobook_with_duration;
DROP TABLE file;
DROP TABLE audiobook;
DROP TABLE user;
