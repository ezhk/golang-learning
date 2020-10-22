package server

import (
	"context"

	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"github.com/golang/protobuf/ptypes/empty"
)

//   // Placement methods.
//   rpc CreatePlacement(PlacementCreateRequest) returns(PlacementResponse) {
//     option(google.api.http) = {post : "/api/v1/placements" body : "*"};
//   }
//   rpc ReadPlacements(google.protobuf.Empty) returns(MultiplePlacementResponse) {
//     option(google.api.http) = {get : "/api/v1/placements"};
//   }
//   rpc UpdatePlacement(SimpleUpdateRequest) returns(PlacementResponse) {
//     option(google.api.http) = {put : "/api/v1/placements/{ID}" body : "*"};
//   }
//   rpc DeletePlacement(SimpleRequestID) returns(SimpleResponseID) {
//     option(google.api.http) = {delete : "/api/v1/placements/{ID}"};
//   }

func ConvertBannerPlacementToPlacementResponse(p structs.BannerPlacement) *PlacementResponse {
	return &PlacementResponse{
		ID:     p.ID,
		Banner: ConvertBannerToSimpleResponse(p.Banner),
		Slot:   ConvertSlotToSimpleResponse(p.Slot),
		Group:  ConvertGroupToSimpleResponse(p.Group),
		Shows:  p.Shows,
		Clicks: p.Clicks,
		Score:  p.Score,
	}
}

func ConvertPlacementUpdateRequestToBannerPlacement(r *PlacementUpdateRequest) structs.BannerPlacement {
	return structs.BannerPlacement{
		ID:     r.ID,
		Banner: ConvertSimpleUpdateRequestToBanner(r.Banner),
		Slot:   ConvertSimpleUpdateRequestToSlot(r.Slot),
		Group:  ConvertSimpleUpdateRequestToGroup(r.Group),
		Shows:  r.Shows,
		Clicks: r.Clicks,
		Score:  r.Score,
	}
}

func (s Server) CreatePlacement(ctx context.Context, r *PlacementCreateRequest) (*PlacementResponse, error) {
	placement, err := s.storage.CreateBannerPlacement(r.BannerID, r.SlotID, r.GroupID)
	if err != nil {
		return nil, err
	}

	return ConvertBannerPlacementToPlacementResponse(placement), nil
}

func (s Server) ReadPlacements(ctx context.Context, empty *empty.Empty) (*MultiplePlacementResponse, error) {
	placements, err := s.storage.ReadBannersPlacements(nil)
	if err != nil {
		return nil, err
	}

	response := make([]*PlacementResponse, 0)
	for _, p := range placements {
		response = append(response, ConvertBannerPlacementToPlacementResponse(*p))
	}

	return &MultiplePlacementResponse{Objects: response}, nil
}

func (s Server) UpdatePlacement(ctx context.Context, r *PlacementUpdateRequest) (*PlacementResponse, error) {
	placement := ConvertPlacementUpdateRequestToBannerPlacement(r)
	p, err := s.storage.UpdateBannerPlacement(placement)
	if err != nil {
		return nil, err
	}

	return ConvertBannerPlacementToPlacementResponse(p), nil
}

func (s Server) DeletePlacement(ctx context.Context, r *SimpleRequestID) (*SimpleResponseID, error) {
	if err := s.storage.DeleteBannerPlacement(r.ID); err != nil {
		return nil, err
	}

	return &SimpleResponseID{ID: r.ID}, nil
}
