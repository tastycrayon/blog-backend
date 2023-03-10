-- +migrate Up
-- --------------------------------------------------------
-- Create Table Posts
--
CREATE TABLE IF NOT EXISTS `posts` (
`ID` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
`post_title` text COLLATE utf8mb4_unicode_520_ci NOT NULL,
`post_slug` varchar(200) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
`post_content` longtext COLLATE utf8mb4_unicode_520_ci NOT NULL,
`post_image` bigint(20) UNSIGNED NOT NULL DEFAULT '0',
`post_author` bigint(20) UNSIGNED NOT NULL DEFAULT '0',
`created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
`updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (`ID`),
INDEX `post_title` (`post_title`(50)),
INDEX `post_slug` (`post_slug`(50)),
INDEX `post_author` (`post_author`),
INDEX `id_date` (`created_at`, `ID`)
);
-- --------------------------------------------------------
-- Create Table Categories
--
CREATE TABLE IF NOT EXISTS `categories` (
`ID` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
`category_title` text COLLATE utf8mb4_unicode_520_ci NOT NULL,
`category_slug` varchar(200) COLLATE utf8mb4_unicode_520_ci NOT NULL,
`description` text COLLATE utf8mb4_unicode_520_ci NOT NULL,
`created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
`updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (`ID`),
INDEX `category_title` (`category_title`(50)),
INDEX `category_slug` (`category_slug`(50))
);
-- --------------------------------------------------------
-- Create Table users
--
CREATE TABLE IF NOT EXISTS `users` (
`ID` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
`user_name` varchar(60) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
`user_email` varchar(100) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
`user_pass` varchar(255) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
`user_activation_key` varchar(255) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
`user_role` tinyint NOT NULL DEFAULT '2',
`display_name` varchar(250) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
`user_image` bigint(20) UNSIGNED NOT NULL DEFAULT '0',
`created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
`updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (`ID`),
INDEX `user_name` (`user_name`(50)),
INDEX `user_email` (`user_email`(50))
);
-- --------------------------------------------------------
-- Create Table images
--
CREATE TABLE IF NOT EXISTS `images` (
`ID` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
`image_title` varchar(200) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
`image_url` varchar(200) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
`thumbnail_url` varchar(200) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
`height` INT UNSIGNED NOT NULL DEFAULT '0',
`width` INT UNSIGNED NOT NULL DEFAULT '0',
`created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
`updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (`ID`)
);
-- --------------------------------------------------------
-- Create Table post_category
--
CREATE TABLE IF NOT EXISTS `post_category` (
`ID` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
`post_id` bigint(20) UNSIGNED NOT NULL,
`category_id` bigint(20) UNSIGNED NOT NULL,
PRIMARY KEY (`ID`),
INDEX `post_id` (`post_id`),
INDEX `category_id` (`category_id`)
);
-- +migrate Down
DROP TABLE IF EXISTS `posts`;
DROP TABLE IF EXISTS `categories`;
DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `images`;