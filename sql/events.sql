CREATE TABLE public.events
(
    id uuid NOT NULL,
    event_owner character varying NOT NULL,
    title character varying NOT NULL,
    event_start_time timestamp with time zone NOT NULL,
    event_end_time timestamp with time zone NOT NULL,
    participants character varying[] NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT event_owner_user_id__foreign_key FOREIGN KEY (event_owner)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);