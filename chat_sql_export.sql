/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


--
CREATE DATABASE IF NOT EXISTS `chat` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;
USE `chat`;

-- Дамп структуры для таблица chat.api_log
CREATE TABLE IF NOT EXISTS `api_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Идентификатор',
  `created` timestamp NOT NULL DEFAULT current_timestamp(),
  `uri` text NOT NULL,
  `request` text NOT NULL,
  `data_size` int(10) unsigned DEFAULT NULL,
  `response` text DEFAULT NULL,
  `ip` varchar(255) NOT NULL,
  `app_version` varchar(32) DEFAULT NULL,
  `auth_token` text DEFAULT NULL,
  `useragent` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `created` (`created`),
  KEY `ip` (`ip`)
) ENGINE=InnoDB AUTO_INCREMENT=834 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица chat.auth
CREATE TABLE IF NOT EXISTS `auth` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `token` varchar(2000) NOT NULL,
  `ip` varchar(50) NOT NULL,
  `time_expiry` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `logout` tinyint(1) NOT NULL DEFAULT 0,
  `updated` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `source` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=444 DEFAULT CHARSET=utf8;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица chat.chats
CREATE TABLE IF NOT EXISTS `chats` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(150) NOT NULL,
  `is_group` int(11) NOT NULL DEFAULT 0,
  `logo` int(11) DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8mb4;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица chat.files
CREATE TABLE IF NOT EXISTS `files` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'код',
  `guid` varchar(32) NOT NULL COMMENT 'имя фото в адр строке',
  `type` varchar(55) DEFAULT NULL,
  `parent_id` int(11) unsigned DEFAULT NULL COMMENT '{type}_id',
  `name` varchar(255) DEFAULT NULL COMMENT 'имя файла',
  `path` varchar(255) DEFAULT NULL COMMENT 'путь',
  `fullname` varchar(255) DEFAULT NULL COMMENT 'полный путь',
  `extension` varchar(50) DEFAULT NULL,
  `mime_type` varchar(255) DEFAULT NULL,
  `size` int(11) unsigned DEFAULT NULL COMMENT 'размер файла',
  `mtype` varchar(255) DEFAULT NULL COMMENT 'тип файла',
  `index` int(11) unsigned NOT NULL DEFAULT 0,
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT 'создан',
  `updated` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'изменен',
  PRIMARY KEY (`id`),
  KEY `parent_id` (`parent_id`),
  KEY `guid` (`guid`)
) ENGINE=InnoDB AUTO_INCREMENT=245 DEFAULT CHARSET=utf8;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица chat.file_parts
CREATE TABLE IF NOT EXISTS `file_parts` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `file_id` int(11) unsigned NOT NULL COMMENT 'код из табл files',
  `name` varchar(255) DEFAULT NULL COMMENT 'имя',
  `fullname` varchar(255) NOT NULL COMMENT 'полный путь',
  `size` int(11) unsigned DEFAULT NULL COMMENT 'размер',
  `mtime` int(11) unsigned DEFAULT NULL,
  `index` int(11) unsigned NOT NULL,
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT 'соднание',
  `updated` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'изменение',
  PRIMARY KEY (`id`),
  KEY `file_id` (`file_id`),
  CONSTRAINT `file_parts_ibfk_1` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица chat.members
CREATE TABLE IF NOT EXISTS `members` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `chat_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `is_read` tinyint(4) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=109 DEFAULT CHARSET=utf8mb4;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица chat.messages
CREATE TABLE IF NOT EXISTS `messages` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `chat_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `message` longtext NOT NULL,
  `attachment` int(11) NOT NULL DEFAULT 0,
  `readed` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`readed`)),
  `created` timestamp NOT NULL DEFAULT current_timestamp(),
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=338 DEFAULT CHARSET=utf8mb4;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица chat.users
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'уникальный код пользователя',
  `presentation` varchar(255) DEFAULT NULL COMMENT 'ФИО',
  `login` varchar(50) DEFAULT NULL COMMENT 'логин',
  `password` varchar(32) NOT NULL COMMENT 'пароль',
  `email` varchar(80) DEFAULT NULL COMMENT 'электронная почта',
  `phone` varchar(20) DEFAULT NULL COMMENT 'телефон',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'статус',
  `role` tinyint(1) NOT NULL COMMENT 'роль',
  `avatar` int(11) DEFAULT NULL,
  `theme` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0-light 1-dark',
  `theme_color` varchar(80) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8;

-- Экспортируемые данные не выделены.

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;


INSERT INTO `users` (`id`, `presentation`, `login`, `password`, `email`, `phone`, `status`, `role`, `avatar`, `theme`, `theme_color`) VALUES (NULL, 'developer', 'test1', '098af2c9929f1f286f7c72b1200bee54', 'test@ya.ru', NULL, '1', '1', '243', '1', '#9538FF')
