CREATE DATABASE IF NOT EXISTS devBook;
USE devbook;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS followers;

CREATE TABLE users(
    id int auto_increment primary key,
    name varchar(255) not null,
    nick varchar(255) not null unique,
    email varchar(255) not null unique,
    password varchar(255) not null,
    created_at timestamp default current_timestamp()
)ENGINE=INNODB;

CREATE TABLE followers(
    user_id int not null,
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

    follower_id int not null,
    FOREIGN KEY (follower_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

    PRIMARY KEY(user_id, follower_id)
)ENGINE=INNODB;