CREATE TABLE `user`
(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ユーザーの識別子',
    `name` VARCHAR(20) NOT NULL COMMENT 'ユーザーのフルネーム',
    `user_name` VARCHAR(20) NOT NULL COMMENT 'ユーザーの名前',
    `password`VARCHAR(80) NOT NULL COMMENT 'パスワードハッシュ',
    `role` VARCHAR(80) NOT NULL COMMENT 'ロール',
    `email` VARCHAR(80) NOT NULL COMMENT 'メールアドレス',
    `address` VARCHAR(200) NOT NULL COMMENT '住所',
    `phone` VARCHAR(20) NOT NULL COMMENT '電話番号',
    `website` VARCHAR(200) NOT NULL COMMENT 'Webサイト',
    `company` VARCHAR(200) NOT NULL COMMENT '会社',
    `created` DATETIME(6) NOT NULL COMMENT 'レコード作成日時',
    `modified` DATETIME(6) NOT NULL COMMENT 'レコード編集日時',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uix_name` (`user_name`) USING BTREE
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ユーザー';
