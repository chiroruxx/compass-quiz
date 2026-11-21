-- name: ListStudentsWithAccounts :many
SELECT students.id, students.name, student_accounts.login_id,
       classes.id AS class_id, classes.name AS class_name
FROM students
JOIN classes ON classes.id = students.class_id
JOIN class_teachers ON class_teachers.class_id = classes.id
JOIN student_accounts ON student_accounts.student_id = students.id
WHERE class_teachers.teacher_id = ?
ORDER BY students.id ASC;
