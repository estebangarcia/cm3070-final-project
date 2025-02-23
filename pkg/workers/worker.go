package workers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// Interface for all workers to implement
type SQSWorker interface {
	Handle(context.Context, types.Message) error
}

type SQSDispatcher interface {
	Start(context.Context, SQSWorker)
}

type PeriodicWorker interface {
	Handle(context.Context) error
}

type PeriodicDispatcher interface {
	Start(context.Context, PeriodicWorker)
}
