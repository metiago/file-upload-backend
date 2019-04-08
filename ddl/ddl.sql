DROP DATABASE zbx1;

CREATE DATABASE zbx1;

USE zbx1;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
    id int(11) unsigned not null auto_increment,
    u_name varchar(50) not null,
    u_email varchar(50) not null,
    u_username varchar(50) not null,
    u_password varchar(255) not null,
    u_created datetime not null,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `UNI_U_NAME` (`u_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles` (
    id int(11) unsigned not null auto_increment,
    r_name varchar(50) not null,
    r_created datetime not null,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `UNI_R_NAME` (`r_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `files`;
CREATE TABLE `files` (
    id int(11) unsigned not null auto_increment,
    f_name varchar(50) not null,
    f_ext varchar(5) not null,
    f_created datetime not null,
    f_data LONGBLOB not null,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `UNI_F_NAME` (`f_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `privileges`;
CREATE TABLE `privileges` (
    id int(11) unsigned not null auto_increment,
    p_name varchar(50) not null,
    p_created datetime not null,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `UNI_P_NAME` (`p_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `users_roles`;
CREATE TABLE `users_roles` (
  `user_id` int(10) unsigned NOT NULL,
  `role_id` int(10) unsigned NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  PRIMARY KEY (`user_id`, `role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `users_files`;
CREATE TABLE `users_files` (
  `user_id` int(10) unsigned NOT NULL,
  `file_id` int(10) unsigned NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  PRIMARY KEY (`user_id`, `file_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `roles_privileges`;
CREATE TABLE `roles_privileges` (
  `privilege_id` int(10) unsigned NOT NULL,
  `role_id` int(10) unsigned NOT NULL,
  FOREIGN KEY (`privilege_id`) REFERENCES `privileges` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  PRIMARY KEY (`privilege_id`, `role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- CREATE DEFAULT USER
INSERT INTO users(u_name, u_email, u_username, u_password, u_created) VALUES ('Tiago','metiago@gmail.com','metiago','$2a$14$LAYYntKvxeR1TBLjTfCXpOhqQwfZhMfW4AcJmq1Rx/IXSMfEbJy1K', now());

-- CREATE DEFAULT ROLE
INSERT INTO roles(r_name, r_created) VALUES ('ADMIN', now());

-- CREATE DEFAULT PERMISSIONS
INSERT INTO privileges(p_name, p_created) VALUES ('READ', now());
INSERT INTO privileges(p_name, p_created) VALUES ('WRITE', now());
INSERT INTO privileges(p_name, p_created) VALUES ('DELETE', now());

-- LINK DEFAULT USERS AND ROLES
INSERT INTO users_roles(user_id, role_id) VALUES ((SELECT id FROM `users` WHERE u_username='metiago'), (SELECT id FROM `roles` WHERE r_name='ADMIN'));

-- LINK DEFAULT PRIVILEGES AND ROLES
INSERT INTO roles_privileges(privilege_id, role_id) VALUES ((SELECT id FROM `privileges` WHERE p_name='READ'), (SELECT id FROM `roles` WHERE r_name='ADMIN'));
INSERT INTO roles_privileges(privilege_id, role_id) VALUES ((SELECT id FROM `privileges` WHERE p_name='WRITE'), (SELECT id FROM `roles` WHERE r_name='ADMIN'));
INSERT INTO roles_privileges(privilege_id, role_id) VALUES ((SELECT id FROM `privileges` WHERE p_name='DELETE'), (SELECT id FROM `roles` WHERE r_name='ADMIN'));