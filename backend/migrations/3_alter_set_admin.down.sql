drop procedure if exists set_admin(uid uuid, admin boolean);

create or replace procedure set_admin(name character varying, admin boolean)
    language plpgsql
as
$$
BEGIN
    UPDATE users SET is_admin = admin WHERE user_name = name;
END;
$$;
