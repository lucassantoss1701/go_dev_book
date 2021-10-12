CREATE DATABASE IF NOT EXISTS devBook;
USE devbook;

DROP TABLE IF EXISTS users;

CREATE TABLE users(
    id int auto_increment primary key,
    name varchar(255) not null,
    nick varchar(255) not null unique,
    email varchar(255) not null unique,
    password varchar(255) not null,
    created_at timestamp default current_timestamp()
)ENGINE=INNODB;