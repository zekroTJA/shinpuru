SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";
CREATE DATABASE IF NOT EXISTS `shinpuru` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE `shinpuru`;

DROP TABLE IF EXISTS `guilds`;
CREATE TABLE `guilds` (
  `guildID` text NOT NULL,
  `prefix` text NOT NULL,
  `autorole` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `permissions`;
CREATE TABLE `permissions` (
  `roleID` text NOT NULL,
  `guildID` text NOT NULL,
  `permission` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `settings`;
CREATE TABLE `settings` (
  `setting` text NOT NULL,
  `value` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;