CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS forum_users;
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS posts;

/*---------------------------------------------------------------------------------------*/

DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    nickname CITEXT UNIQUE NOT NULL,
    fullname CITEXT        NOT NULL,
    about    TEXT          NOT NULL,
    email    CITEXT UNIQUE NOT NULL
);

create index index_users on users (nickname, fullname, email, about);

/*---------------------------------------------------------------------------------------*/

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


/*---------------------------------------------------------------------------------------*/

DROP TABLE IF EXISTS threads CASCADE;

CREATE TABLE threads
(
    id      SERIAL PRIMARY KEY,
    author  CITEXT NOT NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT now(),
    forum   CITEXT NOT NULL,
    message TEXT   NOT NULL,
    slug    CITEXT UNIQUE,
    title   CITEXT NOT NULL,
    votes   INT                      DEFAULT 0,
    FOREIGN KEY (forum) REFERENCES forums (slug) ON DELETE CASCADE,
    FOREIGN KEY (author) REFERENCES users (nickname) ON DELETE CASCADE
);

/*---------------------------------------------------------------------------------------*/

CREATE TABLE posts
(
    id        SERIAL PRIMARY KEY,
    author    CITEXT NOT NULL,
    created   TIMESTAMP WITH TIME ZONE DEFAULT now(),
    forum     CITEXT NOT NULL,
    is_edited BOOLEAN                  DEFAULT FALSE,
    message   TEXT   NOT NULL,
    parent    INT    NOT NULL,
    thread    INT    NOT NULL,
    path      BIGINT[],
    FOREIGN KEY (forum) REFERENCES forums (slug) ON DELETE CASCADE,
    FOREIGN KEY (author) REFERENCES users (nickname) ON DELETE CASCADE,
    FOREIGN KEY (thread) REFERENCES threads (id) ON DELETE CASCADE
);

/*---------------------------------------------------------------------------------------*/

CREATE TABLE votes
(
    thread   INT      NOT NULL,
    voice    INT NOT NULL,
    nickname CITEXT   NOT NULL,
    FOREIGN KEY (thread) REFERENCES threads (id),
    FOREIGN KEY (nickname) REFERENCES users (nickname),
    UNIQUE (thread, nickname)
);


/*
CREATE TABLE forum_users
(
    forum    CITEXT NOT NULL,
    nickname CITEXT NOT NULL,
    FOREIGN KEY (forum) REFERENCES forums (slug),
    FOREIGN KEY (nickname) REFERENCES users (nickname),
    UNIQUE (forum, nickname)
);



TRUNCATE TABLE forum_users;*/
TRUNCATE TABLE posts;
TRUNCATE TABLE votes;
TRUNCATE TABLE threads CASCADE;
TRUNCATE TABLE forums CASCADE;
TRUNCATE TABLE users CASCADE;


CREATE OR REPLACE FUNCTION insert_thread_votes()
    RETURNS TRIGGER AS
$insert_thread_votes$
DECLARE
BEGIN
    IF new.voice > 0 THEN
        UPDATE threads SET votes = (votes + 1) WHERE id = new.thread;
    ELSE
        UPDATE threads SET votes = (votes - 1) WHERE id = new.thread;
    END IF;
    RETURN new;
END;
$insert_thread_votes$ language plpgsql;

create trigger insert_thread_votes
    before insert
    on votes
    for each row
execute procedure insert_thread_votes();


create or replace function update_thread_votes()
    returns trigger as
$update_thread_votes$
begin
    if new.voice > 0 then
        update threads set votes = (votes + 2) where threads.id = new.thread;
    else
        update threads set votes = (votes - 2) where threads.id = new.thread;
    end if;
    return new;
end;
$update_thread_votes$ language plpgsql;

create trigger update_thread_votes
    before update
    on votes
    for each row
execute procedure update_thread_votes();