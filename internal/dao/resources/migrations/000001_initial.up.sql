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
       (4, 'voting_distribution'),
       (5, 'voting_burn'),
       (6, 'unstake')
;

create table voting_types
(
    id   tinyint unsigned primary key,
    name varchar(32) not null
) ENGINE = InnoDB
  default CHARSET = utf8;

insert into voting_types (id, name)
values (1, 'simple'),
       (2, 'slashing'),
       (3, 'kyc'),
       (4, 'repo'),
       (5, 'reputation'),
       (6, 'onboarding'),
       (7, 'admin');


create table reputation_changes
(
    address               binary(32) not null,
    contract_package_hash binary(32) not null,
    amount                bigint   not null,
    voting_id             int unsigned null,
    deploy_hash           binary(32) not null,
    reason                tinyint unsigned not null,
    timestamp             datetime not null,

    primary key (address, contract_package_hash, deploy_hash, reason)
) ENGINE = InnoDB
  default CHARSET = utf8;


create table reputation_totals
(
    address                  binary(32) not null,
    total_liquid_reputation  bigint not null,
    total_staked_reputation  bigint null,
    voting_lost_reputation   bigint null,
    voting_earned_reputation bigint null,
    voting_id                int unsigned null,
    deploy_hash              binary(32) not null,
    reason                   tinyint unsigned not null,
    timestamp             datetime not null,

    primary key (deploy_hash)
) ENGINE = InnoDB
  default CHARSET = utf8;

create table votings
(
    creator                                        binary(32) not null,
    deploy_hash                                    binary(32) not null,
    voting_id                                      int unsigned not null,
    voting_type_id                                 tinyint unsigned not null,
    informal_voting_quorum                         int unsigned not null,
    informal_voting_starts_at                      datetime not null,
    informal_voting_ends_at                        datetime not null,
    formal_voting_quorum                           int unsigned not null,
    formal_voting_starts_at                        datetime null,
    formal_voting_ends_at                          datetime null,
    metadata                                       json     not null,
    is_canceled                                    tinyint unsigned not null,
    informal_voting_result                         tinyint unsigned null,
    formal_voting_result                           tinyint unsigned null,
    config_total_onboarded                         int unsigned not null,
    config_voting_clearness_delta                  int unsigned not null,
    config_time_between_informal_and_formal_voting int unsigned not null,

    primary key (creator, voting_id, deploy_hash)
) ENGINE = InnoDB
  default CHARSET = utf8;


create table votes
(
    deploy_hash  binary(32) not null,
    voting_id    int unsigned not null,
    address      binary(32) not null,
    amount       int unsigned not null,
    is_in_favour tinyint unsigned not null,
    is_canceled  tinyint unsigned not null,
    timestamp    datetime not null,

    primary key (address, voting_id, deploy_hash)
) ENGINE = InnoDB
  default CHARSET = utf8;


create table settings
(
    name            varchar(64) not null,
    value           varchar(64) not null,
    next_value      varchar(64) null,
    activation_time timestamp null,
    primary key (name)
) ENGINE = InnoDB
  default CHARSET = utf8;


create table accounts
(
    hash      binary(32) not null,
    is_kyc    tinyint unsigned not null,
    is_va     tinyint unsigned not null,
    timestamp datetime not null,

    primary key (hash)
) ENGINE = InnoDB
  default CHARSET = utf8;
