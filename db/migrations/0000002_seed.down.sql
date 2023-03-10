-- +migrate Down
TRUNCATE TABLE `posts`;
TRUNCATE TABLE `categories`;
TRUNCATE TABLE `users`;
TRUNCATE TABLE `images`;