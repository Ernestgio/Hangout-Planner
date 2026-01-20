-- Modify "memories" table
ALTER TABLE `memories` ADD COLUMN `file_id` char(36) NULL AFTER `name`, ADD INDEX `idx_memories_file_id` (`file_id`);
