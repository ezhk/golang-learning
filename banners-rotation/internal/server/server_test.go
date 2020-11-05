// +build integration

package server

import (
	"context"
	"net"
	"testing"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/ezhk/golang-learning/banners-rotation/internal/queue"
	"github.com/ezhk/golang-learning/banners-rotation/internal/storage"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
)

func MakeDialer() func(context.Context, string) (net.Conn, error) {
	cfg := config.NewConfig("testdata/config.yaml")
	storage, _ := storage.NewStorage(cfg)
	queue, _ := queue.NewQueue(cfg)

	listener := bufconn.Listen(8192)
	server := grpc.NewServer()

	RegisterBannerServer(server, &Server{storage: storage, queue: queue})
	go server.Serve(listener)

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

type ServerTestSuite struct {
	suite.Suite
	conn   *grpc.ClientConn
	client BannerClient
	ctx    context.Context
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) SetupTest() {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(MakeDialer()), grpc.WithInsecure())
	s.NoError(err)
	s.conn = conn

	client := NewBannerClient(conn)
	s.client = client

	s.ctx = context.Background()
}

func (s *ServerTestSuite) TearDownTest() {
	s.conn.Close()
}

func (s *ServerTestSuite) TestBanners() {
	banners, err := s.client.ReadBanners(s.ctx, &emptypb.Empty{})
	s.NoError(err)
	initBannersLen := len(banners.Objects)

	res, err := s.client.CreateBanner(s.ctx, &SimpleCreateRequest{Name: "test grpc banner"})
	s.NoError(err)
	defer s.client.DeleteBanner(s.ctx, &SimpleRequestID{ID: res.ID})

	s.Equal("test grpc banner", res.Name)

	banners, err = s.client.ReadBanners(s.ctx, &emptypb.Empty{})
	s.NoError(err)
	s.Equal(initBannersLen+1, len(banners.Objects))

	resUpdated, err := s.client.UpdateBanner(s.ctx, &SimpleUpdateRequest{ID: res.ID, Name: "updated test grpc banner"})
	s.NoError(err)
	s.Equal(res.ID, resUpdated.ID)
	s.Equal("updated test grpc banner", resUpdated.Name)

	_, err = s.client.DeleteBanner(s.ctx, &SimpleRequestID{ID: res.ID})
	s.NoError(err)

	banners, err = s.client.ReadBanners(s.ctx, &emptypb.Empty{})
	s.NoError(err)
	s.Equal(initBannersLen, len(banners.Objects))
}

func (s *ServerTestSuite) TestSlots() {
	slots, err := s.client.ReadSlots(s.ctx, &emptypb.Empty{})
	s.NoError(err)
	initSlotsLen := len(slots.Objects)

	res, err := s.client.CreateSlot(s.ctx, &SimpleCreateRequest{Name: "test grpc slot"})
	s.NoError(err)
	defer s.client.DeleteSlot(s.ctx, &SimpleRequestID{ID: res.ID})

	s.Equal("test grpc slot", res.Name)

	slots, err = s.client.ReadSlots(s.ctx, &emptypb.Empty{})
	s.NoError(err)
	s.Equal(initSlotsLen+1, len(slots.Objects))

	resUpdated, err := s.client.UpdateSlot(s.ctx, &SimpleUpdateRequest{
		ID:   res.ID,
		Name: "updated test grpc slot",
	})
	s.NoError(err)
	s.Equal(res.ID, resUpdated.ID)
	s.Equal("updated test grpc slot", resUpdated.Name)

	_, err = s.client.DeleteSlot(s.ctx, &SimpleRequestID{ID: res.ID})
	s.NoError(err)

	slots, err = s.client.ReadSlots(s.ctx, &emptypb.Empty{})
	s.NoError(err)
	s.Equal(initSlotsLen, len(slots.Objects))
}

func (s *ServerTestSuite) TestGroups() {
	groups, err := s.client.ReadGroups(s.ctx, &emptypb.Empty{})
	s.NoError(err)
	initGroupsLen := len(groups.Objects)

	res, err := s.client.CreateGroup(s.ctx, &SimpleCreateRequest{Name: "test grpc group"})
	s.NoError(err)
	defer s.client.DeleteGroup(s.ctx, &SimpleRequestID{ID: res.ID})

	s.Equal("test grpc group", res.Name)

	groups, err = s.client.ReadGroups(s.ctx, &emptypb.Empty{})
	s.NoError(err)
	s.Equal(initGroupsLen+1, len(groups.Objects))

	resUpdated, err := s.client.UpdateGroup(s.ctx, &SimpleUpdateRequest{
		ID:   res.ID,
		Name: "updated test grpc group",
	})
	s.NoError(err)
	s.Equal(res.ID, resUpdated.ID)
	s.Equal("updated test grpc group", resUpdated.Name)

	_, err = s.client.DeleteGroup(s.ctx, &SimpleRequestID{ID: res.ID})
	s.NoError(err)

	groups, err = s.client.ReadGroups(s.ctx, &emptypb.Empty{})
	s.NoError(err)
	s.Equal(initGroupsLen, len(groups.Objects))
}

