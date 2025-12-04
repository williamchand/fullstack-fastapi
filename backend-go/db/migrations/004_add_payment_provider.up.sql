ALTER TABLE public.payment ADD COLUMN provider varchar(20) NOT NULL DEFAULT 'stripe';
CREATE INDEX IF NOT EXISTS ix_payment_provider_tx ON public.payment (provider, transaction_id);
