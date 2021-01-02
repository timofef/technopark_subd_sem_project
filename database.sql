CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS forum_users CASCADE;
DROP TABLE IF EXISTS votes CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS users CASCADE;

DROP FUNCTION IF EXISTS insert_thread_votes();
DROP FUNCTION IF EXISTS update_thread_votes();
DROP FUNCTION IF EXISTS set_post_path();
DROP FUNCTION IF EXISTS update_forum_threads();
DROP FUNCTION IF EXISTS update_forum_posts();
DROP FUNCTION IF EXISTS add_forum_user();

DROP TRIGGER IF EXISTS insert_thread_votes ON votes;
DROP TRIGGER IF EXISTS update_thread_votes ON votes;
DROP TRIGGER IF EXISTS set_post_path ON posts;
DROP TRIGGER IF EXISTS update_forum_threads ON threads;
DROP TRIGGER IF EXISTS update_forum_posts ON posts;
DROP TRIGGER IF EXISTS add_forum_user_new_thread ON threads;
DROP TRIGGER IF EXISTS add_forum_user_new_post ON posts;

/*                                     TABLES                                            */
/*---------------------------------------------------------------------------------------*/

CREATE UNLOGGED TABLE users
(
    id       SERIAL PRIMARY KEY,
    nickname CITEXT COLLATE "C" UNIQUE NOT NULL,
    fullname CITEXT COLLATE "C"        NOT NULL,
    about    TEXT                      NOT NULL,
    email    CITEXT UNIQUE             NOT NULL
);

CREATE INDEX index_users_nickname_hash ON users using hash (nickname);
CREATE INDEX index_users_email_hash ON users using hash (email);

/*---------------------------------------------------------------------------------------*/

CREATE UNLOGGED TABLE forums
(
    id      SERIAL PRIMARY KEY,
    title   TEXT                      NOT NULL,
    owner   CITEXT COLLATE "C"        NOT NULL,
    posts   INT DEFAULT 0,
    threads INT DEFAULT 0,
    slug    CITEXT COLLATE "C" UNIQUE NOT NULL,
    FOREIGN KEY (owner) REFERENCES users (nickname)
);

CREATE INDEX index_forums ON forums (slug, title, owner, posts, threads);
CREATE INDEX index_forums_slug_hash ON forums USING hash (slug);

CREATE INDEX index_forums_users_foreign on forums (owner);

/*---------------------------------------------------------------------------------------*/

CREATE UNLOGGED TABLE threads
(
    id      SERIAL PRIMARY KEY,
    author  CITEXT COLLATE "C" NOT NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT now(),
    forum   CITEXT COLLATE "C" NOT NULL,
    message TEXT               NOT NULL,
    slug    CITEXT COLLATE "C" UNIQUE,
    title   CITEXT COLLATE "C" NOT NULL,
    votes   INT                      DEFAULT 0,
    FOREIGN KEY (forum) REFERENCES forums (slug) ON DELETE CASCADE,
    FOREIGN KEY (author) REFERENCES users (nickname) ON DELETE CASCADE
);

create index index_threads_forum_created on threads (forum, created);
create index index_threads_created on threads (created);

create index index_threads_slug_hash on threads using hash (slug);
create index index_threads_id_hash on threads using hash (id);

/*---------------------------------------------------------------------------------------*/

CREATE UNLOGGED TABLE posts
(
    id        BIGSERIAL PRIMARY KEY,
    author    CITEXT COLLATE "C" NOT NULL,
    created   TIMESTAMP WITH TIME ZONE DEFAULT now(),
    forum     CITEXT COLLATE "C" NOT NULL,
    is_edited BOOLEAN                  DEFAULT FALSE,
    message   TEXT               NOT NULL,
    parent    INT                NOT NULL,
    thread    INT                NOT NULL,
    path      BIGINT[],
    FOREIGN KEY (forum) REFERENCES forums (slug) ON DELETE CASCADE,
    FOREIGN KEY (author) REFERENCES users (nickname) ON DELETE CASCADE,
    FOREIGN KEY (thread) REFERENCES threads (id) ON DELETE CASCADE
);

create index index_posts_id on posts (id);
create index index_posts_thread_created_id on posts (thread, created, id);
create index index_posts_thread_id on posts (thread, id);
create index index_posts_thread_path on posts (thread, path);
create index index_posts_thread_parent_path on posts (thread, parent, path);
create index index_posts_path1_path on posts ((path[1]), path);

/*---------------------------------------------------------------------------------------*/

CREATE UNLOGGED TABLE votes
(
    thread   INT                NOT NULL,
    voice    INT                NOT NULL,
    nickname CITEXT COLLATE "C" NOT NULL,
    FOREIGN KEY (thread) REFERENCES threads (id),
    FOREIGN KEY (nickname) REFERENCES users (nickname),
    UNIQUE (thread, nickname)
);

create unique index index_votes_user_thread on votes (thread, nickname);

/*---------------------------------------------------------------------------------------*/

CREATE UNLOGGED TABLE forum_users
(
    forum    CITEXT COLLATE "C" NOT NULL,
    nickname CITEXT COLLATE "C" NOT NULL,
    FOREIGN KEY (forum) REFERENCES forums (slug) ON DELETE CASCADE,
    FOREIGN KEY (nickname) REFERENCES users (nickname) ON DELETE CASCADE,
    UNIQUE (forum, nickname)
);

create index index_forum_users on forum_users (forum, nickname);
create index index_forum_users_nickname on forum_users (nickname);
cluster forum_users using index_forum_users;

/*---------------------------------------------------------------------------------------*/

TRUNCATE TABLE forum_users CASCADE;
TRUNCATE TABLE posts CASCADE;
TRUNCATE TABLE votes CASCADE;
TRUNCATE TABLE threads CASCADE;
TRUNCATE TABLE forums CASCADE;
TRUNCATE TABLE users CASCADE;


/*                                     TRIGGERS                                          */
/*---------------------------------------------------------------------------------------*/

CREATE OR REPLACE FUNCTION insert_thread_votes()
    RETURNS TRIGGER AS
$insert_thread_votes$
BEGIN
    IF new.voice > 0 THEN
        UPDATE threads SET votes = (votes + 1)
        WHERE id = new.thread;
    ELSE
        UPDATE threads SET votes = (votes - 1)
        WHERE id = new.thread;
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
        UPDATE threads
        SET votes = (votes + 2)
        WHERE threads.id = new.thread;
    else
        UPDATE threads
        SET votes = (votes - 2)
        WHERE threads.id = new.thread;
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
    INSERT INTO forum_users (nickname, forum)
    VALUES (new.author, new.forum)
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