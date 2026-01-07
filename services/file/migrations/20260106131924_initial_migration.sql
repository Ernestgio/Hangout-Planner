-- Create "memory_files" table
CREATE TABLE `memory_files` (
  `id` char(36) NOT NULL,
  `original_name` varchar(255) NOT NULL,
  `file_extension` varchar(10) NOT NULL,
  `storage_path` varchar(500) NOT NULL,
  `file_size` bigint NOT NULL,
  `mime_type` varchar(100) NOT NULL,
  `file_status` varchar(50) NOT NULL,
  `created_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `memory_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_memory_files_deleted_at` (`deleted_at`),
  UNIQUE INDEX `idx_memory_files_memory_id` (`memory_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
