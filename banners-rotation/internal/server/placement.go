package server

import (
	"context"

	"github.com/ezhk/golang-learning/banners-rotation/internal/api"
	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"github.com/golang/protobuf/ptypes/empty"
)

func ConvertBannerPlacementToPlacementResponse(p structs.BannerPlacement) *api.PlacementResponse {
	return &api.PlacementResponse{
		ID:     p.ID,
		Banner: ConvertBannerToSimpleResponse(p.Banner),
		Slot:   ConvertSlotToSimpleResponse(p.Slot),
		Group:  ConvertGroupToSimpleResponse(p.Group),
		Shows:  p.Shows,
		Clicks: p.Clicks,
		Score:  p.Score,
	}
}

func ConvertBannerPlacementToPlacementIDsResponse(p structs.BannerPlacement) *api.PlacementIDsResponse {
	return &api.PlacementIDsResponse{
		ID:       p.ID,
		BannerID: p.BannerID,
		SlotID:   p.SlotID,
		GroupID:  p.GroupID,
	}
}

func ConvertPlacementUpdateRequestToBannerPlacement(r *api.PlacementUpdateRequest) structs.BannerPlacement {
	return structs.BannerPlacement{
		ID:       r.ID,
		BannerID: r.BannerID,
		SlotID:   r.SlotID,
		GroupID:  r.GroupID,
		Shows:    r.Shows,
		Clicks:   r.Clicks,
		Score:    r.Score,
	}
}

func (s Server) CreatePlacement(ctx context.Context, r *api.PlacementCreateRequest) (*api.PlacementIDsResponse, error) {
	placement, err := s.storage.CreateBannerPlacement(r.BannerID, r.SlotID, r.GroupID)
	if err != nil {
		return nil, err
	}

	return ConvertBannerPlacementToPlacementIDsResponse(placement), nil
}

func (s Server) ReadPlacements(ctx context.Context, empty *empty.Empty) (*api.MultiplePlacementResponse, error) {
	placements, err := s.storage.ReadBannersPlacements(nil)
	if err != nil {
		return nil, err
	}

	response := make([]*api.PlacementResponse, 0)
	for _, p := range placements {
		response = append(response, ConvertBannerPlacementToPlacementResponse(*p))
	}

	return &api.MultiplePlacementResponse{Objects: response}, nil
}

func (s Server) UpdatePlacement(ctx context.Context, r *api.PlacementUpdateRequest) (*api.PlacementIDsResponse, error) {
	placement := ConvertPlacementUpdateRequestToBannerPlacement(r)
	p, err := s.storage.UpdateBannerPlacement(placement)
	if err != nil {
		return nil, err
	}

	return ConvertBannerPlacementToPlacementIDsResponse(p), nil
}

func (s Server) DeletePlacement(ctx context.Context, r *api.SimpleRequestID) (*api.SimpleResponseID, error) {
	if err := s.storage.DeleteBannerPlacement(r.ID); err != nil {
		return nil, err
	}

	return &api.SimpleResponseID{ID: r.ID}, nil
}
