create table settings
(
    name            varchar(64) not null,
    value           varchar(64) not null,
    next_value      varchar(64) null,
    activation_time timestamp null,
    primary key (name)
) ENGINE = InnoDB
  default CHARSET = utf8;
