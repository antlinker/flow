CREATE TABLE test_apply_users
(
    id int(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id varchar(50),
    launcher varchar(50)
);
INSERT INTO test_apply_users (user_id, launcher) VALUES ('A002', 'A001');
INSERT INTO test_apply_users (user_id, launcher) VALUES ('A003', 'A001');