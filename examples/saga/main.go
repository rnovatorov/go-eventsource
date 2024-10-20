package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"

	"github.com/rnovatorov/go-eventsource/examples/saga/orchestratorservice"
	"github.com/rnovatorov/go-eventsource/examples/saga/postgresadapter"
	"github.com/rnovatorov/go-eventsource/examples/saga/sagapb"
	"github.com/rnovatorov/go-eventsource/examples/saga/workerservice"
	"github.com/rnovatorov/go-eventsource/pkg/eventstore/eventstorepostgres"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("new database pool: %w", err)
	}
	defer pool.Close()

	eventStore := eventstorepostgres.Start(pool,
		eventstorepostgres.WithLogger(logger))
	defer eventStore.Stop()

	requestSender, err := postgresadapter.StartRequestSender(ctx,
		postgresadapter.RequestSenderParams{
			Pool: pool,
		})
	if err != nil {
		return fmt.Errorf("start postgres request sender: %w", err)
	}
	defer requestSender.Stop()

	orchestrator := orchestratorservice.New(orchestratorservice.Params{
		EventStore:    eventStore,
		RequestSender: requestSender,
	})

	if err := eventStore.Subscribe(
		ctx, "sub_id", orchestrator.HandleSagaEvent,
	); err != nil {
		return fmt.Errorf("subscribe saga event handler: %w", err)
	}

	responseReceiver, err := postgresadapter.StartResponseReceiver(ctx,
		postgresadapter.ResponseReceiverParams{
			Pool:            pool,
			RequestHandler:  workerservice.New(),
			ResponseHandler: orchestrator,
		})
	if err != nil {
		return fmt.Errorf("start postgres response receiver: %w", err)
	}
	defer responseReceiver.Stop()

	tasks, err := loadExampleTasksDefinition()
	if err != nil {
		return fmt.Errorf("load example tasks definition: %w", err)
	}

	id := fmt.Sprint(time.Now().Unix())

	if err := orchestrator.BeginSaga(ctx, id, tasks); err != nil {
		return fmt.Errorf("begin saga: %w", err)
	}

	//if err := orchestrator.BeginSaga(ctx, id, []*sagapb.TaskDefinition{
	//	{Id: "task1"},
	//	{Id: "task2", Arguments: []*sagapb.ArgumentDefinition{
	//		{Name: "arg1", Value: &sagapb.ArgumentDefinition_StaticValue{
	//			StaticValue: structpb.NewNumberValue(10),
	//		}},
	//	}, Dependencies: []string{"task1"}},
	//	{Id: "task3", Arguments: []*sagapb.ArgumentDefinition{
	//		{Name: "arg2", Value: &sagapb.ArgumentDefinition_ResultReference{
	//			ResultReference: "task1",
	//		}},
	//	}, Dependencies: []string{"task1"}},
	//	{Id: "task4", Dependencies: []string{"task2", "task3"}},
	//}); err != nil {
	//	return fmt.Errorf("begin saga: %w", err)
	//}

	<-ctx.Done()
	return nil
}

func loadExampleTasksDefinition() (defs []*sagapb.TaskDefinition, _ error) {
	var payloads []interface{}
	if err := yaml.Unmarshal(exampleTasksDefinitionYml, &payloads); err != nil {
		return nil, fmt.Errorf("unmarshal YAML: %w", err)
	}

	for _, p := range payloads {
		jsonData, err := json.Marshal(p)
		if err != nil {
			return nil, fmt.Errorf("marshal JSON: %w", err)
		}
		var def sagapb.TaskDefinition
		if err := protojson.Unmarshal(jsonData, &def); err != nil {
			return nil, fmt.Errorf("unmarshal JSON: %w", err)
		}
		defs = append(defs, &def)
	}

	return defs, nil
}

//go:embed example_tasks_definition.yml
var exampleTasksDefinitionYml []byte
