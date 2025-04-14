CREATE DATABASE IF NOT EXISTS gogateway;

USE gogateway;

create table go_gateway_admin
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3) null,
    updated_at datetime(3) null,
    deleted_at datetime(3) null,
    username   longtext    null,
    salt       longtext    null,
    password   longtext    null
);

create index idx_go_gateway_admin_deleted_at
    on go_gateway_admin (deleted_at);

create table go_gateway_application
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)  null,
    updated_at datetime(3)  null,
    deleted_at datetime(3)  null,
    app_id     varchar(128) not null,
    name       varchar(255) not null,
    secret     longtext     not null,
    white_ip_s longtext     null,
    qpd        bigint       null,
    qps        bigint       null
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

create index idx_go_gateway_application_deleted_at
    on go_gateway_application (deleted_at);

create table go_gateway_service_access_control
(
    id                   bigint unsigned auto_increment
        primary key,
    service_id           bigint unsigned null,
    open_auth            tinyint         null,
    black_list           varchar(600)    null,
    white_list           varchar(600)    null,
    white_host_name      varchar(255)    null,
    client_ip_flow_limit bigint          null,
    service_flow_limit   bigint          null
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

create table go_gateway_service_grpc_rule
(
    id              bigint unsigned auto_increment
        primary key,
    service_id      bigint unsigned null,
    port            bigint          null,
    header_transfer longtext        null
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

create table go_gateway_service_http_rule
(
    id              bigint unsigned auto_increment
        primary key,
    service_id      bigint unsigned null,
    rule_type       int default 0   null,
    rule            varchar(300)    null,
    need_https      tinyint         null,
    need_websocket  tinyint         null,
    need_strip_uri  tinyint         null,
    url_rewrite     varchar(500)    null,
    header_transfer varchar(500)    null
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

create table go_gateway_service_info
(
    id           bigint unsigned auto_increment
        primary key,
    created_at   datetime(3)  null,
    updated_at   datetime(3)  null,
    deleted_at   datetime(3)  null,
    load_type    bigint       null,
    service_name varchar(256) null,
    service_desc varchar(256) null
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

create index idx_go_gateway_service_info_deleted_at
    on go_gateway_service_info (deleted_at);

create table go_gateway_service_load_balance
(
    id                       bigint unsigned auto_increment
        primary key,
    service_id               bigint unsigned null,
    check_method             bigint          null,
    check_timeout            bigint          null,
    check_interval           bigint          null,
    round_type               bigint          null,
    ip_list                  longtext        null,
    weight_list              longtext        null,
    forbid_list              longtext        null,
    upstream_connect_timeout bigint          null,
    upstream_header_timeout  bigint          null,
    upstream_idle_timeout    bigint          null,
    upstream_max_idle        bigint          null
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

create table go_gateway_service_tcp_rule
(
    id         bigint unsigned auto_increment
        primary key,
    service_id bigint unsigned null,
    port       bigint          null
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# 初始密码
INSERT INTO go_gateway_admin VALUES (1,
                                     current_time(),
                                     current_time(),
                                     null,
                                     'admin',
                                     '123456',
                                     '9b1063951d443cfac15cc879efb4054f4f4fd599e1b1a9aee67b0301e19e40fe')

