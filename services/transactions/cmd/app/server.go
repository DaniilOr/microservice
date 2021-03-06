package app

import (
	"context"
	"go.opencensus.io/trace"
	serverPb "transactions/pkg/server"
	"log"
	"transactions/pkg/transactions"
)

type Server struct {
	transactionsSvc *transactions.Service
	ctx context.Context
}

func NewServer(transactionsSvc *transactions.Service, ctx context.Context) *Server {
	return &Server{transactionsSvc: transactionsSvc, ctx: ctx}
}

func (s *Server) Transactions(ctx context.Context, request * serverPb.TransactionsRequest) (*serverPb.TransactionsResponse, error){
	ctx, span := trace.StartSpan(ctx, "route: transaction")
	defer span.End()
	userID := request.UserID

	records, err := s.transactionsSvc.Transactions(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var response serverPb.TransactionsResponse
	for _, record := range records {
		response.Items = append(response.Items, &serverPb.Transaction{
			Id:       record.ID,
			UserId:   record.UserID,
			Category: record.Category,
			Amount:   record.Amount,
			Created:  record.Created,
		})
	}
	return &response, nil
}