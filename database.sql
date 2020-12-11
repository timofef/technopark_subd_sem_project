CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS forum_users;
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS threads;


DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    nickname CITEXT UNIQUE NOT NULL,
    fullname CITEXT        NOT NULL,
    about    TEXT          NOT NULL,
    email    CITEXT UNIQUE NOT NULL
);



DROP TABLE IF EXISTS forums CASCADE;

CREATE TABLE forums
(
    id      SERIAL PRIMARY KEY,
    title   TEXT          NOT NULL,
    owner   CITEXT        NOT NULL,
    posts   INT DEFAULT 0,
    threads INT DEFAULT 0,
    slug    CITEXT UNIQUE NOT NULL,
    FOREIGN KEY (owner) REFERENCES users (nickname)
);

/*CREATE TABLE threads
(
    author      CITEXT NOT NULL,
    create_date timestamptz DEFAULT now(),
    forum       CITEXT NOT NULL,
    id          SERIAL PRIMARY KEY,
    msg         TEXT   NOT NULL,
    slug        CITEXT UNIQUE,
    title       CITEXT NOT NULL,
    votes       INT         DEFAULT 0,
    FOREIGN KEY (forum) REFERENCES forums (slug),
    FOREIGN KEY (author) REFERENCES users (nickname)
);

CREATE TABLE posts
(
    author      CITEXT NOT NULL,
    create_date timestamptz DEFAULT now(),
    forum       CITEXT NOT NULL,
    id          SERIAL PRIMARY KEY,
    is_edited   BOOLEAN     DEFAULT FALSE,
    msg         TEXT   NOT NULL,
    parent      INT    NOT NULL,
    thread      INT    NOT NULL,
    path        BIGINT ARRAY,
    FOREIGN KEY (forum) REFERENCES forums (slug),
    FOREIGN KEY (author) REFERENCES users (nickname),
    FOREIGN KEY (thread) REFERENCES threads (id)
);

CREATE TABLE votes
(
    thread   INT    NOT NULL,
    voice    INT    NOT NULL,
    nickname CITEXT NOT NULL,
    FOREIGN KEY (thread) REFERENCES threads (id),
    FOREIGN KEY (nickname) REFERENCES users (nickname),
    UNIQUE (thread, nickname)
);

CREATE TABLE forum_users
(
    forum    CITEXT NOT NULL,
    nickname CITEXT NOT NULL,
    FOREIGN KEY (forum) REFERENCES forums (slug),
    FOREIGN KEY (nickname) REFERENCES users (nickname),
    UNIQUE (forum, nickname)
);


TRUNCATE TABLE posts;
TRUNCATE TABLE votes;
TRUNCATE TABLE forum_users;
TRUNCATE TABLE threads CASCADE;*/
TRUNCATE TABLE forums CASCADE;
TRUNCATE TABLE users CASCADE;