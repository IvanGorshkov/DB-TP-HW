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
    created   timestamp with time zone default now(),

    FOREIGN KEY (author) REFERENCES Users (nickname),
    FOREIGN KEY (forum) REFERENCES forum (slug),
    FOREIGN KEY (thread) REFERENCES thread (id) 
);