func (s *ServerTestSuite) TestPlacements() {
	banner, err := s.client.CreateBanner(s.ctx, &SimpleCreateRequest{Name: "test placements banner"})
	s.NoError(err)
	defer s.client.DeleteBanner(s.ctx, &SimpleRequestID{ID: banner.ID})

	slot, err := s.client.CreateSlot(s.ctx, &SimpleCreateRequest{Name: "test placements slot"})
	s.NoError(err)
	defer s.client.DeleteSlot(s.ctx, &SimpleRequestID{ID: slot.ID})

	group, err := s.client.CreateGroup(s.ctx, &SimpleCreateRequest{Name: "test placements group"})
	s.NoError(err)
	defer s.client.DeleteGroup(s.ctx, &SimpleRequestID{ID: group.ID})

	placements, err := s.client.ReadPlacements(s.ctx, &emptypb.Empty{})
	s.NoError(err)
	initPlacementsLen := len(placements.Objects)

	placement, err := s.client.CreatePlacement(s.ctx, &PlacementCreateRequest{
		BannerID: banner.ID,
		SlotID:   slot.ID,
		GroupID:  group.ID,
	})
	s.NoError(err)
	defer s.client.DeletePlacement(s.ctx, &SimpleRequestID{ID: placement.ID})

	placements, err = s.client.ReadPlacements(s.ctx, &emptypb.Empty{})
	s.NoError(err)
	s.Equal(initPlacementsLen+1, len(placements.Objects))

	secondBanner, err := s.client.CreateBanner(s.ctx, &SimpleCreateRequest{Name: "second test placements banner"})
	s.NoError(err)
	defer s.client.DeleteBanner(s.ctx, &SimpleRequestID{ID: secondBanner.ID})

	resUpdated, err := s.client.UpdatePlacement(s.ctx, &PlacementUpdateRequest{
		ID:       placement.ID,
		BannerID: secondBanner.ID,
		SlotID:   slot.ID,
		GroupID:  group.ID,
	})
	s.NoError(err)

	s.Equal(resUpdated.BannerID, secondBanner.ID)

	_, err = s.client.DeletePlacement(s.ctx, &SimpleRequestID{ID: placement.ID})
	s.NoError(err)

	placements, err = s.client.ReadPlacements(s.ctx, &emptypb.Empty{})
	s.NoError(err)
	s.Equal(initPlacementsLen, len(placements.Objects))
}

func (s *ServerTestSuite) TestEvents() {
	banner, err := s.client.CreateBanner(s.ctx, &SimpleCreateRequest{Name: "test placements banner"})
	s.NoError(err)
	defer s.client.DeleteBanner(s.ctx, &SimpleRequestID{ID: banner.ID})

	slot, err := s.client.CreateSlot(s.ctx, &SimpleCreateRequest{Name: "test placements slot"})
	s.NoError(err)
	defer s.client.DeleteSlot(s.ctx, &SimpleRequestID{ID: slot.ID})

	group, err := s.client.CreateGroup(s.ctx, &SimpleCreateRequest{Name: "test placements group"})
	s.NoError(err)
	defer s.client.DeleteGroup(s.ctx, &SimpleRequestID{ID: group.ID})

	placement, err := s.client.CreatePlacement(s.ctx, &PlacementCreateRequest{
		BannerID: banner.ID,
		SlotID:   slot.ID,
		GroupID:  group.ID,
	})
	s.NoError(err)
	defer s.client.DeletePlacement(s.ctx, &SimpleRequestID{ID: placement.ID})

	showedBanner, err := s.client.BannerShow(s.ctx, &BannerShowRequest{
		SlotID:  slot.ID,
		GroupID: group.ID,
	})
	s.NoError(err)
	s.Equal(banner.ID, showedBanner.Banner.ID)

	res, err := s.client.BannerClick(s.ctx, &SimpleRequestID{ID: placement.ID})
	s.NoError(err)
	s.Equal(placement.ID, res.ID)
}
