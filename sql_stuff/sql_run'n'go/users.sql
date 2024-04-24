-- Table: public.users

-- DROP TABLE IF EXISTS public.users;

CREATE TABLE IF NOT EXISTS public.users
(
    username text COLLATE pg_catalog."default" NOT NULL,
    "encPassword" "char" NOT NULL,
    "IP" character varying(45) COLLATE pg_catalog."default",
    CONSTRAINT users_pkey PRIMARY KEY (username)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;