-- name: UpsertSubscription :one
INSERT INTO subscription (user_id, stripe_customer_id, stripe_subscription_id, status, current_period_end)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id)
DO UPDATE SET stripe_customer_id = EXCLUDED.stripe_customer_id,
              stripe_subscription_id = EXCLUDED.stripe_subscription_id,
              status = EXCLUDED.status,
              current_period_end = EXCLUDED.current_period_end,
              updated_at = now()
RETURNING *;

-- name: GetSubscriptionByUser :one
SELECT * FROM subscription WHERE user_id = $1 LIMIT 1;

-- name: ListAllSubscriptions :many
SELECT * FROM subscription;
