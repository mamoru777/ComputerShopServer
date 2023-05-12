-- +migrate Up
create table usr
(
    ID int Primary key,
    Login char(15) not null,
    Password char(20) not null,
    Name char(15) null,
    LastName char(20) null,
    SurName char(20) null,
    Email char(20) null,
    Avatar bytea null
);

-- +migare Down
drop table usr;