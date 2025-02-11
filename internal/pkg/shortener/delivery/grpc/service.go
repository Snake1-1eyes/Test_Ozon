package grpc

import (
	"context"

	"github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener"
	shortenerpb "github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener/delivery/grpc/shortenerpb"
)

// Server реализует интерфейс ShortenerService gRPC.
type Server struct {
	shortenerpb.UnimplementedShortenerServiceServer
	usecase shortener.ShortenerUsecase
}

func NewServer(uc shortener.ShortenerUsecase) *Server {
	return &Server{usecase: uc}
}

func (s *Server) CreateShortLink(ctx context.Context, req *shortenerpb.CreateShortLinkRequest) (*shortenerpb.CreateShortLinkResponse, error) {
	shortURL, err := s.usecase.CreateAndSaveShortLink(req.GetOriginalUrl())
	if err != nil {
		return nil, err
	}
	return &shortenerpb.CreateShortLinkResponse{ShortUrl: shortURL}, nil
}

func (s *Server) GetShortLink(ctx context.Context, req *shortenerpb.GetShortLinkRequest) (*shortenerpb.GetShortLinkResponse, error) {
	originalURL, err := s.usecase.GetShortLink(req.GetShortUrl())
	if err != nil {
		return nil, err
	}
	return &shortenerpb.GetShortLinkResponse{OriginalUrl: originalURL}, nil
}
