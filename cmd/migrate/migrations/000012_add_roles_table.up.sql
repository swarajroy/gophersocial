CREATE TABLE IF NOT EXISTS roles(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level int NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO
    roles(name, description, level)
VALUES
    ('user', 'a user can create posts and comments', 1);

INSERT INTO
    roles(name, description, level)
VALUES
    (
        'moderator',
        'a moderator can update other user posts and comments',
        2
    );

INSERT INTO
    roles(name, description, level)
VALUES
    (
        'admin',
        'a admin can update and delete other user posts and comments',
        3
    );