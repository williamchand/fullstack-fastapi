-- Enhanced schema with better constraints and indexing

CREATE TABLE public.role (
    id serial4 NOT NULL,
    name varchar(50) NOT NULL,
    description varchar(255) NULL,
    created_at timestamptz DEFAULT now() NULL,
    updated_at timestamptz DEFAULT now() NULL,
    CONSTRAINT role_pkey PRIMARY KEY (id),
    CONSTRAINT role_name_key UNIQUE (name)
);

CREATE TABLE public."user" (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    email varchar(255) NOT NULL,
    phone_number varchar(32) NULL,
    full_name varchar(255) NULL,
    hashed_password varchar NULL,
    is_active bool DEFAULT true NOT NULL,
    is_email_verified bool DEFAULT false NOT NULL,
    is_phone_verified bool DEFAULT false NOT NULL,
    is_totp_enabled bool DEFAULT false NOT NULL,
    totp_secret varchar NULL,
    created_at timestamptz DEFAULT now() NULL,
    updated_at timestamptz DEFAULT now() NULL,
    CONSTRAINT user_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX ix_user_email ON public."user" USING btree (email) WHERE (is_active = true);
CREATE UNIQUE INDEX ix_user_phone ON public."user" USING btree (phone_number) WHERE (is_active = true AND phone_number IS NOT NULL);

CREATE TABLE public.user_role (
    user_id uuid NOT NULL,
    role_id int4 NOT NULL,
    assigned_at timestamptz DEFAULT now() NULL,
    CONSTRAINT user_role_pkey PRIMARY KEY (user_id, role_id),
    CONSTRAINT user_role_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.role(id) ON DELETE CASCADE,
    CONSTRAINT user_role_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE
);

CREATE TABLE public.oauth_account (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    provider varchar(50) NOT NULL,
    provider_user_id varchar(255) NOT NULL,
    access_token varchar NULL,
    refresh_token varchar NULL,
    token_expires_at timestamptz NULL,
    created_at timestamptz DEFAULT now() NULL,
    updated_at timestamptz DEFAULT now() NULL,
	provider_data jsonb DEFAULT '{}'::jsonb NOT NULL,
    CONSTRAINT oauth_account_pkey PRIMARY KEY (id),
    CONSTRAINT oauth_account_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE,
    CONSTRAINT uq_oauth_provider_user UNIQUE (provider, provider_user_id)
);

CREATE TABLE public.verification_code (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    verification_code varchar(64) NOT NULL,
    verification_type varchar(16) NOT NULL,
    extra_metadata jsonb DEFAULT '{}'::jsonb NOT NULL,
    created_at timestamptz DEFAULT now() NULL,
    expires_at timestamptz NOT NULL,
    used_at timestamptz NULL,
    CONSTRAINT verification_code_pkey PRIMARY KEY (id),
    CONSTRAINT verification_code_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE
);

CREATE INDEX idx_verification_code_user_type ON public.verification_code (user_id, verification_type);
CREATE INDEX idx_verification_code_expires ON public.verification_code (expires_at) WHERE used_at IS NULL;
-- Index to speed lookups by verification type + code for unused verification codes
CREATE INDEX idx_verification_code_type_code ON public.verification_code (verification_type, verification_code) WHERE used_at IS NULL;

INSERT INTO public.role (name, description) VALUES
('salon_owner','Beauty salon owner'),
('salon_employee','salon employee role'),
('customer','Default customer role'),
('superuser','Superuser role');

CREATE TABLE IF NOT EXISTS email_template (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name TEXT NOT NULL UNIQUE,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_email_template_name ON public.email_template USING btree (name) WHERE (is_active = true);

INSERT INTO email_template (name, subject, body)
VALUES 
(
  'verification_email',
  'Verify Your Email Address',
  '<p>Hello,</p><p>Your verification code is: <strong>{{.code}}</strong></p>'
),
(
  'password_reset',
  'Reset Your Password',
  '<p>Hello,</p><p>Click the link to reset your password: <a href="{{.link}}">Reset Password</a></p>'
),
(
  'verification_phone',
  'Verify Your Phone Number',
  'Your verification code is {{.code}}'
)
ON CONFLICT (name) DO NOTHING;