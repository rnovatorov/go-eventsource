package eventstorepostgres

import _ "embed"

var (
	//go:embed queries/list_events.sql
	listEventsQuery string

	//go:embed queries/create_aggregate.sql
	createAggregateQuery string

	//go:embed queries/update_aggregate_version.sql
	updateAggregateVersionQuery string

	//go:embed queries/save_event.sql
	saveEventQuery string

	//go:embed queries/sequence_events.sql
	sequenceEventsQuery string

	//go:embed queries/notify_events_inserted.sql
	notifyEventsInsertedQuery string

	//go:embed queries/acquire_events_advisory_lock.sql
	acquireEventsAdvisoryLockQuery string

	//go:embed queries/notify_events_sequenced.sql
	notifyEventsSequencedQuery string

	//go:embed queries/create_subscription.sql
	createSubscriptionQuery string

	//go:embed queries/populate_subscription_backlog.sql
	populateSubscriptionBacklogQuery string

	//go:embed queries/select_subscription_event_for_processing.sql
	selectSubscriptionEventForProcessingQuery string

	//go:embed queries/complete_subscription_event_processing.sql
	completeSubscriptionEventProcessingQuery string
)
