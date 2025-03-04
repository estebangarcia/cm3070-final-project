package workers

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"golang.org/x/sync/errgroup"
)

type SQSWorkerDispatcher struct {
	maxMessages int32
	queueUrl    string
	sqsClient   *sqs.Client
	wg          *sync.WaitGroup
	mux         *sync.Mutex
	errGroup    *errgroup.Group
	dbClient    *ent.Client
}

func NewSQSWorkerDispatcher(sqsQueueUrl string, sqsClient *sqs.Client, maxMessages int32, dbClient *ent.Client) *SQSWorkerDispatcher {
	return &SQSWorkerDispatcher{
		queueUrl:    sqsQueueUrl,
		sqsClient:   sqsClient,
		maxMessages: maxMessages,
		wg:          &sync.WaitGroup{},
		mux:         &sync.Mutex{},
		errGroup:    &errgroup.Group{},
		dbClient:    dbClient,
	}
}

func (w *SQSWorkerDispatcher) Start(ctx context.Context, worker SQSWorker) {
	w.errGroup.Go(func() error {
		for {
			if ctx.Err() != nil {
				log.Println("context has been cancelled")
				return nil
			}

			results, err := w.sqsClient.ReceiveMessage(
				ctx,
				&sqs.ReceiveMessageInput{
					QueueUrl:            &w.queueUrl,
					MaxNumberOfMessages: 10,
					WaitTimeSeconds:     20,
				},
			)

			if err != nil {
				log.Printf("error consuming from sqs %v", err)
				continue
			}

			msgAck := []types.DeleteMessageBatchRequestEntry{}

			for _, message := range results.Messages {
				w.wg.Add(1)
				go func() {
					defer w.wg.Done()

					tx, err := w.dbClient.Tx(ctx)
					if err != nil {
						return
					}
					ctxTx := context.WithValue(ctx, "dbClient", tx.Client())
					defer tx.Rollback()

					if err := worker.Handle(ctxTx, message); err == nil {
						if err := tx.Commit(); err != nil {
							fmt.Printf("error commiting transaction %v", err)
							return
						}
						w.mux.Lock()
						msgAck = append(msgAck, types.DeleteMessageBatchRequestEntry{
							Id:            message.MessageId,
							ReceiptHandle: message.ReceiptHandle,
						})
						w.mux.Unlock()
					}
				}()
			}
			w.wg.Wait()

			if len(msgAck) > 0 {
				_, err = w.sqsClient.DeleteMessageBatch(ctx, &sqs.DeleteMessageBatchInput{
					Entries:  msgAck,
					QueueUrl: &w.queueUrl,
				})
				if err != nil {
					log.Printf("error acking messages %v", err)
				}
			}
		}
	})
}
