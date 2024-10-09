package sagaservice

type TaskExecutor interface {
	BeginTask(ctx SagaContext, id string) error
	BeginCompensation(ctx SagaContext, id string) error
}
