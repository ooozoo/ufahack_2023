create or replace procedure insert_user(IN user_uid uuid, IN name character varying, IN password character varying)
    language plpgsql
as
$$
BEGIN
    INSERT INTO users (users_uid, user_name, pass) VALUES (user_uid, name, password);
END;
$$;