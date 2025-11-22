CREATE TYPE public.rsvp_status AS ENUM ('yes', 'no', 'maybe');

CREATE TYPE public.payment_status AS ENUM ('pending', 'paid', 'failed');

CREATE TYPE public.wedding_status AS ENUM ('draft', 'active', 'archived');

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

CREATE TABLE public."template" (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	"name" varchar(100) NOT NULL,
	theme_config jsonb DEFAULT '{}'::jsonb NOT NULL,
	config_schema jsonb DEFAULT '{}'::jsonb NOT NULL,
	preview_url varchar(512) NULL,
	price numeric(10, 2) DEFAULT 0.00 NOT NULL,
	created_at timestamptz DEFAULT now() NULL,
	CONSTRAINT template_name_key UNIQUE (name),
	CONSTRAINT template_pkey PRIMARY KEY (id)
);

CREATE TABLE public.guest (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	wedding_id uuid NOT NULL,
	"name" varchar(255) NOT NULL,
	contact varchar(255) NOT NULL,
	"rsvp_status" public."rsvp_status" DEFAULT 'maybe'::rsvp_status NOT NULL,
	message text NULL,
	created_at timestamptz DEFAULT now() NULL,
	deleted_at timestamptz NULL,
	CONSTRAINT guest_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX uix_guest_wedding_contact_active ON public.guest USING btree (wedding_id, contact) WHERE (deleted_at IS NULL);

CREATE TABLE public.item (
	description varchar(255) NOT NULL,
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	title varchar(255) NOT NULL,
	owner_id uuid NOT NULL,
	deleted_at timestamptz NULL,
	CONSTRAINT item_pkey PRIMARY KEY (id)
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
	created_at timestamptz DEFAULT now() NULL,
	CONSTRAINT payment_pkey PRIMARY KEY (id),
	CONSTRAINT payment_transaction_id_key UNIQUE (transaction_id)
);

CREATE TABLE public.wedding (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	user_id uuid NOT NULL,
	template_id uuid NULL,
	payment_id uuid NULL,
	status public."wedding_status" DEFAULT 'draft'::wedding_status NOT NULL,
	custom_domain varchar(255) NULL,
	slug varchar(150) NULL,
	config_data jsonb DEFAULT '{}'::jsonb NOT NULL,
	created_at timestamptz DEFAULT now() NULL,
	deleted_at timestamptz NULL,
	CONSTRAINT wedding_custom_domain_key UNIQUE (custom_domain),
	CONSTRAINT wedding_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX ix_wedding_slug_active ON public.wedding USING btree (slug) WHERE (deleted_at IS NULL);

ALTER TABLE public.guest ADD CONSTRAINT guest_wedding_id_fkey FOREIGN KEY (wedding_id) REFERENCES public.wedding(id) ON DELETE CASCADE;

ALTER TABLE public.item ADD CONSTRAINT item_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public."user"(id) ON DELETE CASCADE;

ALTER TABLE public.payment ADD CONSTRAINT payment_payment_method_id_fkey FOREIGN KEY (payment_method_id) REFERENCES public.payment_method(id) ON DELETE SET NULL;
ALTER TABLE public.payment ADD CONSTRAINT payment_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;

ALTER TABLE public.wedding ADD CONSTRAINT wedding_payment_id_fkey FOREIGN KEY (payment_id) REFERENCES public.payment(id) ON DELETE SET NULL;
ALTER TABLE public.wedding ADD CONSTRAINT wedding_template_id_fkey FOREIGN KEY (template_id) REFERENCES public."template"(id) ON DELETE SET NULL;
ALTER TABLE public.wedding ADD CONSTRAINT wedding_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;