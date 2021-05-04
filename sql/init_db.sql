GRANT ALL PRIVILEGES ON database formdb TO postgres;

create table if not exists users
(
    id          serial primary key,
    nickname    varchar(128) UNIQUE,
    fullname    varchar(128) NOT NULL,
    email       varchar(128) UNIQUE,
    about       TEXT
);


create table if not exists forum
(
    id          serial primary key,
    title       varchar(128) NOT NULL,
    nickname    varchar(128) NOT NULL,
    slug        varchar(128) NOT NULL UNIQUE,
    post        bigint default 0,
    threads     int default 0,

    FOREIGN KEY (nickname) REFERENCES Users (nickname) 
);

create table if not exists thread
(
    id          serial primary key,
    title       varchar(128) NOT NULL,
    author      varchar(128) NOT NULL,
    forum       varchar(128) NOT NULL,
    message     TEXT NOT NULL,
    votes       int default 0,
    created     timestamp with time zone default now(),
    slug        varchar(128),

    FOREIGN KEY (author) REFERENCES Users (nickname),
    FOREIGN KEY (forum) REFERENCES forum (slug) 
);

create table if not exists  posts
(
    id          serial primary key,
    parent      int default 0,
    author      varchar(128) NOT NULL,
    message     TEXT NOT NULL,
    is_edited   boolean default false,
    forum       varchar(128) NOT NULL,
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
    nickname    varchar(128) NOT NULL,
    thread      int  NOT NULL,
    voice       int  not null,

    FOREIGN KEY (nickname) REFERENCES Users (nickname),
    FOREIGN KEY (thread) REFERENCES thread(id),
    unique (nickname, thread)
);

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

CREATE UNLOGGED TABLE users_forum
(
    id          serial primary key,
    nickname    varchar(128) NOT NULL,
    slug        varchar(128) NOT NULL,

    FOREIGN KEY (nickname) REFERENCES Users (nickname),
    FOREIGN KEY (slug) REFERENCES forum(slug),
    UNIQUE (nickname, slug)
);

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