-- Table: public.expressions

-- DROP TABLE IF EXISTS public.expressions;

CREATE TABLE IF NOT EXISTS public.expressions
(
    username character varying(50) COLLATE pg_catalog."default" NOT NULL,
    expression character varying(100) COLLATE pg_catalog."default" NOT NULL,
    worker character varying(30) COLLATE pg_catalog."default" NOT NULL,
    "startTime" time without time zone NOT NULL,
    "endTime" time without time zone
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.expressions
    OWNER to postgres;