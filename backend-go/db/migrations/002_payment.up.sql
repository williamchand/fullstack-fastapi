CREATE TYPE public.payment_status AS ENUM ('pending', 'paid', 'failed');

CREATE TABLE public.payment_method (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	"name" varchar(100) NOT NULL,
	provider varchar(50) NOT NULL,
	config jsonb DEFAULT '{}'::jsonb NOT NULL,
	is_active bool DEFAULT true NOT NULL,
	created_at timestamptz DEFAULT now() NULL,
	CONSTRAINT payment_method_name_key UNIQUE (name),
	CONSTRAINT payment_method_pkey PRIMARY KEY (id)
);

CREATE TABLE public.payment (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	user_id uuid NOT NULL,
	payment_method_id uuid NULL,
	amount numeric(10, 2) NOT NULL,
	currency varchar(10) DEFAULT 'USD'::character varying NOT NULL,
	status public."payment_status" DEFAULT 'pending'::payment_status NOT NULL,
	transaction_id varchar(255) NOT NULL,
	extra_metadata jsonb DEFAULT '{}'::jsonb NOT NULL,
	provider varchar(20) NOT NULL DEFAULT 'stripe',
	created_at timestamptz DEFAULT now() NULL,
	CONSTRAINT payment_pkey PRIMARY KEY (id),
	CONSTRAINT payment_transaction_id_key UNIQUE (transaction_id)
);

ALTER TABLE public.payment ADD CONSTRAINT payment_payment_method_id_fkey FOREIGN KEY (payment_method_id) REFERENCES public.payment_method(id) ON DELETE SET NULL;
ALTER TABLE public.payment ADD CONSTRAINT payment_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;
CREATE INDEX IF NOT EXISTS ix_payment_provider_tx ON public.payment (provider, transaction_id);

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
