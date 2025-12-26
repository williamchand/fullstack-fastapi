-- Revert user_id to NOT NULL
ALTER TABLE public.verification_code ALTER COLUMN user_id SET NOT NULL;