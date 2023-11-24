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
    is_deleted boolean   default false not null
);

create table if not exists subjects
(
    subjects_uid uuid                    not null
        constraint subjects_pk
            primary key,
    subject_name varchar(255)            not null,
    created_at   timestamp default now() not null,
    updated_at   timestamp default now() not null,
    is_deleted   boolean   default false not null
);

create table topics
(
    topics_uid  uuid                    not null
        constraint topics_pk
            primary key,
    subject_uid uuid                    not null
        constraint topics_subjects_subjects_uid_fk
            references subjects,
    theme_name  varchar(255)            not null,
    created_at  timestamp default now() not null,
    updated_at  timestamp default now() not null,
    is_deleted  boolean   default false not null
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
    is_deleted    boolean   default false not null
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
    vs_valid     boolean                 not null,
    created_at   timestamp default now() not null,
    updated_at   timestamp default now() not null,
    is_deleted   boolean   default false not null
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
    is_confirmed boolean   default false not null
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
    updated_at   timestamp default now() not null
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

create or replace trigger update_lastup
    before update
    on users
    for each row
execute procedure update_lastup_col();

create or replace trigger update_lastup
    before update
    on subjects
    for each row
execute procedure update_lastup_col();

create or replace trigger update_lastup
    before update
    on topics
    for each row
execute procedure update_lastup_col();

create or replace trigger update_lastup
    before update
    on questions
    for each row
execute procedure update_lastup_col();

create or replace trigger update_lastup
    before update
    on answers
    for each row
execute procedure update_lastup_col();

create or replace trigger update_lastup
    before update
    on attempts
    for each row
execute procedure update_lastup_col();

create or replace trigger update_lastup
    before update
    on topic_results
    for each row
execute procedure update_lastup_col();

create or replace procedure insert_user(user_uid uuid, name character varying, password character varying)
    language plpgsql
as
$$
BEGIN
    INSERT INTO users (users_uid, user_name, pass, is_admin) VALUES (user_uid, name, password);
END;
$$;


create or replace procedure set_admin(name character varying, admin boolean)
    language plpgsql
as
$$
BEGIN
    UPDATE users SET is_admin = admin WHERE user_name = name;
END;
$$;

