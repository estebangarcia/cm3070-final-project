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

func (a *DbTxMiddleware) HandleTransaction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		srw := &statusResponseWriter{ResponseWriter: w}

		isMutable := isMutableMethod(r)

		var tx *ent.Tx
		var err error

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

		defer func(hasTx bool, tx *ent.Tx) {
			if !isMutable {
				return
			}

			if srw.status >= http.StatusBadRequest {
				_ = tx.Rollback()
				return
			}

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
