package app

import (
	"github.com/DaniilOr/microservice/auth/pkg/auth"
	serverPb "github.com/DaniilOr/microservice/auth/pkg/server"
	"context"
	"go.opencensus.io/trace"
	"log"
)

type Server struct {
	authSvc *auth.Service
	ctx context.Context
}

func NewServer(authSvc *auth.Service, ctx context.Context) *Server {
	return &Server{authSvc: authSvc, ctx: ctx }
}

func (s *Server) Token(ctx context.Context, request *serverPb.TokenRequest) ( * serverPb.TokenResponse, error) {
	ctx, span := trace.StartSpan(ctx, "route: token")
	defer span.End()
	token, err := s.authSvc.Login(ctx, request.Login, request.Password)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	response := serverPb.TokenResponse{Token: token}
	return &response, nil
}

// Доступно всем
func (s *Server) Id (ctx context.Context, request *serverPb.IdRequest) (*serverPb.IdResponse, error) {
	userID, err := s.authSvc.UserID(ctx, request.Token)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	response := serverPb.IdResponse{UserId: userID}
	return &response, nil
}