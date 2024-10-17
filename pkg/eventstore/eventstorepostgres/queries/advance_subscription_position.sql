UPDATE
    subscriptions
SET
    position = position + 1
WHERE
    id = @subscription_id;
