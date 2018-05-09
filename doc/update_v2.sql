-- 增加流程状态
ALTER TABLE f_flow ADD status INT DEFAULT 1 NULL;
ALTER TABLE f_flow
  MODIFY COLUMN status INT DEFAULT 1 AFTER parent_id;
