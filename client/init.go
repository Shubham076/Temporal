package client

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/worker"
	"log"
	logrusadapter "logur.dev/adapter/logrus"
	"logur.dev/logur"
	"sync"
)

var temporalClient client.Client
var QUEUE_NAME = "payment"

func Init() {
	var err error
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.InfoLevel)
	logger := logur.LoggerToKV(logrusadapter.New(logrusLogger))

	meter := otel.GetMeterProvider().Meter("app")
	metricsHandler := opentelemetry.NewMetricsHandler(opentelemetry.MetricsHandlerOptions{
		Meter: meter,
	})
	temporalClient, err = client.Dial(client.Options{
		HostPort:       "127.0.0.1:7233",
		Namespace:      "default",
		MetricsHandler: metricsHandler,
		Logger:         logger,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal Client", err)
	}
}

func StartWorkflow(message interface{}, workflow interface{}, options client.StartWorkflowOptions) error {
	_, err := temporalClient.ExecuteWorkflow(context.Background(), options, workflow, message)
	if err != nil {
		return err
	}
	return nil
}

func StartWorkers(wg *sync.WaitGroup, count int, queue string, options worker.Options, workflows []interface{}, activities []interface{}) {
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			w := worker.New(temporalClient, queue, options)
			for _, workflow := range workflows {
				w.RegisterWorkflow(workflow)
			}
			for _, activity := range activities {
				w.RegisterActivity(activity)
			}
			err := w.Run(worker.InterruptCh())
			if err != nil {
				log.Fatalln("Unable to start worker", err)
			}
		}()
	}
}
