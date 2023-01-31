CREATE TABLE users (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    firstname text,
    lastname text,
    username character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    is_verified boolean,
    profile_path text,
    customer_id text,
    subscription_id text,
    token text,
    secret text,
    tow_fa boolean,
    two_fa boolean,
    country text,
    timezone text,
    price_id text,
    active boolean DEFAULT false
);



ALTER TABLE ONLY plans
    ADD CONSTRAINT plans_name_key UNIQUE (name);


ALTER TABLE ONLY plans
    ADD CONSTRAINT plans_pkey PRIMARY KEY (id);


ALTER TABLE ONLY prices
    ADD CONSTRAINT prices_pkey PRIMARY KEY (id);


ALTER TABLE ONLY prices
    ADD CONSTRAINT prices_price_id_key UNIQUE (price_id);


ALTER TABLE ONLY reset_passwords
    ADD CONSTRAINT reset_passwords_pkey PRIMARY KEY (id);


ALTER TABLE ONLY public.storage_types
    ADD CONSTRAINT storage_types_name_key UNIQUE (name);


ALTER TABLE ONLY storage_types
    ADD CONSTRAINT storage_types_pkey PRIMARY KEY (id);


ALTER TABLE ONLY storages
    ADD CONSTRAINT storages_pkey PRIMARY KEY (id);


ALTER TABLE ONLY users
    ADD CONSTRAINT users_email_key UNIQUE (email);


ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


ALTER TABLE ONLY users
    ADD CONSTRAINT users_username_key UNIQUE (username);
