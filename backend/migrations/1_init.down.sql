drop trigger if exists update_lastup on users;
drop trigger if exists update_lastup on subjects;
drop trigger if exists update_lastup on topics;
drop trigger if exists update_lastup on questions;
drop trigger if exists update_lastup on answers;
drop trigger if exists update_lastup on attempts;
drop trigger if exists update_lastup on topic_results;

drop function if exists update_lastup_col();
drop procedure if exists insert_user(user_uid uuid, name varchar, password varchar);
drop procedure if exists set_admin(name varchar, admin boolean);

drop table if exists topic_results;
drop table if exists attempts;
drop table if exists answers;
drop table if exists questions;
drop table if exists topics;
drop table if exists subjects;
drop table if exists users;