create table if not exists users
(
    users_uid  uuid                    not null
        constraint users_pk
            primary key,
    user_name  varchar(255)            not null
        constraint users_pk2
            unique,
    pass       varchar(255)            not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    is_admin   boolean   default false not null,
    is_deleted boolean   default false not null,
    deleted_at timestamp
);

create table if not exists subjects
(
    subjects_uid uuid                    not null
        constraint subjects_pk
            primary key,
    subject_name varchar(255)            not null
        constraint subjects_pk2
            unique,
    created_at   timestamp default now() not null,
    updated_at   timestamp default now() not null,
    is_deleted   boolean   default false not null,
    deleted_at   timestamp
);

create table if not exists topics
(
    topics_uid  uuid                    not null
        constraint topics_pk
            primary key,
    subject_uid uuid                    not null
        constraint topics_subjects_subjects_uid_fk
            references subjects,
    topic_name  varchar(255)            not null,
    created_at  timestamp default now() not null,
    updated_at  timestamp default now() not null,
    is_deleted  boolean   default false not null,
    deleted_at  timestamp
);

create table if not exists questions
(
    questions_uid uuid                    not null
        constraint questions_pk
            primary key,
    topic_uid     uuid                    not null
        constraint questions_themes_themes_uid_fk
            references topics,
    text          text                    not null,
    created_at    timestamp default now() not null,
    updated_at    timestamp default now() not null,
    is_deleted    boolean   default false not null,
    deleted_at    timestamp
);

create table if not exists answers
(
    answers_uid  uuid                    not null
        constraint answers_pk
            primary key,
    question_uid uuid                    not null
        constraint answers_questions_questions_uid_fk
            references questions,
    text         text                    not null,
    is_valid     boolean                 not null,
    created_at   timestamp default now() not null,
    updated_at   timestamp default now() not null,
    is_deleted   boolean   default false not null,
    deleted_at   timestamp
);

create table if not exists attempts
(
    attempts_uid uuid                    not null
        constraint attempts_pk
            primary key,
    topic_uid    uuid                    not null
        constraint attempts_topics_topics_uid_fk
            references topics,
    user_uid     uuid                    not null
        constraint attempts_users_users_uid_fk
            references users,
    created_at   timestamp default now() not null,
    updated_at   timestamp default now() not null,
    is_deleted   boolean   default false not null,
    is_confirmed boolean   default false not null,
    deleted_at   timestamp
);

create table if not exists topic_results
(
    attempt_uid  uuid                    not null
        constraint topic_results_attempts_attempts_uid_fk
            references attempts,
    question_uid uuid                    not null
        constraint topic_results_questions_questions_uid_fk
            references questions,
    answer_uid   uuid                    not null
        constraint topic_results_answers_answers_uid_fk
            references answers,
    created_at   timestamp default now() not null,
    updated_at   timestamp default now() not null,
    is_valid     boolean                 not null
);

create or replace function update_lastup_col() returns trigger
    language plpgsql
as
$$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;

create trigger update_lastup
    before update
    on users
    for each row
execute procedure update_lastup_col();

create trigger update_lastup
    before update
    on subjects
    for each row
execute procedure update_lastup_col();

create trigger update_lastup
    before update
    on topics
    for each row
execute procedure update_lastup_col();

create trigger update_lastup
    before update
    on questions
    for each row
execute procedure update_lastup_col();

create trigger update_lastup
    before update
    on answers
    for each row
execute procedure update_lastup_col();

create trigger update_lastup
    before update
    on attempts
    for each row
execute procedure update_lastup_col();

create trigger update_lastup
    before update
    on topic_results
    for each row
execute procedure update_lastup_col();

create or replace function deleted_at() returns trigger
    language plpgsql
as
$$
BEGIN
    NEW.deleted_at = now();
    RETURN NEW;
END;
$$;

create trigger deleted_at
    before update
        of is_deleted
    on users
    for each row
execute procedure deleted_at();

