package auth

import (
	"context"
	"github.com/upassed/upassed-assignment-service/internal/config"
	"github.com/upassed/upassed-assignment-service/internal/logging"
	"github.com/upassed/upassed-assignment-service/internal/messanging"
	"github.com/upassed/upassed-assignment-service/internal/middleware/amqp"
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
)

var amqpAuthenticationRules = map[string]tokenAuthFunc{}

func (wrapper *clientImpl) AmqpMiddleware(config *config.Config, log *slog.Logger) amqp.Middleware {
	amqpAuthenticationRules[config.Rabbit.Queues.AssignmentCreate.Name] = wrapper.TeacherAccountTypeAuthenticationFunc

	return func(ctx context.Context, next messanging.HandlerWithContext) messanging.HandlerWithContext {
		return func(ctx context.Context, d rabbitmq.Delivery) (action rabbitmq.Action) {
			log = logging.Wrap(
				log,
				logging.WithOp(wrapper.AmqpMiddleware),
				logging.WithCtx(ctx),
			)

			authenticationFunc, ok := amqpAuthenticationRules[d.RoutingKey]
			if !ok {
				log.Info("authentication function is not provided, using AnyAccountTypeAuthenticationFunc")
				authenticationFunc = wrapper.AnyAccountTypeAuthenticationFunc
			}

			token, ok := d.Headers[AuthenticationHeaderKey]
			if !ok {
				log.Error("authentication header is not passed, discarding the message")
				return rabbitmq.NackDiscard
			}

			enrichedCtx, err := authenticationFunc(ctx, token.(string))
			if err != nil {
				log.Error("authentication failed, discarding the message", slog.String("err", err.Error()))
				return rabbitmq.NackDiscard
			}

			return next(enrichedCtx, d)
		}
	}
}
