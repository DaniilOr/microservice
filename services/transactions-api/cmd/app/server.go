package app

import (
	"context"
	"go.opencensus.io/trace"
	"log"
	serverPb "transactions-api/pkg/server"
	"transactions-api/pkg/transactions"
)

type Server struct {
	transactionsSvc *transactions.Service
	ctx 			context.Context
}

func NewServer(transactionsSvc *transactions.Service, ctx context.Context) *Server {
	return &Server{transactionsSvc: transactionsSvc, ctx: ctx}
}

func (s *Server) Transactions(ctx context.Context, request * serverPb.TransactionsRequest) (*serverPb.TransactionsResponse, error){
	ctx, span := trace.StartSpan(ctx, "route: transactions")
	defer span.End()
	userID := request.UserID
	data, err := s.transactionsSvc.Transactions(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &serverPb.TransactionsResponse{Response: data}, nil
}