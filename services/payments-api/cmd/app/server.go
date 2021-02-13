package app

import (
	"context"
	"log"
	"microservice/services/payments-api/pkg/payments"
	serverPb "microservice/services/payments-api/pkg/server"

)

type Server struct {
	paymentsSvc *payments.Service
	ctx 			context.Context
}

func NewServer(transactionsSvc *payments.Service, ctx context.Context) *Server {
	return &Server{paymentsSvc: transactionsSvc, ctx: ctx}
}

func (s *Server) Pay(ctx context.Context, request * serverPb.PaymentTransactionsRequest) (*serverPb.PaymentTransactionsResponse, error){

	userID := request.UserID
	amount := request.Amount
	category := request.Category
	err := s.paymentsSvc.Pay(ctx, userID, amount, category)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &serverPb.PaymentTransactionsResponse{Ok: true}, nil
}