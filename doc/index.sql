ALTER TABLE `f_node_instance` ADD UNIQUE INDEX `record_id` (`record_id`);
ALTER TABLE `f_node_instance` ADD INDEX `flow_instance_id` (`flow_instance_id`);
ALTER TABLE `f_node_instance` ADD INDEX `node_id` (`node_id`);
ALTER TABLE `f_node_instance` ADD INDEX `deleted` (`deleted`);
ALTER TABLE `f_node_instance` ADD INDEX `status` (`status`);

ALTER TABLE `f_flow_instance` ADD UNIQUE INDEX `record_id` (`record_id`);
ALTER TABLE `f_flow_instance` ADD INDEX `flow_id` (`flow_id`);
ALTER TABLE `f_flow_instance` ADD INDEX `status` (`status`);
ALTER TABLE `f_flow_instance` ADD INDEX `deleted` (`deleted`);

ALTER TABLE `f_node` ADD UNIQUE INDEX `record_id` (`record_id`);
ALTER TABLE `f_node` ADD INDEX `deleted` (`deleted`);

ALTER TABLE `f_form` ADD UNIQUE INDEX `record_id` (`record_id`);
ALTER TABLE `f_form` ADD INDEX `deleted` (`deleted`);

ALTER TABLE `f_flow` ADD UNIQUE INDEX `record_id` (`record_id`);
ALTER TABLE `f_flow` ADD INDEX `code` (`code`);
ALTER TABLE `f_flow` ADD INDEX `flag` (`flag`);
ALTER TABLE `f_flow` ADD INDEX `deleted` (`deleted`);

ALTER TABLE `f_node_candidate` ADD UNIQUE INDEX `record_id` (`record_id`);
ALTER TABLE `f_node_candidate` ADD INDEX `candidate_id` (`candidate_id`);
ALTER TABLE `f_node_candidate` ADD INDEX `deleted` (`deleted`);

ALTER TABLE `f_node_router` ADD UNIQUE INDEX `record_id` (`record_id`);
ALTER TABLE `f_node_router` ADD INDEX `source_node_id` (`source_node_id`);
ALTER TABLE `f_node_router` ADD INDEX `deleted` (`deleted`);

ALTER TABLE `f_node_assignment` ADD UNIQUE INDEX `record_id` (`record_id`);
ALTER TABLE `f_node_assignment` ADD INDEX `node_id` (`node_id`);
ALTER TABLE `f_node_assignment` ADD INDEX `deleted` (`deleted`);
