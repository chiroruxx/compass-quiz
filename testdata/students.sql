INSERT INTO teachers (id) VALUES (1), (2);
INSERT INTO classes (id, name) VALUES (1, 'クラスA'), (2, 'クラスB');
INSERT INTO class_teachers (class_id, teacher_id) VALUES (1, 1), (2, 2);
INSERT INTO students (id, name, class_id) VALUES
    (10, 'Williams', 1), (2, '鈴木', 2), (7, '小田原', 1), (4, '加藤', 1), (1, '佐藤', 1);
INSERT INTO student_accounts (student_id, login_id) VALUES
    (1, 'foo123'), (2, 'bar456'), (4, 'hoge0000'), (7, 'fizz9999'), (10, 'xxx:42');
