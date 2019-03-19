CREATE TABLE IF NOT EXISTS `guilds` (
  `iid` INTEGER PRIMARY KEY AUTOINCREMENT,
  `guildID` text NOT NULL DEFAULT '',
  `prefix` text NOT NULL DEFAULT '',
  `autorole` text NOT NULL DEFAULT '',
  `modlogchanID` text NOT NULL DEFAULT '',
  `voicelogchanID` text NOT NULL DEFAULT '',
  `muteRoleID` text NOT NULL DEFAULT '',
  `ghostPingMsg` text NOT NULL DEFAULT '',
  `jdoodleToken` text NOT NULL DEFAULT '',
  `backup` text NOT NULL DEFAULT '',
  `inviteBlock` text NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS `permissions` (
  `iid` INTEGER PRIMARY KEY AUTOINCREMENT,
  `roleID` text NOT NULL DEFAULT '',
  `guildID` text NOT NULL DEFAULT '',
  `permission` int(11) NOT NULL DEFAULT '0'
);

CREATE TABLE IF NOT EXISTS `reports` (
  `iid` INTEGER PRIMARY KEY AUTOINCREMENT,
  `id` text NOT NULL DEFAULT '',
  `type` int(11) NOT NULL DEFAULT '3',
  `guildID` text NOT NULL DEFAULT '',
  `executorID` text NOT NULL DEFAULT '',
  `victimID` text NOT NULL DEFAULT '',
  `msg` text NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS `settings` (
  `iid` INTEGER PRIMARY KEY AUTOINCREMENT,
  `setting` text NOT NULL DEFAULT '',
  `value` text NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS `starboard` (
  `iid` INTEGER PRIMARY KEY AUTOINCREMENT,
  `guildID` text NOT NULL DEFAULT '',
  `chanID` text NOT NULL DEFAULT '',
  `enabled` tinyint(1) NOT NULL DEFAULT '1',
  `minimum` int(11) NOT NULL DEFAULT '5'
);

CREATE TABLE IF NOT EXISTS `votes` (
  `iid` INTEGER PRIMARY KEY AUTOINCREMENT,
  `id` text NOT NULL DEFAULT '',
  `data` mediumtext NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS `twitchnotify` (
  `iid` INTEGER PRIMARY KEY AUTOINCREMENT,
  `guildID` text NOT NULL DEFAULT '',
  `channelID` text NOT NULL DEFAULT '',
  `twitchUserID` text NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS `backups` (
  `iid` INTEGER PRIMARY KEY AUTOINCREMENT,
  `guildID` text NOT NULL DEFAULT '',
  `timestamp` bigint(20) NOT NULL DEFAULT 0,
  `fileID` text NOT NULL DEFAULT ''
);