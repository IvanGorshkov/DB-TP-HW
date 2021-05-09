GRANT ALL PRIVILEGES ON database formdb TO postgres;
CREATE EXTENSION IF NOT EXISTS citext;

create table if not exists users
(
    id          serial primary key,
    nickname    citext UNIQUE,
    fullname    varchar(128) NOT NULL,
    email       citext UNIQUE,
    about       TEXT
);


create table if not exists forum
(
    id          serial primary key,
    title       varchar(128) NOT NULL,
    nickname    citext NOT NULL,
    slug        citext NOT NULL UNIQUE,
    post        bigint default 0,
    threads     int default 0,

    FOREIGN KEY (nickname) REFERENCES Users (nickname) 
);

create table if not exists thread
(
    id          serial primary key,
    title       varchar(128) NOT NULL,
    author      citext NOT NULL,
    forum       citext NOT NULL,
    message     TEXT NOT NULL,
    votes       int default 0,
    created     timestamp with time zone default now(),
    slug        citext,

    FOREIGN KEY (author) REFERENCES Users (nickname),
    FOREIGN KEY (forum) REFERENCES forum (slug) 
);

create table if not exists  posts
(
    id          serial primary key,
    parent      int default 0,
    author      citext NOT NULL,
    message     TEXT NOT NULL,
    is_edited   boolean default false,
    forum       citext NOT NULL,
    thread      integer  NOT NULL,
    created     timestamp with time zone default now(),
    Path        INTEGER[]  DEFAULT ARRAY []::INTEGER[],

    FOREIGN KEY (author) REFERENCES Users (nickname),
    FOREIGN KEY (forum) REFERENCES forum (slug),
    FOREIGN KEY (thread) REFERENCES thread (id) 
);

create table if not exists  votes
(
    id          serial primary key,
    nickname    citext NOT NULL,
    thread      int  NOT NULL,
    voice       int  not null,

    FOREIGN KEY (nickname) REFERENCES Users (nickname),
    FOREIGN KEY (thread) REFERENCES thread(id),
    unique (nickname, thread)
);

create table if not exists users_forum
(
    id          serial primary key,
    nickname    citext NOT NULL,
    slug        citext NOT NULL,

    FOREIGN KEY (nickname) REFERENCES Users (nickname),
    FOREIGN KEY (slug) REFERENCES forum(slug),
    UNIQUE (nickname, slug)
);


CREATE OR REPLACE FUNCTION thread_insert() RETURNS TRIGGER AS
$thread_insert$
DECLARE
    Slug_thread         TEXT;
BEGIN
    SELECT slug from thread where slug = NEW.slug INTO Slug_thread;
    IF (Slug_thread <> '') THEN
        RAISE EXCEPTION 'Cant find thread by slug' USING ERRCODE = '23505';
    end if;

    return NEW;
end
$thread_insert$ LANGUAGE plpgsql;

CREATE TRIGGER thread_insert
    BEFORE INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE thread_insert();


CREATE OR REPLACE FUNCTION insert_votes_threads()
    RETURNS TRIGGER AS
$insert_votes_threads$
BEGIN
    UPDATE thread
    SET votes = votes + NEW.voice
    WHERE id = NEW.thread;
    RETURN NEW;
END;
$insert_votes_threads$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS insert_votes_threads ON votes;
CREATE TRIGGER insert_votes_threads
    AFTER INSERT 
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE insert_votes_threads();

CREATE OR REPLACE FUNCTION update_votes_threads()
    RETURNS TRIGGER AS
$update_votes_threads$
BEGIN
    UPDATE thread
    SET votes = votes + (2 * NEW.voice)
    WHERE id = NEW.thread;
    RETURN NEW;
END;
$update_votes_threads$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS update_votes_threads ON votes;
CREATE TRIGGER update_votes_threads
    AFTER UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE update_votes_threads();

CREATE OR REPLACE FUNCTION updatePath() RETURNS TRIGGER AS
$update_path$
DECLARE
    parentPath         INTEGER[];
    first_parent_thread INT;
BEGIN
    IF (NEW.parent IS NULL) THEN
        NEW.path := array_append(new.path, new.id);
    ELSE
        SELECT path FROM posts WHERE id = new.parent INTO parentPath;
        SELECT thread FROM posts WHERE id = parentPath[1] INTO first_parent_thread;
        IF NOT FOUND OR first_parent_thread != NEW.thread THEN
            RAISE EXCEPTION 'parent is from different thread' USING ERRCODE = '00409';
        end if;

        NEW.path := NEW.path || parentPath || new.id;
    end if;
    UPDATE forum SET post=post + 1 WHERE forum.slug = new.forum;
    RETURN new;
end
$update_path$ LANGUAGE plpgsql;

CREATE TRIGGER update_path_trigger
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE updatePath();



CREATE OR REPLACE FUNCTION updateCountOfThreads() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    UPDATE forum SET Threads=(Threads+1) WHERE slug=NEW.forum;
    return NEW;
end
$update_users_forum$ LANGUAGE plpgsql;

CREATE TRIGGER addThreadInForum
    BEFORE INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE updateCountOfThreads();

CREATE OR REPLACE FUNCTION update_user_forum() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    INSERT INTO users_forum (nickname, slug)
    VALUES (NEW.author, NEW.forum) on conflict do nothing;
    return NEW;
end
$update_users_forum$ LANGUAGE plpgsql;

CREATE TRIGGER thread_insert_user_forum
    AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE update_user_forum();

CREATE TRIGGER post_insert_user_forum
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_user_forum();


CREATE INDEX if not exists user_nickname_index  ON users using hash (nickname);
CREATE INDEX if not exists user_email_index     ON users using hash (email);

CREATE INDEX if not exists forum_slug_index     ON forum using hash (slug);

create unique index if not exists forum_users_unique_index on users_forum (slug, nickname);
cluster users_forum using forum_users_unique_index;

CREATE INDEX if not exists thread_slug_index        ON thread using hash (slug);
CREATE INDEX if not exists thread_forum_index       ON thread using hash (forum);
CREATE INDEX if not exists thread_date_index        ON thread (created);
CREATE INDEX if not exists thread_forum_date_index  ON thread (forum, created);

create index if not exists post_id_path_index  on posts (id, (path[1]));
create index if not exists post_thread_id_path1_parent_index on posts (thread, id, (path[1]), parent);
create index if not exists post_thread_path_id_index  on posts (thread, path, id);


create index if not exists post_path1_index on posts ((path[1]));
create index if not exists post_thread_id_index on posts (thread, id);
CREATE INDEX if not exists post_thread_index ON posts (thread);

create unique index if not exists vote_unique_index on votes (nickname, Thread);
