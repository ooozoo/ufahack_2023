drop procedure if exists insert_user(uuid, varchar, varchar);

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