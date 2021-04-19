
GRANT ALL PRIVILEGES ON database formDB TO postgres;

create table if not exists users
(
    id          serial primary key,
    nickname    varchar(128) UNIQUE,
    fullname    varchar(128) NOT NULL,
    email       varchar(128) UNIQUE,
    about       TEXT
);
