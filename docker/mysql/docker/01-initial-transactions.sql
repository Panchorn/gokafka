create table local_db.transactions
(
    id           varchar(191)  not null
        primary key,
    ref_id       longtext      null,
    from_id      longtext      null,
    to_id        longtext      null,
    status       longtext      null,
    remark       longtext      null,
    amount       double        null,
    secret_token varchar(1000) null,
    created_date datetime(3)   null,
    updated_date datetime(3)   null
);

