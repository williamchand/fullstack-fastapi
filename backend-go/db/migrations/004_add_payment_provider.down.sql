ALTER TABLE public.payment DROP COLUMN provider;
DROP INDEX IF EXISTS ix_payment_provider_tx;