create trigger deleted_at
    before update
        of is_deleted
    on subjects
    for each row
execute procedure deleted_at();

create trigger deleted_at
    before update
        of is_deleted
    on topics
    for each row
execute procedure deleted_at();

create trigger deleted_at
    before update
        of is_deleted
    on questions
    for each row
execute procedure deleted_at();

create trigger deleted_at
    before update
        of is_deleted
    on answers
    for each row
execute procedure deleted_at();

create trigger deleted_at
    before update
        of is_deleted
    on attempts
    for each row
execute procedure deleted_at();

create or replace procedure set_admin(IN name character varying, IN admin boolean)
    language plpgsql
as
$$
BEGIN
    UPDATE users SET is_admin = admin WHERE user_name = name;
END;
$$;

create or replace procedure set_admin(IN uid uuid, IN admin boolean)
    language plpgsql
as
$$
BEGIN
    UPDATE users SET is_admin = admin WHERE users_uid = uid;
END;
$$;

create or replace procedure delete_user(IN user_name character varying)
    language plpgsql
as
$$
begin
    update users set is_deleted = true where user_name = delete_user.user_name;
end;
$$;

create or replace procedure delete_user(IN user_uid uuid)
    language plpgsql
as
$$
begin
    update users set is_deleted = true where users_uid = delete_user.user_uid;
end;
$$;

create or replace function insert_user(name character varying, password character varying) returns uuid
    language plpgsql
as
$$
DECLARE
    user_uid uuid := gen_random_uuid();
BEGIN
    INSERT INTO users (users_uid, user_name, pass) VALUES (user_uid, name, password);
    RETURN user_uid;
END;
$$;

create or replace function insert_subject(subject_name character varying) returns uuid
    language plpgsql
as
$$
declare
    uid uuid := gen_random_uuid();
begin
    insert into subjects (subjects_uid, subject_name) values (uid, insert_subject.subject_name);
    return uid;
end;
$$;

create or replace procedure delete_subject(IN uid uuid)
    language plpgsql
as
$$
begin
    update subjects set is_deleted = true where subjects_uid = uid;
    update topics set is_deleted = true where subject_uid = uid;
    update questions set is_deleted = true where topic_uid in (select topics_uid from topics where subject_uid = uid);
    update answers set is_deleted = true where question_uid in(select questions_uid from questions where topic_uid in (select topics_uid from topics where subject_uid = uid));
end;
$$;

create or replace function insert_topic(subject_uid uuid, topic_name character varying) returns uuid
    language plpgsql
as
$$
declare
    uid uuid := gen_random_uuid();
begin
    insert into topics (topics_uid, subject_uid, topic_name) VALUES (uid, insert_topic.subject_uid, insert_topic.topic_name);
    return uid;
end;
$$;

create or replace procedure delete_topic(IN topic_uid uuid)
    language plpgsql
as
$$
begin
    update topics set is_deleted = true where topics_uid = topic_uid;
    update questions set is_deleted = true where topic_uid = delete_topic.topic_uid;
    update answers set is_deleted = true where question_uid in (select questions_uid from questions where questions.topic_uid = delete_topic.topic_uid);
end;
$$;

create or replace function is_admin(uid uuid) returns boolean
    language sql
as
$$
select is_admin from users where users_uid = uid limit 1;
$$;

create or replace function insert_question(topic_uid uuid, text text) returns uuid
    language plpgsql
as
$$
declare
    uid uuid := gen_random_uuid();
begin
    insert into questions (questions_uid, topic_uid, text) values (uid, insert_question.topic_uid, insert_question.text);
    return uid;
end;
$$;

create or replace procedure delete_question(IN question_uid uuid)
    language plpgsql
as
$$
begin
    update questions set is_deleted = true where questions_uid = question_uid;
    update answers set is_deleted = true where question_uid = delete_question.question_uid;
end;
$$;

create or replace function insert_answer(question_uid uuid, text text, is_val boolean) returns uuid
    language plpgsql
