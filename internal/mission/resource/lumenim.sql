CREATE TABLE `admin` (
 `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
 `username` varchar(20) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户昵称',
 `password` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户密码',
 `avatar` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户头像',
 `gender` tinyint unsigned NOT NULL DEFAULT '3' COMMENT '用户性别[1:男;2:女;3:未知;]',
 `mobile` varchar(11) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
 `email` varchar(30) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户邮箱',
 `motto` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户座右铭',
 `last_login_at` datetime NOT NULL COMMENT '最后一次登录时间',
 `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态[1:正常;2:停用;]',
 `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
 `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
 PRIMARY KEY (`id`) USING BTREE,
 UNIQUE KEY `uk_username` (`username`) USING BTREE,
 UNIQUE KEY `uk_email` (`email`) USING BTREE,
 KEY `idx_created_at` (`created_at`) USING BTREE,
 KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='管理员表';;

CREATE TABLE IF NOT EXISTS `article`
(
    `id`          int unsigned     NOT NULL AUTO_INCREMENT COMMENT '文章ID',
    `user_id`     int unsigned     NOT NULL COMMENT '用户ID',
    `class_id`    int unsigned     NOT NULL DEFAULT '0' COMMENT '分类ID',
    `tags_id`     varchar(128)     NOT NULL DEFAULT '' COMMENT '笔记关联标签',
    `title`       varchar(255)     NOT NULL COMMENT '文章标题',
    `abstract`    varchar(255)     NOT NULL DEFAULT '' COMMENT '文章摘要',
    `image`       varchar(255)     NOT NULL DEFAULT '' COMMENT '文章首图',
    `is_asterisk` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '是否星标文章[1:是;2:否;]',
    `status`      tinyint unsigned NOT NULL DEFAULT '1' COMMENT '笔记状态[1:正常;2:已删除;]',
    `md_content`  longtext         NOT NULL COMMENT 'markdown 内容',
    `created_at`  datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`  datetime                  DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_userid_classid_title` (`user_id`, `class_id`, `title`),
    KEY `idx_updated_at` (`updated_at`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户文章表';;


CREATE TABLE IF NOT EXISTS `article_annex`
(
    `id`            int unsigned     NOT NULL AUTO_INCREMENT COMMENT '文件ID',
    `user_id`       int unsigned     NOT NULL COMMENT '上传文件的用户ID',
    `article_id`    int unsigned     NOT NULL COMMENT '笔记ID',
    `drive`         tinyint unsigned NOT NULL DEFAULT '1' COMMENT '文件驱动[1:local;2:cos;]',
    `suffix`        varchar(10)      NOT NULL DEFAULT '' COMMENT '文件后缀名',
    `size`          bigint unsigned  NOT NULL DEFAULT '0' COMMENT '文件大小',
    `path`          varchar(500)     NOT NULL COMMENT '文件地址（相对地址）',
    `original_name` varchar(100)     NOT NULL DEFAULT '' COMMENT '原文件名',
    `status`        tinyint unsigned NOT NULL DEFAULT '1' COMMENT '附件状态[1:正常;2:已删除;]',
    `created_at`    datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`    datetime                  DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_userid_articleid` (`user_id`, `article_id`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='文章附件信息表';;


CREATE TABLE IF NOT EXISTS `article_class`
(
    `id`         int unsigned     NOT NULL AUTO_INCREMENT COMMENT '文章分类ID',
    `user_id`    int unsigned     NOT NULL COMMENT '用户ID',
    `class_name` varchar(64)      NOT NULL COMMENT '分类名',
    `sort`       tinyint unsigned NOT NULL DEFAULT '1' COMMENT '排序',
    `is_default` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '默认分类[1:是;2:否]',
    `created_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id_class_name` (`user_id`, `class_name`) USING BTREE,
    KEY `uk_user_id_sort` (`user_id`, `sort`),
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='文章分类表';;


CREATE TABLE IF NOT EXISTS `article_tag`
(
    `id`         int unsigned     NOT NULL AUTO_INCREMENT COMMENT '文章分类ID',
    `user_id`    int unsigned     NOT NULL COMMENT '用户ID',
    `tag_name`   varchar(20)      NOT NULL COMMENT '标签名',
    `sort`       tinyint unsigned NOT NULL DEFAULT '1' COMMENT '排序',
    `created_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_userid` (`user_id`),
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='文章标签表';;


CREATE TABLE IF NOT EXISTS `contact`
(
    `id`         int unsigned     NOT NULL AUTO_INCREMENT COMMENT '关系ID',
    `user_id`    int unsigned     NOT NULL DEFAULT '0' COMMENT '用户id',
    `friend_id`  int unsigned     NOT NULL DEFAULT '0' COMMENT '好友id',
    `remark`     varchar(64)      NOT NULL DEFAULT '' COMMENT '好友的备注',
    `status`     tinyint unsigned NOT NULL DEFAULT '0' COMMENT '好友状态 [0:否;1:是]',
    `group_id`   int unsigned     NOT NULL DEFAULT '0' COMMENT '分组ID',
    `created_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_friend` (`user_id`, `friend_id`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户好友关系表';;


CREATE TABLE IF NOT EXISTS `contact_apply`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT COMMENT '申请ID',
    `user_id`    int unsigned NOT NULL COMMENT '申请人ID',
    `friend_id`  int unsigned NOT NULL COMMENT '被申请人',
    `remark`     varchar(64)  NOT NULL DEFAULT '' COMMENT '申请备注',
    `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`) USING BTREE,
    KEY `idx_friend_id` (`friend_id`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户添加好友申请表';;

CREATE TABLE IF NOT EXISTS `contact_group`
(
    `id`         int          NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id`    int unsigned NOT NULL COMMENT '用户ID',
    `name`       varchar(64)  NOT NULL COMMENT '分组名称',
    `sort`       int unsigned NOT NULL DEFAULT '1' COMMENT '排序',
    `num`        int unsigned NOT NULL DEFAULT '0' COMMENT '成员总数',
    `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id_name` (`user_id`, `name`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='联系人分组';;


CREATE TABLE IF NOT EXISTS `emoticon`
(
    `id`         int unsigned     NOT NULL AUTO_INCREMENT COMMENT '表情分组ID',
    `name`       varchar(64)      NOT NULL COMMENT '分组名称',
    `icon`       varchar(255)     NOT NULL DEFAULT '' COMMENT '分组图标',
    `status`     tinyint unsigned NOT NULL DEFAULT '0' COMMENT '分组状态[1:正常;2:已禁用;]',
    `created_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name` (`name`),
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='表情包分组';;

CREATE TABLE IF NOT EXISTS `emoticon_item`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT COMMENT '表情包详情ID',
    `emoticon_id` int unsigned NOT NULL COMMENT '表情分组ID（0:用户自定义上传）',
    `user_id`     int unsigned NOT NULL COMMENT '用户ID（0:代码系统表情包）',
    `describe`    varchar(64)  NOT NULL DEFAULT '' COMMENT '表情描述',
    `url`         varchar(255) NOT NULL COMMENT '图片链接',
    `created_at`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE,
    KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='表情包详情表';;

CREATE TABLE IF NOT EXISTS `file_upload`
(
    `id`            int unsigned     NOT NULL AUTO_INCREMENT COMMENT '临时文件ID',
    `type`          tinyint unsigned NOT NULL DEFAULT '1' COMMENT '文件属性[1:合并文件;2:拆分文件]',
    `drive`         tinyint unsigned NOT NULL DEFAULT '1' COMMENT '驱动类型[1:local;2:cos;]',
    `upload_id`     varchar(128)     NOT NULL DEFAULT '' COMMENT '临时文件hash名',
    `user_id`       int unsigned     NOT NULL DEFAULT '0' COMMENT '上传的用户ID',
    `original_name` varchar(64)      NOT NULL DEFAULT '' COMMENT '原文件名',
    `split_index`   int unsigned     NOT NULL DEFAULT '0' COMMENT '当前索引块',
    `split_num`     int unsigned     NOT NULL DEFAULT '0' COMMENT '总上传索引块',
    `path`          varchar(255)     NOT NULL DEFAULT '' COMMENT '临时保存路径',
    `file_ext`      varchar(16)      NOT NULL DEFAULT '' COMMENT '文件后缀名',
    `file_size`     int unsigned     NOT NULL COMMENT '文件大小',
    `is_delete`     tinyint unsigned NOT NULL DEFAULT '0' COMMENT '文件是否删除[1:是;2:否;] ',
    `attr`          json             NOT NULL COMMENT '额外参数json',
    `created_at`    datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `updated_at`    datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id_hash_name` (`user_id`, `upload_id`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='文件拆分数据表';;


CREATE TABLE IF NOT EXISTS `group`
(
    `id`         int unsigned      NOT NULL AUTO_INCREMENT COMMENT '群ID',
    `type`       tinyint unsigned  NOT NULL DEFAULT '1' COMMENT '群类型[1:普通群;2:企业群;]',
    `name`       varchar(64)       NOT NULL DEFAULT '' COMMENT '群名称',
    `profile`    varchar(128)      NOT NULL DEFAULT '' COMMENT '群介绍',
    `avatar`     varchar(255)      NOT NULL DEFAULT '' COMMENT '群头像',
    `max_num`    smallint unsigned NOT NULL DEFAULT '200' COMMENT '最大群成员数量',
    `is_overt`   tinyint unsigned  NOT NULL DEFAULT '2' COMMENT '是否公开可见[1:是;2:否;]',
    `is_mute`    tinyint unsigned  NOT NULL DEFAULT '2' COMMENT '是否全员禁言 [1:是;2:否;] 提示:不包含群主或管理员',
    `is_dismiss` tinyint unsigned  NOT NULL DEFAULT '2' COMMENT '是否已解散[1:是;2:否;]',
    `creator_id` int unsigned      NOT NULL COMMENT '创建者ID(群主ID)',
    `created_at` datetime          NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime          NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='聊天群';;

CREATE TABLE IF NOT EXISTS `group_apply`
(
    `id`         int unsigned     NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `group_id`   int unsigned     NOT NULL COMMENT '群组ID',
    `user_id`    int unsigned     NOT NULL COMMENT '用户ID',
    `status`     tinyint unsigned NOT NULL DEFAULT '1' COMMENT '申请状态[1:待审核;2:已通过;3:不通过;]',
    `remark`     varchar(255)     NOT NULL DEFAULT '' COMMENT '备注信息',
    `reason`     varchar(255)     NOT NULL DEFAULT '' COMMENT '拒绝原因',
    `created_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_group_id_user_id` (`group_id`, `user_id`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='入群申请';;


CREATE TABLE IF NOT EXISTS `group_member`
(
    `id`         int unsigned     NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `group_id`   int unsigned     NOT NULL COMMENT '群组ID',
    `user_id`    int unsigned     NOT NULL COMMENT '用户ID',
    `leader`     tinyint unsigned NOT NULL DEFAULT '3' COMMENT '成员属性[1:群主;1:管理员;3:普通成员]',
    `user_card`  varchar(64)      NOT NULL DEFAULT '' COMMENT '群名片',
    `is_quit`    tinyint unsigned NOT NULL DEFAULT '2' COMMENT '是否退群[1:是;2:否;]',
    `is_mute`    tinyint unsigned NOT NULL DEFAULT '2' COMMENT '是否禁言[1:是;2:否;]',
    `join_time`  datetime                  DEFAULT NULL COMMENT '入群时间',
    `created_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_group_id_user_id` (`group_id`, `user_id`) USING BTREE,
    KEY `idx_user_id` (`user_id`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='群聊成员';;


CREATE TABLE IF NOT EXISTS `group_notice`
(
    `id`            int unsigned     NOT NULL AUTO_INCREMENT COMMENT '公告ID',
    `group_id`      int unsigned     NOT NULL COMMENT '群组ID',
    `creator_id`    int unsigned     NOT NULL COMMENT '创建者用户ID',
    `modify_id`     int              NOT NULL COMMENT '修改者ID',
    `content`       longtext         NOT NULL COMMENT '公告内容',
    `confirm_users` json                      DEFAULT NULL COMMENT '已确认成员',
    `is_confirm`    tinyint unsigned NOT NULL DEFAULT '1' COMMENT '是否需群成员确认公告[1:是;2:否;]',
    `created_at`    datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`, `modify_id`) USING BTREE,
    UNIQUE KEY `un_group_id` (`group_id`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='群组公告表';;


CREATE TABLE IF NOT EXISTS `group_vote`
(
    `id`            int unsigned NOT NULL AUTO_INCREMENT COMMENT '投票ID',
    `group_id`      int unsigned NOT NULL COMMENT '群组ID',
    `user_id`       int unsigned NOT NULL COMMENT '用户ID(创建人)',
    `title`         varchar(64)  NOT NULL COMMENT '投票标题',
    `answer_mode`   int unsigned NOT NULL COMMENT '答题模式[1:单选;2:多选;]',
    `answer_option` json         NOT NULL COMMENT '答题选项',
    `answer_num`    int unsigned NOT NULL DEFAULT '0' COMMENT '应答人数',
    `answered_num`  int unsigned NOT NULL DEFAULT '0' COMMENT '已答人数',
    `is_anonymous`  int unsigned NOT NULL DEFAULT '2' COMMENT '匿名投票[1:是;2:否;]',
    `status`        int unsigned NOT NULL DEFAULT '1' COMMENT '投票状态[1:投票中;2:已完成;]',
    `created_at`    datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE,
    KEY `idx_groupid` (`group_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='群投票表';;


CREATE TABLE IF NOT EXISTS `group_vote_answer`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT COMMENT '答题ID',
    `vote_id`    int unsigned NOT NULL COMMENT '投票ID',
    `user_id`    int unsigned NOT NULL COMMENT '用户ID',
    `option`     char(1)      NOT NULL COMMENT '投票选项[A、B、C 、D、E、F]',
    `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '答题时间',
    PRIMARY KEY (`id`),
    KEY `idx_vote_id_user_id` (`vote_id`, `user_id`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='投票详情统计表';;


CREATE TABLE IF NOT EXISTS `organize`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `user_id`     int unsigned NOT NULL COMMENT '用户id',
    `dept_id`     int unsigned NOT NULL COMMENT '部门ID',
    `position_id` int unsigned NOT NULL COMMENT '岗位ID',
    `created_at`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`, `position_id`, `dept_id`) USING BTREE,
    UNIQUE KEY `uk_user_id` (`user_id`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='组织表';;


CREATE TABLE IF NOT EXISTS `organize_dept`
(
    `dept_id`    int              NOT NULL AUTO_INCREMENT COMMENT '部门id',
    `parent_id`  int              NOT NULL DEFAULT '0' COMMENT '父部门id',
    `ancestors`  varchar(128)     NOT NULL DEFAULT '' COMMENT '祖级列表',
    `dept_name`  varchar(64)      NOT NULL DEFAULT '' COMMENT '部门名称',
    `order_num`  int unsigned     NOT NULL DEFAULT '1' COMMENT '显示顺序',
    `leader`     varchar(64)      NOT NULL COMMENT '负责人',
    `phone`      varchar(11)      NOT NULL COMMENT '联系电话',
    `email`      varchar(64)      NOT NULL COMMENT '邮箱',
    `status`     tinyint          NOT NULL DEFAULT '1' COMMENT '部门状态[1:正常;2:停用]',
    `is_deleted` tinyint unsigned NOT NULL DEFAULT '2' COMMENT '是否删除[1:是;2:否;]',
    `created_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`dept_id`),
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='部门表';;


CREATE TABLE IF NOT EXISTS `organize_position`
(
    `position_id` int              NOT NULL AUTO_INCREMENT COMMENT '岗位ID',
    `post_code`   varchar(32)      NOT NULL COMMENT '岗位编码',
    `post_name`   varchar(64)      NOT NULL COMMENT '岗位名称',
    `sort`        int unsigned     NOT NULL DEFAULT '1' COMMENT '显示顺序',
    `status`      tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态[1:正常;2:停用;]',
    `remark`      varchar(255)     NOT NULL DEFAULT '' COMMENT '备注',
    `created_at`  datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`position_id`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='岗位信息表';;


CREATE TABLE IF NOT EXISTS `robot`
(
    `id`         int unsigned     NOT NULL AUTO_INCREMENT COMMENT '机器人ID',
    `user_id`    int unsigned     NOT NULL COMMENT '关联用户ID',
    `robot_name` varchar(64)      NOT NULL DEFAULT '' COMMENT '机器人名称',
    `describe`   varchar(255)     NOT NULL DEFAULT '' COMMENT '描述信息',
    `logo`       varchar(255)     NOT NULL DEFAULT '' COMMENT '机器人logo',
    `is_talk`    tinyint unsigned NOT NULL DEFAULT '2' COMMENT '可发送消息[1:是;2:否;]',
    `status`     tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态[1:正常;2:已禁用;3:已删除;]',
    `type`       tinyint unsigned NOT NULL DEFAULT '0' COMMENT '机器人类型',
    `created_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_type` (`type`),
    UNIQUE KEY `uk_user_id` (`user_id`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='聊天机器人表';;


CREATE TABLE IF NOT EXISTS `talk_group_message`
(
    `id`         bigint unsigned  NOT NULL AUTO_INCREMENT COMMENT '聊天记录ID',
    `msg_id`     varchar(64)      NOT NULL COMMENT '消息ID',
    `sequence`   bigint unsigned  NOT NULL COMMENT '消息时序ID（消息排序）',
    `msg_type`   int unsigned     NOT NULL DEFAULT '1' COMMENT '消息类型',
    `group_id`   int unsigned     NOT NULL COMMENT '群组ID',
    `from_id`    int unsigned     NOT NULL COMMENT '消息发送者ID',
    `is_revoked` tinyint unsigned NOT NULL DEFAULT '2' COMMENT '是否撤回[1:是;2:否;]',
    `extra`      json             NOT NULL COMMENT '消息扩展字段',
    `quote`      json             NOT NULL COMMENT '引用消息',
    `send_time`  datetime         NOT NULL COMMENT '发送时间',
    `created_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_group_id_sequence` (`group_id`, `sequence`) USING BTREE,
    UNIQUE KEY `uk_msgid` (`msg_id`),
    KEY `idx_updated_at` (`updated_at`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='群聊消息记录表';;

CREATE TABLE IF NOT EXISTS `talk_group_message_del`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT,
    `user_id`    int unsigned NOT NULL COMMENT '用户ID',
    `group_id`   int unsigned NOT NULL COMMENT '群ID',
    `msg_id`     varchar(64)  NOT NULL COMMENT '聊天记录ID',
    `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id_msg_id` (`user_id`, `msg_id`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='群聊消息记录表-删除记录关系表';;

CREATE TABLE IF NOT EXISTS `talk_session`
(
    `id`         int unsigned     NOT NULL AUTO_INCREMENT COMMENT '聊天列表ID',
    `talk_mode`  tinyint unsigned NOT NULL DEFAULT '1' COMMENT '聊天类型[1:私信;2:群聊;]',
    `user_id`    int unsigned     NOT NULL DEFAULT '0' COMMENT '用户ID',
    `to_from_id` int unsigned     NOT NULL COMMENT '接收者ID（用户ID 或 群ID）',
    `is_top`     tinyint unsigned NOT NULL DEFAULT '2' COMMENT '是否置顶[1:是;2:否]',
    `is_disturb` tinyint unsigned NOT NULL DEFAULT '2' COMMENT '消息免打扰[1:是;2:否]',
    `is_delete`  tinyint unsigned NOT NULL DEFAULT '2' COMMENT '是否删除[1:是;2:否]',
    `is_robot`   tinyint unsigned NOT NULL DEFAULT '2' COMMENT '是否机器人[1:是;2:否]',
    `created_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id_receiver_id_talk_type` (`user_id`, `to_from_id`, `talk_mode`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='会话列表';;


CREATE TABLE IF NOT EXISTS `talk_user_message`
(
    `id`         bigint unsigned  NOT NULL AUTO_INCREMENT COMMENT '聊天记录ID',
    `msg_id`     varchar(64)      NOT NULL COMMENT '消息ID',
    `org_msg_id` varchar(64)      NOT NULL COMMENT '原消息ID',
    `sequence`   bigint           NOT NULL COMMENT '消息时序ID（消息排序）',
    `msg_type`   int unsigned     NOT NULL DEFAULT '1' COMMENT '消息类型',
    `user_id`    int unsigned     NOT NULL COMMENT '用户ID',
    `from_id`    int unsigned     NOT NULL COMMENT '消息发送者ID',
    `to_from_id` int unsigned     NOT NULL COMMENT '接收者ID',
    `is_revoked` tinyint unsigned NOT NULL DEFAULT '2' COMMENT '是否撤回[1:是;2:否;]',
    `is_deleted` tinyint unsigned NOT NULL DEFAULT '2' COMMENT '是否删除[1:是;2:否;]',
    `extra`      json             NOT NULL COMMENT '消息扩展字段',
    `quote`      json             NOT NULL COMMENT '引用消息',
    `send_time`  datetime         NOT NULL COMMENT '发送时间',
    `created_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id_friend_id_sequence` (`user_id`, `to_from_id`, `sequence`) USING BTREE,
    UNIQUE KEY `uk_msgid` (`msg_id`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE,
    KEY `idx_updated_at` (`updated_at`) USING BTREE,
    KEY `idx_org_msg_id` (`org_msg_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='私有消息记录表';;


CREATE TABLE `users` (
 `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
 `mobile` varchar(11) NOT NULL DEFAULT '' COMMENT '手机号',
 `nickname` varchar(64) NOT NULL DEFAULT '' COMMENT '用户昵称',
 `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '用户头像',
 `gender` tinyint unsigned NOT NULL DEFAULT '3' COMMENT '用户性别[1:男 ;2:女;3:未知]',
 `password` varchar(255) NOT NULL COMMENT '用户密码',
 `motto` varchar(500) NOT NULL DEFAULT '' COMMENT '用户座右铭',
 `email` varchar(30) NOT NULL DEFAULT '' COMMENT '用户邮箱',
 `birthday` varchar(10) NOT NULL DEFAULT '' COMMENT '生日',
 `status` int NOT NULL DEFAULT '1' COMMENT '用户状态[1:正常;2:停用;3:注销]',
 `is_robot` tinyint unsigned NOT NULL DEFAULT '2' COMMENT '是否机器人[1:是;2:否;]',
 `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
 `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
 PRIMARY KEY (`id`) USING BTREE,
 UNIQUE KEY `uk_mobile` (`mobile`) USING BTREE,
 KEY `idx_created_at` (`created_at`) USING BTREE,
 KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4531 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='用户表';;

CREATE TABLE IF NOT EXISTS `users_emoticon`
(
    `id`           int unsigned NOT NULL AUTO_INCREMENT COMMENT '表情包收藏ID',
    `user_id`      int unsigned NOT NULL COMMENT '用户ID',
    `emoticon_ids` json         NOT NULL COMMENT '表情包ID',
    `created_at`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id` (`user_id`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户收藏表情包';;



CREATE TABLE IF NOT EXISTS `article_history`
(
    `id`         int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id`    int(11) unsigned NOT NULL COMMENT '用户ID',
    `article_id` int(11) unsigned NOT NULL COMMENT '笔记ID',
    `content`    longtext NOT NULL COMMENT 'markdown 内容',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY          `idx_created_at` (`created_at`) USING BTREE,
    KEY          `idx_article_id` (`article_id`) USING BTREE,
    KEY          `idx_user_id_article_id` (`user_id`,`article_id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COMMENT='笔记历史记录表';;
