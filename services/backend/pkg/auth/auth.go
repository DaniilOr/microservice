package auth

import (
	serverPb "backend/pkg/auth_client"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"time"
)

type Service struct {
	client *http.Client
	url    string
}

func NewService(client *http.Client, url string) *Service {
	return &Service{client: client, url: url}
}

func (s *Service) Token(ctx context.Context, login string, password string) (token string, err error) {
	// for simplicity just define locally

	conn, err := grpc.Dial(s.url, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()

	client := serverPb.NewAuthServerClient(conn)
	resp, err := client.Token(ctx, &serverPb.TokenRequest{
		Login: login,
		Password: password,
	})

	if err != nil {
		if st, ok := status.FromError(err); ok {
			log.Print(st.Code())
			log.Print(st.Message())
		}
		return "", err
	}


	return resp.Token, nil
}

func (s *Service) Auth(ctx context.Context, token string) (userID int64, err error) {
	// for simplicity just define locally

	conn, err := grpc.Dial(s.url, grpc.WithInsecure())
	if err != nil {
		return 0, err
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()
	client := serverPb.NewAuthServerClient(conn)
	ctx, _ = context.WithTimeout(context.Background(), time.Second)
	resp, err := client.Id(ctx, &serverPb.IdRequest{
		Token: token,
	})
	if err != nil{
		log.Println(err)
		return 0, err
	}
	return resp.UserId, nil
}