as
$$
declare
    uid uuid := gen_random_uuid();
begin
    insert into answers (answers_uid, question_uid, text, is_valid) VALUES (uid, insert_answer.question_uid, insert_answer.text, is_val);
    return uid;
end;
$$;

create or replace procedure delete_answer(IN answer_uid uuid)
    language plpgsql
as
$$
begin
    update answers set is_deleted = true where answers_uid = answer_uid;
end;
$$;

create or replace function insert_attempt(user_uid uuid, topic_uid uuid) returns uuid
    language plpgsql
as
$$
declare
    uid uuid := gen_random_uuid();
begin
    insert into attempts (attempts_uid, user_uid, topic_uid) values (uid, insert_attempt.user_uid, insert_attempt.topic_uid);
    return uid;
end;
$$;

create or replace procedure confirm_attempt(IN attempt_uid uuid)
    language plpgsql
as
$$
begin
    update attempts set is_confirmed = true where attempts_uid = attempt_uid;
end;
$$;

create or replace procedure insert_topic_res(IN attempt_uid uuid, IN question_uid uuid, IN answer_uid uuid, IN is_val boolean)
    language plpgsql
as
$$
begin
    insert into topic_results (attempt_uid, question_uid, answer_uid, is_valid) values (insert_topic_res.attempt_uid, insert_topic_res.question_uid, insert_topic_res.answer_uid, is_val);
end;
$$;

create or replace function get_topics_stat(subject uuid, user_uuid uuid)
    returns TABLE(topics_id uuid, topic_name character varying, created_at timestamp without time zone, updated_at timestamp without time zone, all_count integer, valid_count integer)
    language plpgsql
as
$$
begin
    select topics_uid, topic_name, topics.created_at, topics.updated_at, stat.all_count, stat.valid_count
    from topics
             left join (select user_uid, topic_uid, count(tr.attempt_uid) as all_count, count(case when tr.is_valid then 1 end) as valid_count
                        from (select *
                              from attempts a
                              where a.created_at = (select max(created_at) from attempts a2 where a.topic_uid = a2.topic_uid and a.user_uid = a2.user_uid)) last_att
                                 left join topic_results tr on last_att.attempts_uid = tr.attempt_uid
                        where user_uid = user_uuid
                        group by user_uid, topic_uid) stat on stat.topic_uid = topics.topics_uid
    where topics.is_deleted = false
      and topics.subject_uid = subject
    ;
end;
$$;

create or replace function get_topics(subject uuid)
    returns TABLE(topics_id uuid, topic_name character varying, created_at timestamp without time zone, updated_at timestamp without time zone)
    language plpgsql
as
$$
begin
    select topics_uid, topic_name, topics.created_at, topics.updated_at
    from topics
    where topics.is_deleted = false
      and topics.subject_uid = subject
    ;
end;
$$;

create or replace function get_topic_test(topic_uuid uuid)
    returns TABLE(questions_uid uuid, answers_uid uuid, is_valid boolean)
    language plpgsql
as
$$
begin
    select questions_uid, a.answers_uid, a.is_valid
    from questions q
             left join answers a on q.questions_uid = a.question_uid and a.is_deleted = false
    where q.is_deleted = false
      and q.topic_uid = topic_uuid
    ;
end;
$$;

create or replace function get_user_by_username(username character varying)
    returns TABLE(users_uid uuid, user_name character varying, pass character varying)
    language sql
as
$$
select users_uid,
       user_name,
       pass
from users
where user_name = username limit 1;
$$;

create or replace function get_subjects()
    returns TABLE(subjects_uid uuid, subject_name character varying)
    language sql
as
$$
select subjects_uid,
       subjects.subject_name
from subjects
where is_deleted = false;
$$;

create or replace function get_subject_by_id(id uuid)
    returns TABLE(subjects_uid uuid, subject_name character varying)
    language sql
as
$$
select subjects_uid, subject_name from subjects where subjects_uid = id limit 1;
$$;

