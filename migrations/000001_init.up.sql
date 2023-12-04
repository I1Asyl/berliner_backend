CREATE TABLE IF NOT EXISTS `user` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(255) UNIQUE NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `first_name` VARCHAR(255) NOT NULL,
    `last_name` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `team`(
  `id` INT NOT NULL AUTO_INCREMENT,
  `team_name` VARCHAR(255) UNIQUE NOT NULL,
  `team_leader_id` INT DEFAULT NULL,
  `team_description` TEXT NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`team_leader_id`) REFERENCES `user`(`id`) ON DELETE SET DEFAULT
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `membership`(
    `id` INT NOT NULL AUTO_INCREMENT,
    `team_id` INT NOT NULL,
    `user_id` INT NOT NULL,
    `is_editor` TINYINT(1) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`team_id`) REFERENCES `team`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `request`(
    `id` INT NOT NULL AUTO_INCREMENT,
    `team_id` INT NOT NULL,
    `user_id` INT NOT NULL,
    `is_accepted` TINYINT(1) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`team_id`) REFERENCES `team`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `following`(
    `id` INT NOT NULL AUTO_INCREMENT,
    `user_id` INT NOT NULL,
    `follower_id` INT NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`follower_id`) REFERENCES `user`(`id`) ON DELETE CASCADE
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `user_post` (
-- mandatory columns
    `id` INT NOT NULL AUTO_INCREMENT,
    `content` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `author_type` ENUM('user', 'team') NOT NULL,
    `is_public` TINYINT(1) NOT NULL DEFAULT 1,
    `user_id` INT NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    PRIMARY KEY (`id`)

    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `team_post` (
-- mandatory columns
    `id` INT NOT NULL AUTO_INCREMENT,
    `content` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `author_type` ENUM('user', 'team') NOT NULL,
    `is_public` TINYINT(1) NOT NULL DEFAULT 1,
    `team_id` INT NOT NULL,
    FOREIGN KEY (`team_id`) REFERENCES `team`(`id`) ON DELETE CASCADE,
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

