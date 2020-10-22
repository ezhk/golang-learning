package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

func (s Server) BannerShow(ctx context.Context, empty *empty.Empty) (*PlacementResponse, error) {
	bannerPlacement, err := s.storage.ReadBannerHighestScore(nil)
	if err != nil {
		return nil, err
	}

	err = s.storage.ProcessBannerEvent(bannerPlacement.ID, "show")
	if err != nil {
		return nil, err
	}

	return ConvertBannerPlacementToPlacementResponse(bannerPlacement), nil
}

func (s Server) BannerClick(ctx context.Context, r *SimpleRequestID) (*SimpleResponseID, error) {
	err := s.storage.ProcessBannerEvent(r.ID, "click")
	if err != nil {
		return nil, err
	}

	return &SimpleResponseID{ID: r.ID}, nil
}
