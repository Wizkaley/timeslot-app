CREATE TABLE public.time_slots
(
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    time_slot character varying NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT user_id_foreign_key FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);