package workers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// Interface for all workers to implement
type SQSWorker interface {
	Handle(context.Context, types.Message) error
}

type Dispatcher interface {
	Start(context.Context, SQSWorker)
}
