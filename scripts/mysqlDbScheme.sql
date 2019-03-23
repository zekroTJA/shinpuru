SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

CREATE TABLE IF NOT EXISTS `guilds` (
  `iid` int(11) NOT NULL AUTO_INCREMENT,
  `guildID` text NOT NULL,
  `prefix` text NOT NULL,
  `autorole` text NOT NULL,
  `modlogchanID` text NOT NULL,
  `voicelogchanID` text NOT NULL,
  `muteRoleID` text NOT NULL,
  `ghostPingMsg` text NOT NULL,
  `jdoodleToken` text NOT NULL,
  `backup` text NOT NULL,
  `inviteBlock` text NOT NULL,
  PRIMARY KEY (`iid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `permissions` (
  `iid` int(11) NOT NULL AUTO_INCREMENT,
  `roleID` text NOT NULL,
  `guildID` text NOT NULL,
  `permission` int(11) NOT NULL,
  PRIMARY KEY (`iid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `reports` (
  `iid` int(11) NOT NULL AUTO_INCREMENT,
  `id` text NOT NULL,
  `type` int(11) NOT NULL,
  `guildID` text NOT NULL,
  `executorID` text NOT NULL,
  `victimID` text NOT NULL,
  `msg` text NOT NULL,
  PRIMARY KEY (`iid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `settings` (
  `iid` int(11) NOT NULL AUTO_INCREMENT,
  `setting` text NOT NULL,
  `value` text NOT NULL,
  PRIMARY KEY (`iid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `starboard` (
  `iid` int(11) NOT NULL AUTO_INCREMENT,
  `guildID` text NOT NULL,
  `chanID` text NOT NULL,
  `enabled` tinyint(1) NOT NULL DEFAULT '1',
  `minimum` int(11) NOT NULL DEFAULT '5',
  PRIMARY KEY (`iid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `votes` (
  `iid` int(11) NOT NULL AUTO_INCREMENT,
  `id` text NOT NULL,
  `data` mediumtext NOT NULL,
  PRIMARY KEY (`iid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `twitchnotify` (
  `iid` int(11) NOT NULL AUTO_INCREMENT,
  `guildID` text NOT NULL,
  `channelID` text NOT NULL,
  `twitchUserID` text NOT NULL,
  PRIMARY KEY (`iid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `backups` (
  `iid` int(11) NOT NULL AUTO_INCREMENT,
  `guildID` text NOT NULL,
  `timestamp` bigint(20) NOT NULL,
  `fileID` text NOT NULL,
  PRIMARY KEY (`iid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
