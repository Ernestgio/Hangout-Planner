-- Create "users" table
CREATE TABLE `users` (
  `id` char(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_users_deleted_at` (`deleted_at`),
  UNIQUE INDEX `idx_users_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "activities" table
CREATE TABLE `activities` (
  `id` char(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `user_id` char(36) NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_activities_user` (`user_id`),
  INDEX `idx_activities_deleted_at` (`deleted_at`),
  UNIQUE INDEX `idx_activities_name` (`name`),
  CONSTRAINT `fk_activities_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "hangouts" table
CREATE TABLE `hangouts` (
  `id` char(36) NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` text NULL,
  `date` datetime(3) NOT NULL,
  `status` varchar(50) NOT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `user_id` char(36) NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_users_hangouts` (`user_id`),
  INDEX `idx_hangouts_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_users_hangouts` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "hangout_activities" table
CREATE TABLE `hangout_activities` (
  `activity_id` char(36) NOT NULL,
  `hangout_id` char(36) NOT NULL,
  PRIMARY KEY (`activity_id`, `hangout_id`),
  INDEX `fk_hangout_activities_hangout` (`hangout_id`),
  CONSTRAINT `fk_hangout_activities_activity` FOREIGN KEY (`activity_id`) REFERENCES `activities` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_hangout_activities_hangout` FOREIGN KEY (`hangout_id`) REFERENCES `hangouts` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
