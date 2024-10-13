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
)
