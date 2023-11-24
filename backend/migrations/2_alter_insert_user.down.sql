create or replace procedure insert_user(user_uid uuid, name character varying, password character varying)
    language plpgsql
as
$$
BEGIN
    INSERT INTO users (users_uid, user_name, pass, is_admin) VALUES (user_uid, name, password);
END;
$$;