CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS forum_users;
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS posts;

/*---------------------------------------------------------------------------------------*/

DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    nickname CITEXT COLLATE "C" UNIQUE NOT NULL,
    fullname CITEXT COLLATE "C"        NOT NULL,
    about    TEXT                      NOT NULL,
    email    CITEXT UNIQUE             NOT NULL
);

create index index_users on users (nickname, fullname, email, about);

/*---------------------------------------------------------------------------------------*/

DROP TABLE IF EXISTS forums CASCADE;

CREATE TABLE forums
(
    id      SERIAL PRIMARY KEY,
    title   TEXT          NOT NULL,
    owner   CITEXT   COLLATE "C"     NOT NULL,
    posts   INT DEFAULT 0,
    threads INT DEFAULT 0,
    slug    CITEXT COLLATE "C" UNIQUE NOT NULL,
    FOREIGN KEY (owner) REFERENCES users (nickname)
);


/*---------------------------------------------------------------------------------------*/

DROP TABLE IF EXISTS threads CASCADE;

CREATE TABLE threads
(
    id      SERIAL PRIMARY KEY,
    author  CITEXT COLLATE "C" NOT NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT now(),
    forum   CITEXT COLLATE "C" NOT NULL,
    message TEXT   NOT NULL,
    slug    CITEXT COLLATE "C" UNIQUE,
    title   CITEXT COLLATE "C" NOT NULL,
    votes   INT                      DEFAULT 0,
    FOREIGN KEY (forum) REFERENCES forums (slug) ON DELETE CASCADE,
    FOREIGN KEY (author) REFERENCES users (nickname) ON DELETE CASCADE
);

/*---------------------------------------------------------------------------------------*/

CREATE TABLE posts
(
    id        BIGSERIAL PRIMARY KEY,
    author    CITEXT COLLATE "C" NOT NULL,
    created   TIMESTAMP WITH TIME ZONE DEFAULT now(),
    forum     CITEXT COLLATE "C" NOT NULL,
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
    thread   INT    NOT NULL,
    voice    INT    NOT NULL,
    nickname CITEXT COLLATE "C" NOT NULL,
    FOREIGN KEY (thread) REFERENCES threads (id),
    FOREIGN KEY (nickname) REFERENCES users (nickname),
    UNIQUE (thread, nickname)
);



CREATE TABLE forum_users
(
    forum    CITEXT COLLATE "C" NOT NULL,
    nickname CITEXT COLLATE "C" NOT NULL,
    FOREIGN KEY (forum) REFERENCES forums (slug) ON DELETE CASCADE,
    FOREIGN KEY (nickname) REFERENCES users (nickname) ON DELETE CASCADE,
    UNIQUE (forum, nickname)
);



TRUNCATE TABLE forum_users;
TRUNCATE TABLE posts;
TRUNCATE TABLE votes;
TRUNCATE TABLE threads CASCADE;
TRUNCATE TABLE forums CASCADE;
TRUNCATE TABLE users CASCADE;



CREATE OR REPLACE FUNCTION insert_thread_votes()
    RETURNS TRIGGER AS
$insert_thread_votes$
BEGIN
    IF new.voice > 0 THEN
        UPDATE threads SET votes = (votes + 1) WHERE id = new.thread;
    ELSE
        UPDATE threads SET votes = (votes - 1) WHERE id = new.thread;
    END IF;
    RETURN new;
END;
$insert_thread_votes$ language plpgsql;

CREATE TRIGGER insert_thread_votes
    BEFORE INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE insert_thread_votes();



CREATE OR REPLACE FUNCTION update_thread_votes()
    RETURNS TRIGGER AS
$update_thread_votes$
BEGIN
    IF new.voice > 0 THEN
        UPDATE threads SET votes = (votes + 2) WHERE threads.id = new.thread;
    else
        UPDATE threads SET votes = (votes - 2) WHERE threads.id = new.thread;
    END IF;
    RETURN new;
END;
$update_thread_votes$ LANGUAGE plpgsql;

CREATE TRIGGER update_thread_votes
    BEFORE UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE update_thread_votes();



CREATE OR REPLACE FUNCTION set_post_path()
    RETURNS TRIGGER AS
$set_post_path$
DECLARE
    parent_thread BIGINT;
    parent_path   BIGINT[];
BEGIN
    IF (new.parent = 0) THEN
        new.path := new.path || new.id;
    ELSE
        SELECT thread, path
        FROM posts p
        WHERE p.thread = new.thread
          AND p.id = new.parent
        INTO parent_thread , parent_path;
        IF parent_thread != new.thread OR NOT FOUND THEN
            RAISE EXCEPTION USING ERRCODE = '00404';
        END IF;
        new.path := parent_path || new.id;
    END IF;
    RETURN new;
END;
$set_post_path$ LANGUAGE plpgsql;

CREATE TRIGGER set_post_path
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE set_post_path();



CREATE OR REPLACE FUNCTION update_forum_threads()
    RETURNS TRIGGER AS
$update_forum_threads$
BEGIN
    UPDATE forums SET threads = threads + 1 WHERE slug = new.forum;
    RETURN new;
END;
$update_forum_threads$ LANGUAGE plpgsql;

CREATE TRIGGER update_forum_threads
    BEFORE INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_threads();



CREATE OR REPLACE FUNCTION update_forum_posts()
    RETURNS TRIGGER AS
$update_forum_posts$
BEGIN
    UPDATE forums SET posts = posts + 1 WHERE slug = new.forum;
    RETURN new;
END;
$update_forum_posts$ LANGUAGE plpgsql;

CREATE TRIGGER update_forum_posts
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_posts();


CREATE OR REPLACE FUNCTION add_forum_user()
    RETURNS TRIGGER AS
$add_forum_user$
BEGIN
    INSERT INTO forum_users (nickname, forum) VALUES (new.author, new.forum)
    ON CONFLICT DO NOTHING;
    RETURN new;
END;
$add_forum_user$ LANGUAGE plpgsql;

CREATE TRIGGER add_forum_user_new_thread
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE add_forum_user();

CREATE TRIGGER add_forum_user_new_post
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE add_forum_user();