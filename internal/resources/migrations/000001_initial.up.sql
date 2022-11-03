create table reputation_change_reasons
(
    id   tinyint unsigned primary key,
    name varchar(32) not null
) ENGINE = InnoDB
  default CHARSET = utf8;

insert into reputation_change_reasons (id, name)
values (1, 'mint'),
       (2, 'burn'),
       (3, 'vote'),
       (4, 'voting_distribution');


create table reputation_changes
(
    address               binary(32)       not null,
    contract_package_hash binary(32)       not null,
    amount                bigint           not null,
    voting_id             int unsigned     null,
    deploy_hash           binary(32)       not null,
    reason                tinyint unsigned not null,
    timestamp             datetime         not null,

    primary key (address, contract_package_hash, deploy_hash, reason)
) ENGINE = InnoDB
  default CHARSET = utf8;
