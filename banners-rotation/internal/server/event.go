package server

import (
	"context"

	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s Server) BannerShow(ctx context.Context, empty *empty.Empty) (*PlacementResponse, error) {
	bannerPlacement, err := s.storage.ReadBannerHighestScore(nil)
	if err != nil {
		return nil, err
	}

	err = s.queue.ProduceEvent(structs.QueueEvent{
		PlacementID: bannerPlacement.ID,
		EventType:   "show",
	})
	if err != nil {
		return nil, err
	}

	// err = s.storage.ProcessBannerEvent(bannerPlacement.ID, "show")
	// if err != nil {
	// 	return nil, err
	// }

	return ConvertBannerPlacementToPlacementResponse(bannerPlacement), nil
}

func (s Server) BannerClick(ctx context.Context, r *SimpleRequestID) (*SimpleResponseID, error) {
	err := s.queue.ProduceEvent(structs.QueueEvent{
		PlacementID: r.ID,
		EventType:   "click",
	})
	if err != nil {
		return nil, err
	}

	// err := s.storage.ProcessBannerEvent(r.ID, "click")
	// if err != nil {
	// 	return nil, err
	// }

	return &SimpleResponseID{ID: r.ID}, nil
}
