drop procedure if exists set_admin(name character varying, admin boolean);

create or replace procedure set_admin(IN uid uuid, IN admin boolean)
    language plpgsql
as
$$
BEGIN
    UPDATE users SET is_admin = admin WHERE users_uid == uid;
END;
$$;
