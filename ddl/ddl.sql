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

CREATE TABLE zbx1.roles (
    id serial not null primary key,
    r_name varchar(50) not null,
    r_created timestamp not null,
    UNIQUE(r_name)
);

CREATE TABLE zbx1.files (
    id serial not null primary key,
    f_name varchar(50) not null,
    f_ext varchar(5) not null,
    f_created timestamp not null,
    f_data BYTEA not null,
    UNIQUE (f_name)
);

CREATE TABLE zbx1.privileges (
    id serial not null primary key,
    p_name varchar(50) not null,
    p_created timestamp not null,
    UNIQUE (p_name)
);

CREATE TABLE zbx1.users_roles (
  user_id INTEGER REFERENCES zbx1.users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  role_id INTEGER REFERENCES zbx1.roles(id) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT user_role_pkey PRIMARY KEY (user_id, role_id)
);


CREATE TABLE zbx1.users_files (
  user_id INTEGER REFERENCES zbx1.users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  file_id INTEGER REFERENCES zbx1.files(id) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT user_file_pkey PRIMARY KEY (user_id, file_id)
);

CREATE TABLE zbx1.roles_privileges (
  privilege_id INTEGER REFERENCES zbx1.privileges(id) ON UPDATE CASCADE ON DELETE CASCADE,
  role_id INTEGER REFERENCES zbx1.roles(id) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT role_privilege_pkey PRIMARY KEY (privilege_id, role_id)
);


-- CREATE DEFAULT USER
INSERT INTO zbx1.users(u_name, u_email, u_username, u_password, u_created) VALUES ('Tiago','metiago@gmail.com','metiago','$2a$14$LAYYntKvxeR1TBLjTfCXpOhqQwfZhMfW4AcJmq1Rx/IXSMfEbJy1K', now());

-- CREATE DEFAULT ROLE
INSERT INTO zbx1.roles(r_name, r_created) VALUES ('ADMIN', now());

-- CREATE DEFAULT PERMISSIONS
INSERT INTO zbx1.privileges(p_name, p_created) VALUES ('READ', now());
INSERT INTO zbx1.privileges(p_name, p_created) VALUES ('WRITE', now());
INSERT INTO zbx1.privileges(p_name, p_created) VALUES ('DELETE', now());

-- LINK DEFAULT USERS AND ROLES
INSERT INTO zbx1.users_roles(user_id, role_id) VALUES ((SELECT id FROM zbx1.users WHERE u_username='metiago'), (SELECT id FROM zbx1.roles WHERE r_name='ADMIN'));

-- LINK DEFAULT PRIVILEGES AND ROLES
INSERT INTO zbx1.roles_privileges(privilege_id, role_id) VALUES ((SELECT id FROM zbx1.privileges WHERE p_name='READ'), (SELECT id FROM zbx1.roles WHERE r_name='ADMIN'));
INSERT INTO zbx1.roles_privileges(privilege_id, role_id) VALUES ((SELECT id FROM zbx1.privileges WHERE p_name='WRITE'), (SELECT id FROM zbx1.roles WHERE r_name='ADMIN'));
INSERT INTO zbx1.roles_privileges(privilege_id, role_id) VALUES ((SELECT id FROM zbx1.privileges WHERE p_name='DELETE'), (SELECT id FROM zbx1.roles WHERE r_name='ADMIN'));