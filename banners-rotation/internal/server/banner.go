package server

import (
	"context"

	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertBannerToSimpleResponse(b structs.Banner) *SimpleResponse {
	return &SimpleResponse{
		ID:        b.ID,
		Name:      b.Name,
		CreatedAt: timestamppb.New(b.CreatedAt),
		UpdatedAt: timestamppb.New(b.UpdatedAt),
	}
}

func ConvertSimpleUpdateRequestToBanner(r *SimpleUpdateRequest) structs.Banner {
	return structs.Banner{
		ID:   r.ID,
		Name: r.Name,
	}
}

func (s Server) CreateBanner(ctx context.Context, r *SimpleCreateRequest) (*SimpleResponse, error) {
	banner, err := s.storage.CreateBanner(r.Name)
	if err != nil {
		return nil, err
	}

	return ConvertBannerToSimpleResponse(banner), nil
}

func (s Server) ReadBanners(ctx context.Context, empty *empty.Empty) (*MultipleSimpleResponse, error) {
	banners, err := s.storage.ReadBanners()
	if err != nil {
		return nil, err
	}

	simpleResponses := make([]*SimpleResponse, 0)
	for _, banner := range banners {
		simpleResponses = append(simpleResponses, ConvertBannerToSimpleResponse(*banner))
	}

	return &MultipleSimpleResponse{Objects: simpleResponses}, nil
}

func (s Server) UpdateBanner(ctx context.Context, r *SimpleUpdateRequest) (*SimpleResponse, error) {
	banner := ConvertSimpleUpdateRequestToBanner(r)
	b, err := s.storage.UpdateBanner(banner)
	if err != nil {
		return nil, err
	}

	return ConvertBannerToSimpleResponse(b), nil
}

func (s Server) DeleteBanner(ctx context.Context, r *SimpleDeleteRequest) (*SimpleDeleteResponse, error) {
	if err := s.storage.DeleteBanner(r.ID); err != nil {
		return nil, err
	}

	return &SimpleDeleteResponse{ID: r.ID}, nil
}
