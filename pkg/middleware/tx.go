package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
)

type DbTxMiddleware struct {
	DBClient *ent.Client
}

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// This middleware checks if the request method is POST, PUT, DELETE or PATCH
// it then considers the request will mutate data in the database and wraps
// the whole request inside a database transaction, this transaction is automatically
// committed or rolled back based on the outcome of the request.
func (a *DbTxMiddleware) HandleTransaction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		srw := &statusResponseWriter{ResponseWriter: w}

		isMutable := isMutableMethod(r)

		var tx *ent.Tx
		var err error

		// If request is mutable then create a transaction and store it in the context
		if isMutable {
			tx, err = a.DBClient.Tx(ctx)
			if err != nil {
				log.Println(err)
				http.Error(w, "Internal Server Error", 500)
				return
			}

			ctx = context.WithValue(ctx, "dbClient", tx.Client())
		} else {
			ctx = context.WithValue(ctx, "dbClient", a.DBClient)
		}

		// Defer this function that will be executed after the request is completed
		defer func(hasTx bool, tx *ent.Tx) {
			if !isMutable {
				return
			}

			// If the http status of the request is greater than or equal to 400
			// rollback the transaction
			if srw.status >= http.StatusBadRequest {
				_ = tx.Rollback()
				return
			}

			// if the status is successful then we commit
			if err := tx.Commit(); err != nil {
				http.Error(srw, "Internal Server Error", http.StatusInternalServerError)
			}
		}(isMutable, tx)

		next.ServeHTTP(srw, r.WithContext(ctx))
	})
}

func isMutableMethod(r *http.Request) bool {
	return (r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" || r.Method == "PATCH")
}
