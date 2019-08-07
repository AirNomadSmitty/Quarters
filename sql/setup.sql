CREATE TABLE `bet_closes` (
  `bet_close_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `bet_id` int(10) unsigned NOT NULL,
  `winning_position_id` int(10) unsigned DEFAULT '0',
  `losing_position_id` int(10) unsigned DEFAULT '0',
  PRIMARY KEY (`bet_close_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `bets` (
  `bet_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `base_value` decimal(6,2) unsigned NOT NULL DEFAULT '0.25',
  `closed` tinyint(3) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`bet_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `positions` (
  `position_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `bet_id` int(11) unsigned NOT NULL,
  `description` varchar(128) NOT NULL DEFAULT '',
  `odds_multiplier` decimal(6,2) unsigned NOT NULL DEFAULT '1.00',
  PRIMARY KEY (`position_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `users` (
  `user_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(128) NOT NULL DEFAULT '',
  `password_hash` varchar(128) NOT NULL DEFAULT '',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;


CREATE TABLE `users_to_positions` (
  `users_to_position_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL DEFAULT '0',
  `position_id` int(10) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`users_to_position_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
