DROP SCHEMA zbx1 cascade;

CREATE SCHEMA zbx1;

CREATE TABLE zbx1.users (
    id serial not null primary key,
    u_name varchar(50) not null,
    u_email varchar(50) not null,
    u_username varchar(50) not null,
    u_password varchar(255) not null,
    u_created timestamp not null,
    UNIQUE(u_name)
);

CREATE TABLE zbx1.files (
    id serial not null primary key,
    f_name varchar(50) not null,
    f_ext varchar(5) not null,
    f_created timestamp not null,
    f_data BYTEA not null
);

CREATE TABLE zbx1.privileges (
    id serial not null primary key,
    p_name varchar(50) not null,
    p_created timestamp not null,
    UNIQUE (p_name)
);

CREATE TABLE zbx1.users_files (
  user_id INTEGER REFERENCES zbx1.users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  file_id INTEGER REFERENCES zbx1.files(id) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT user_file_pkey PRIMARY KEY (user_id, file_id)
);

-- CREATE DEFAULT USER
INSERT INTO zbx1.users(u_name, u_email, u_username, u_password, u_created) VALUES ('Tiago','metiago@gmail.com','metiago','$2a$14$LAYYntKvxeR1TBLjTfCXpOhqQwfZhMfW4AcJmq1Rx/IXSMfEbJy1K', now());

-- CREATE DEFAULT PERMISSIONS
INSERT INTO zbx1.privileges(p_name, p_created) VALUES ('READ', now());
INSERT INTO zbx1.privileges(p_name, p_created) VALUES ('WRITE', now());
INSERT INTO zbx1.privileges(p_name, p_created) VALUES ('DELETE', now());