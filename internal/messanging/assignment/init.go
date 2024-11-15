package assignment

import (
	"errors"
	"github.com/upassed/upassed-assignment-service/internal/logging"
	"github.com/wagslane/go-rabbitmq"
)

var (
	errCreatingAssignmentCreateQueueConsumer = errors.New("unable to create assignment queue consumer")
	errRunningAssignmentCreateQueueConsumer  = errors.New("unable to run assignment queue consumer")
)

func InitializeCreateQueueConsumer(client *rabbitClient) error {
	log := logging.Wrap(client.log,
		logging.WithOp(InitializeCreateQueueConsumer),
	)

	log.Info("started crating assignment create queue consumer")
	assignmentCreateGroupConsumer, err := rabbitmq.NewConsumer(
		client.rabbitConnection,
		client.cfg.Rabbit.Queues.AssignmentCreate.Name,
		rabbitmq.WithConsumerOptionsRoutingKey(client.cfg.Rabbit.Queues.AssignmentCreate.RoutingKey),
		rabbitmq.WithConsumerOptionsExchangeName(client.cfg.Rabbit.Exchange.Name),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)

	if err != nil {
		log.Error("unable to create assignment queue consumer", logging.Error(err))
		return errCreatingAssignmentCreateQueueConsumer
	}

	defer assignmentCreateGroupConsumer.Close()
	if err := assignmentCreateGroupConsumer.Run(client.CreateQueueConsumer()); err != nil {
		log.Error("unable to run assignment queue consumer")
		return errRunningAssignmentCreateQueueConsumer
	}

	log.Info("assignment queue consumer successfully initialized")
	return nil
}
