CREATE IF NOT EXISTS "role"(
       id BIGSERIAL PRIMARY KEY,
       name VARCHAR(255) NOT NULL UNIQUE,
       level INT NOT NULL DEFAULT 0,
       description TEXT
);

INSERT INTO
    "role"(name, level, description)
    VALUES ('user', 'A user can create post and comments', 1),
           ('moderator', 'A moderator can update other users posts', 2),
           ('admin', 'An admin can update and delete other users posts', 3);