create user 'myotp_user'@'%' identified by 'mypass!';
create database myotp;

use myotp;

create table users
(
    `user_id` int not null auto_increment,
    `name` varchar(255),
    `privilege` tinyint default 0,
    primary key (`user_id`)
);

# group groups tickets
create table `groups` (
                          `group_id` int not null auto_increment,
                          `name` char,
                          `user_id` int not null,
                          primary key (`group_id`),
                          foreign key (`user_id`)
                              references users (`user_id`)
                              on delete cascade
                              on update cascade
);

# ticket
create table `ticket` (
                          `ticket_index` int not null auto_increment,
                          `id` varchar(512),
                          `token` text,
                          `group_id` int not null,
                          primary key (`ticket_index`),
                          foreign key (`group_id`)
                              references `groups`(`group_id`)
                              on delete cascade
                              on update cascade
);

grant select, delete, update, insert on myotp.* to 'myotp_user'@'%';
flush privileges;
