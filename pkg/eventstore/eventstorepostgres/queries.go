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

	//go:embed queries/select_next_subscription_event.sql
	selectNextSubscriptionEventQuery string

	//go:embed queries/advance_subscription_position.sql
	advanceSubscriptionPositionQuery string
)
