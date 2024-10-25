DELETE FROM es_subscription_backlogs
WHERE subscription_id = @subscription_id
    AND event_id = @event_id;
