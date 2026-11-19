SET NAMES utf8mb4;

START TRANSACTION;

INSERT INTO teachers (id) VALUES (1), (2), (3);

INSERT INTO classes (id, name) VALUES
    (1, 'クラスA'),
    (2, 'クラスB'),
    (3, 'クラスC');

INSERT INTO students (id, name, class_id) VALUES
    (1, '佐藤', 1),
    (2, '鈴木', 2),
    (3, '田中', 3),
    (4, '加藤', 1),
    (5, '太田', 2),
    (6, '佐々木', 3),
    (7, '小田原', 1),
    (8, 'Smith', 2),
    (9, 'Johnson', 3),
    (10, 'Williams', 1);

INSERT INTO class_teachers (class_id, teacher_id) VALUES
    (1, 1),
    (2, 2),
    (3, 3);

INSERT INTO student_accounts (student_id, login_id) VALUES
    (1, 'foo123'),
    (2, 'bar456'),
    (3, 'baz789'),
    (4, 'hoge0000'),
    (5, 'fuga1234'),
    (6, 'piyo5678'),
    (7, 'fizz9999'),
    (8, 'buzz777'),
    (9, 'fizzbuzz#123'),
    (10, 'xxx:42');

COMMIT;
