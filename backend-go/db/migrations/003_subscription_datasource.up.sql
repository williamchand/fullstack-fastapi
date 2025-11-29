CREATE TABLE public.data_source (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    name varchar(100) NOT NULL,
    type varchar(30) NOT NULL,
    host varchar(255) NOT NULL,
    port int4 NOT NULL,
    database_name varchar(255) NOT NULL,
    username varchar(255) NOT NULL,
    password_enc varchar(1024) NOT NULL,
    options jsonb DEFAULT '{}'::jsonb NOT NULL,
    created_at timestamptz DEFAULT now() NULL,
    updated_at timestamptz DEFAULT now() NULL,
    CONSTRAINT data_source_pkey PRIMARY KEY (id),
    CONSTRAINT data_source_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX uix_data_source_user_name ON public.data_source USING btree (user_id, name);

CREATE TABLE public.ai_credential (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    provider varchar(50) NOT NULL,
    api_key_enc varchar(1024) NOT NULL,
    created_at timestamptz DEFAULT now() NULL,
    updated_at timestamptz DEFAULT now() NULL,
    CONSTRAINT ai_credential_pkey PRIMARY KEY (id),
    CONSTRAINT ai_credential_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX uix_ai_credential_user_provider ON public.ai_credential USING btree (user_id, provider);

CREATE TABLE public.subscription (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    stripe_customer_id varchar(255) NULL,
    stripe_subscription_id varchar(255) NULL,
    status varchar(50) NOT NULL,
    current_period_end timestamptz NULL,
    created_at timestamptz DEFAULT now() NULL,
    updated_at timestamptz DEFAULT now() NULL,
    CONSTRAINT subscription_pkey PRIMARY KEY (id),
    CONSTRAINT subscription_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX uix_subscription_user ON public.subscription USING btree (user_id);
