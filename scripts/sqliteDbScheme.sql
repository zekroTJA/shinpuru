DROP TABLE IF EXISTS `guilds`;
CREATE TABLE `guilds` (
  `guildID` text NOT NULL DEFAULT '',
  `prefix` text NOT NULL DEFAULT '',
  `autorole` text NOT NULL DEFAULT '',
  `modlogchanID` text NOT NULL DEFAULT '',
  `muteRoleID` text NOT NULL DEFAULT ''
);

DROP TABLE IF EXISTS `permissions`;
CREATE TABLE `permissions` (
  `roleID` text NOT NULL DEFAULT '',
  `guildID` text NOT NULL DEFAULT '',
  `permission` int(11) NOT NULL DEFAULT '0'
);

DROP TABLE IF EXISTS `reports`;
CREATE TABLE `reports` (
  `id` text NOT NULL DEFAULT '',
  `type` int(11) NOT NULL DEFAULT '3',
  `guildID` text NOT NULL DEFAULT '',
  `executorID` text NOT NULL DEFAULT '',
  `victimID` text NOT NULL DEFAULT '',
  `msg` text NOT NULL DEFAULT ''
);

DROP TABLE IF EXISTS `settings`;
CREATE TABLE `settings` (
  `setting` text NOT NULL DEFAULT '',
  `value` text NOT NULL DEFAULT ''
);

DROP TABLE IF EXISTS `starboard`;
CREATE TABLE `starboard` (
  `guildID` text NOT NULL DEFAULT '',
  `chanID` text NOT NULL DEFAULT '',
  `enabled` tinyint(1) NOT NULL DEFAULT '1',
  `minimum` int(11) NOT NULL DEFAULT '5'
);

DROP TABLE IF EXISTS `votes`;
CREATE TABLE `votes` (
  `ID` text NOT NULL DEFAULT '',
  `data` mediumtext NOT NULL DEFAULT ''
);
