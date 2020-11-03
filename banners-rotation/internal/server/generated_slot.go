// generated by generate-grpc-methods; DO NOT EDIT
package server

import (
	"context"

	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertSlotToSimpleResponse(b structs.Slot) *SimpleResponse {
	return &SimpleResponse{
		ID:          b.ID,
		Name:        b.Name,
		Description: b.Description,
		CreatedAt:   timestamppb.New(b.CreatedAt),
		UpdatedAt:   timestamppb.New(b.UpdatedAt),
	}
}

func ConvertSimpleUpdateRequestToSlot(r *SimpleUpdateRequest) structs.Slot {
	return structs.Slot{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}
}

func (s Server) CreateSlot(ctx context.Context, r *SimpleCreateRequest) (*SimpleResponse, error) {
	slot, err := s.storage.CreateSlot(r.Name, r.Description)
	if err != nil {
		return nil, err
	}

	return ConvertSlotToSimpleResponse(slot), nil
}

func (s Server) ReadSlots(ctx context.Context, empty *empty.Empty) (*MultipleSimpleResponse, error) {
	slots, err := s.storage.ReadSlots()
	if err != nil {
		return nil, err
	}

	simpleResponses := make([]*SimpleResponse, 0)
	for _, slot := range slots {
		simpleResponses = append(simpleResponses, ConvertSlotToSimpleResponse(*slot))
	}

	return &MultipleSimpleResponse{Objects: simpleResponses}, nil
}

func (s Server) UpdateSlot(ctx context.Context, r *SimpleUpdateRequest) (*SimpleResponse, error) {
	slot := ConvertSimpleUpdateRequestToSlot(r)
	b, err := s.storage.UpdateSlot(slot)
	if err != nil {
		return nil, err
	}

	return ConvertSlotToSimpleResponse(b), nil
}

func (s Server) DeleteSlot(ctx context.Context, r *SimpleRequestID) (*SimpleResponseID, error) {
	if err := s.storage.DeleteSlot(r.ID); err != nil {
		return nil, err
	}

	return &SimpleResponseID{ID: r.ID}, nil
}
