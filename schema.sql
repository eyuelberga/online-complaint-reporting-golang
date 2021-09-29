
CREATE FUNCTION public.update_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$;

CREATE TABLE public.complaints (
    id integer NOT NULL,
    email text NOT NULL,
    file text,
    comment text NOT NULL,
    fullname text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE SEQUENCE public.complaints_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.complaints_id_seq OWNED BY public.complaints.id;


CREATE TABLE public.users (
    email text NOT NULL,
    password text NOT NULL,
    role text NOT NULL,
    fullname text NOT NULL,
    enabled boolean NOT NULL
);


ALTER TABLE ONLY public.complaints ALTER COLUMN id SET DEFAULT nextval('public.complaints_id_seq'::regclass);


ALTER TABLE ONLY public.complaints
    ADD CONSTRAINT complaints_pkey PRIMARY KEY (id);



ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (email);


CREATE TRIGGER update_complaints_time BEFORE UPDATE ON public.complaints FOR EACH ROW EXECUTE PROCEDURE public.update_column();


