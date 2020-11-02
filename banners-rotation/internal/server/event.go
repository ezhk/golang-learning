package server

import (
	"context"

	"github.com/ezhk/golang-learning/banners-rotation/internal/api"
	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
)

func (s Server) BannerShow(ctx context.Context, r *api.BannerShowRequest) (*api.PlacementResponse, error) {
	bannerPlacement, err := s.storage.ReadBannerHighestScore(structs.BannerFilter{
		"slot_id":  r.SlotID,
		"group_id": r.GroupID,
	})
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

	return ConvertBannerPlacementToPlacementResponse(bannerPlacement), nil
}

func (s Server) BannerClick(ctx context.Context, r *api.SimpleRequestID) (*api.SimpleResponseID, error) {
	err := s.queue.ProduceEvent(structs.QueueEvent{
		PlacementID: r.ID,
		EventType:   "click",
	})
	if err != nil {
		return nil, err
	}

	return &api.SimpleResponseID{ID: r.ID}, nil
}
