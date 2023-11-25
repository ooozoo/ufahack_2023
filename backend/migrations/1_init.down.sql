drop table if exists topic_results cascade;

drop table if exists answers cascade;

drop table if exists questions cascade;

drop table if exists attempts cascade;

drop table if exists users cascade;

drop table if exists topics cascade;

drop table if exists subjects cascade;

drop function if exists update_lastup_col() cascade;

drop function if exists deleted_at() cascade;

drop procedure if exists set_admin(varchar, boolean) cascade;

drop procedure if exists set_admin(uuid, boolean) cascade;

drop procedure if exists delete_user(varchar) cascade;

drop procedure if exists delete_user(uuid) cascade;

drop function if exists insert_user(varchar, varchar) cascade;

drop function if exists insert_subject(varchar) cascade;

drop procedure if exists delete_subject(uuid) cascade;

drop function if exists insert_topic(uuid, varchar) cascade;

drop procedure if exists delete_topic(uuid) cascade;

drop function if exists is_admin(uuid) cascade;

drop function if exists insert_question(uuid, text) cascade;

drop procedure if exists delete_question(uuid) cascade;

drop function if exists insert_answer(uuid, text, boolean) cascade;

drop procedure if exists delete_answer(uuid) cascade;

drop function if exists insert_attempt(uuid, uuid) cascade;

drop procedure if exists confirm_attempt(uuid) cascade;

drop procedure if exists insert_topic_res(uuid, uuid, uuid, boolean) cascade;

drop function if exists get_topics_stat(uuid, uuid) cascade;

drop function if exists get_topics(uuid) cascade;

drop function if exists get_topic_test(uuid) cascade;

drop function if exists get_user_by_username(varchar) cascade;

drop function if exists get_subjects() cascade;

drop function if exists get_subject_by_id(uuid) cascade;

