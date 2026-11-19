CREATE TABLE classes (
    id BIGINT UNSIGNED PRIMARY KEY,
    name VARCHAR(191) NOT NULL
);

CREATE TABLE students (
    id BIGINT UNSIGNED PRIMARY KEY,
    name VARCHAR(191) NOT NULL,
    class_id BIGINT UNSIGNED NOT NULL,
    FOREIGN KEY (class_id) REFERENCES classes (id)
);

CREATE TABLE teachers (
    id BIGINT UNSIGNED PRIMARY KEY
);

CREATE TABLE class_teachers (
    class_id BIGINT UNSIGNED NOT NULL,
    teacher_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (class_id, teacher_id),
    FOREIGN KEY (class_id) REFERENCES classes (id),
    FOREIGN KEY (teacher_id) REFERENCES teachers (id)
);

CREATE TABLE student_accounts (
    student_id BIGINT UNSIGNED PRIMARY KEY,
    login_id VARCHAR(191) COLLATE utf8mb4_0900_bin NOT NULL UNIQUE,
    FOREIGN KEY (student_id) REFERENCES students (id)
);
