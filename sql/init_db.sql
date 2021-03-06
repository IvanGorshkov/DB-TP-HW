GRANT ALL PRIVILEGES ON DATABASE FORMDB TO DOCKER;

--GRANT ALL PRIVILEGES ON DATABASE FORMDB TO POSTGRES;

CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE UNLOGGED TABLE IF NOT EXISTS USERS (
    ID          SERIAL       PRIMARY KEY,
    NICKNAME    CITEXT       UNIQUE,
    FULLNAME    VARCHAR(128) NOT NULL,
    EMAIL       CITEXT       UNIQUE,
    ABOUT       TEXT
);

CREATE UNLOGGED TABLE IF NOT EXISTS FORUM (
    ID          SERIAL       PRIMARY KEY,
    TITLE       VARCHAR(128) NOT NULL,
    NICKNAME    CITEXT       NOT NULL,
    SLUG        CITEXT       NOT NULL UNIQUE,
    POST        INT          DEFAULT 0,
    THREADS     INT          DEFAULT 0,

    FOREIGN KEY (NICKNAME) REFERENCES USERS (NICKNAME)
);

CREATE UNLOGGED TABLE IF NOT EXISTS THREAD (
    ID          SERIAL       PRIMARY KEY,
    TITLE       VARCHAR(128) NOT NULL,
    AUTHOR      CITEXT       NOT NULL,
    FORUM       CITEXT       NOT NULL,
    MESSAGE     TEXT         NOT NULL,
    VOTES       INT          DEFAULT 0,
    CREATED     TIMESTAMP    WITH TIME ZONE DEFAULT NOW(),
    SLUG        CITEXT,

    FOREIGN KEY (AUTHOR) REFERENCES USERS (NICKNAME),
    FOREIGN KEY (FORUM)  REFERENCES FORUM (SLUG)
);

CREATE UNLOGGED TABLE IF NOT EXISTS  POSTS (
    ID          SERIAL       PRIMARY KEY,
    PARENT      INT          DEFAULT 0,
    AUTHOR      CITEXT       NOT NULL,
    MESSAGE     TEXT         NOT NULL,
    IS_EDITED   BOOLEAN      DEFAULT FALSE,
    FORUM       CITEXT       NOT NULL,
    THREAD      INT          NOT NULL,
    CREATED     TIMESTAMP    WITH TIME ZONE DEFAULT NOW(),
    PATH        INTEGER[]    DEFAULT ARRAY []::INTEGER[],

    FOREIGN KEY (AUTHOR) REFERENCES USERS  (NICKNAME),
    FOREIGN KEY (FORUM)  REFERENCES FORUM  (SLUG),
    FOREIGN KEY (THREAD) REFERENCES THREAD (ID)
);

CREATE UNLOGGED TABLE IF NOT EXISTS  VOTES (
    ID          SERIAL       PRIMARY KEY,
    NICKNAME    CITEXT       NOT NULL,
    THREAD      INT          NOT NULL,
    VOICE       INT          NOT NULL,

    FOREIGN KEY (NICKNAME) REFERENCES USERS  (NICKNAME),
    FOREIGN KEY (THREAD)   REFERENCES THREAD (ID),
    UNIQUE (NICKNAME, THREAD)
);

CREATE UNLOGGED TABLE IF NOT EXISTS USERS_FORUM (
    ID          SERIAL       PRIMARY KEY,
    NICKNAME    CITEXT       NOT NULL,
    SLUG        CITEXT       NOT NULL,

    FOREIGN KEY (NICKNAME) REFERENCES USERS (NICKNAME),
    FOREIGN KEY (SLUG)     REFERENCES FORUM (SLUG),
    UNIQUE (NICKNAME, SLUG)
);

CREATE OR REPLACE FUNCTION THREAD_INSERT() RETURNS TRIGGER AS
$THREAD_INSERT$
    DECLARE
        SLUG_THREAD TEXT;
    BEGIN
        SELECT SLUG FROM THREAD WHERE SLUG = NEW.SLUG INTO SLUG_THREAD;
        IF (SLUG_THREAD <> '') THEN
            RAISE EXCEPTION 'CANT FIND THREAD BY SLUG' USING ERRCODE = '21212';
        END IF;
        RETURN NEW;
    END
$THREAD_INSERT$
LANGUAGE PLPGSQL;

CREATE OR REPLACE FUNCTION INSERT_VOTES_THREADS() RETURNS TRIGGER AS
$INSERT_VOTES_THREADS$
    BEGIN
        UPDATE THREAD SET VOTES = VOTES + NEW.VOICE WHERE ID = NEW.THREAD;
        RETURN NEW;
    END;
$INSERT_VOTES_THREADS$
LANGUAGE PLPGSQL;

CREATE OR REPLACE FUNCTION UPDATE_VOTES_THREADS() RETURNS TRIGGER AS
$UPDATE_VOTES_THREADS$
    BEGIN
        UPDATE THREAD SET VOTES = VOTES + (2 * NEW.VOICE) WHERE ID = NEW.THREAD;
        RETURN NEW;
    END;
$UPDATE_VOTES_THREADS$
LANGUAGE PLPGSQL;

