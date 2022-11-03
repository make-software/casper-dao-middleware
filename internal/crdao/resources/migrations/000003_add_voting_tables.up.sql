create table votings
(
    creator            binary(32) not null,
    deploy_hash        binary(32) not null,
    voting_id          int unsigned not null,
    informal_voting_id int unsigned null,
    is_formal          tinyint unsigned not null,
    has_ended          tinyint unsigned not null,
    voting_quorum      int unsigned not null,
    voting_time        int unsigned not null,
    timestamp          datetime not null,

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
    timestamp    datetime not null,

    primary key (address, voting_id, deploy_hash)
) ENGINE = InnoDB
  default CHARSET = utf8;

