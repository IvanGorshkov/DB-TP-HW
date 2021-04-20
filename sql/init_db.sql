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