CREATE OR REPLACE FUNCTION UPDATEPATH() RETURNS TRIGGER AS
$UPDATE_PATH$
    DECLARE
        PARENTPATH          INTEGER[];
        FIRST_PARENT_THREAD INT;
    BEGIN
        IF (NEW.PARENT IS NULL) THEN
            NEW.PATH := ARRAY_APPEND(NEW.PATH, NEW.ID);
        ELSE
            SELECT PATH FROM POSTS WHERE ID = NEW.PARENT INTO PARENTPATH;
            SELECT THREAD FROM POSTS WHERE ID = PARENTPATH[1] INTO FIRST_PARENT_THREAD;
            IF NOT FOUND OR FIRST_PARENT_THREAD != NEW.THREAD THEN
                RAISE EXCEPTION 'PARENT IS FROM DIFFERENT THREAD' USING ERRCODE = '00409';
            END IF;
            NEW.PATH := NEW.PATH || PARENTPATH || NEW.ID;
        END IF;
        UPDATE FORUM SET POST=POST + 1 WHERE FORUM.SLUG = NEW.FORUM;
        RETURN NEW;
    END
$UPDATE_PATH$
LANGUAGE PLPGSQL;

CREATE OR REPLACE FUNCTION UPDATE_COUNT_OF_THREADS() RETURNS TRIGGER AS
$UPDATE_USERS_FORUM$
    BEGIN
        UPDATE FORUM SET THREADS=(THREADS+1) WHERE SLUG=NEW.FORUM;
    RETURN NEW;
    END
$UPDATE_USERS_FORUM$ LANGUAGE PLPGSQL;

CREATE OR REPLACE FUNCTION UPDATE_USER_FORUM() RETURNS TRIGGER AS
$UPDATE_USERS_FORUM$
    BEGIN
        INSERT INTO USERS_FORUM (NICKNAME, SLUG) VALUES (NEW.AUTHOR, NEW.FORUM) ON CONFLICT DO NOTHING;
        RETURN NEW;
    END
$UPDATE_USERS_FORUM$
LANGUAGE PLPGSQL;

CREATE TRIGGER THREAD_INSERT
    BEFORE INSERT
    ON THREAD
    FOR EACH ROW
    EXECUTE PROCEDURE THREAD_INSERT();

CREATE TRIGGER INSERT_VOTES_THREADS
    AFTER INSERT
    ON VOTES
    FOR EACH ROW
    EXECUTE PROCEDURE INSERT_VOTES_THREADS();

CREATE TRIGGER UPDATE_VOTES_THREADS
    AFTER UPDATE
    ON VOTES
    FOR EACH ROW
    EXECUTE PROCEDURE UPDATE_VOTES_THREADS();

CREATE TRIGGER UPDATE_PATH_TRIGGER
    BEFORE INSERT
    ON POSTS
    FOR EACH ROW
    EXECUTE PROCEDURE UPDATEPATH();

CREATE TRIGGER ADD_THREAD_IN_FORUM
    BEFORE INSERT
    ON THREAD
    FOR EACH ROW
    EXECUTE PROCEDURE UPDATE_COUNT_OF_THREADS();

CREATE TRIGGER THREAD_INSERT_USER_FORUM
    AFTER INSERT
    ON THREAD
    FOR EACH ROW
    EXECUTE PROCEDURE UPDATE_USER_FORUM();

CREATE TRIGGER POST_INSERT_USER_FORUM
    AFTER INSERT
    ON POSTS
    FOR EACH ROW
    EXECUTE PROCEDURE UPDATE_USER_FORUM();

CREATE INDEX IF NOT EXISTS USER_EMAIL_NICKNAME_INDEX         ON USERS (NICKNAME);
CREATE INDEX IF NOT EXISTS USER_CAST_NICKNAME_INDEX          ON USERS (CAST(LOWER(nickname) AS bytea));

CREATE INDEX IF NOT EXISTS FORUM_SLUG_INDEX                  ON FORUM USING HASH (SLUG);

CREATE UNIQUE INDEX IF NOT EXISTS FORUM_USERS_UNIQUE_INDEX   ON USERS_FORUM (SLUG, NICKNAME);

CREATE INDEX IF NOT EXISTS THREAD_DATE_INDEX                 ON THREAD (CREATED);
CREATE INDEX IF NOT EXISTS THREAD_FORUM_DATE_INDEX           ON THREAD (FORUM, CREATED);
CREATE INDEX IF NOT EXISTS THREAD_SLUG_INDEX                 ON THREAD USING HASH (SLUG);
CREATE INDEX IF NOT EXISTS THREAD_FORUM_INDEX                ON THREAD USING HASH (FORUM);

CREATE INDEX IF NOT EXISTS POST_ID_PATH_INDEX                ON POSTS (ID, (PATH[1]));
CREATE INDEX IF NOT EXISTS POST_THREAD_ID_PATH1_PARENT_INDEX ON POSTS (THREAD, ID, (PATH[1]), PARENT);
CREATE INDEX IF NOT EXISTS POST_THREAD_PATH_ID_INDEX         ON POSTS (THREAD, PATH, ID);

CREATE INDEX IF NOT EXISTS POST_PATH1_INDEX                  ON POSTS ((PATH[1]));
CREATE INDEX IF NOT EXISTS POST_THREAD_ID_INDEX              ON POSTS (THREAD, ID);
CREATE INDEX IF NOT EXISTS POST_THREAD_INDEX                 ON POSTS (THREAD);

CREATE INDEX IF NOT EXISTS VOTE_INDEX                        ON VOTES (VOICE, NICKNAME, THREAD);
