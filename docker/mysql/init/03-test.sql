CREATE DATABASE compass_quiz_test CHARACTER SET utf8mb4;
USE compass_quiz_test;

SOURCE /docker-entrypoint-initdb.d/01-schema.sql

-- Escape underscores so the database name in GRANT is matched literally.
GRANT SELECT ON `compass\_quiz\_test`.* TO 'compass'@'%';
