package listing_client

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	core_errors "messenger-service/internal/core/errors"
	"messenger-service/internal/grpc/listingpb"
)



// GRPCClient реализует тот же ListingClient-интерфейс что и HTTP Client,
// но ходит в listing-service по gRPC вместо HTTP.
type GRPCClient struct {
	stub listingpb.ListingServiceClient
}

func NewGRPCClient(stub listingpb.ListingServiceClient) *GRPCClient {
	return &GRPCClient{stub: stub}
}

func (c *GRPCClient) GetListing(ctx context.Context, id uuid.UUID) (Listing, error) {
	resp, err := c.stub.GetListing(ctx, &listingpb.GetListingRequest{Id: id.String()})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return Listing{}, core_errors.ErrNotFound 
		}
		return Listing{}, fmt.Errorf("grpc get listing: %w", err)
	}

	listingID, err := uuid.Parse(resp.GetId())
	if err != nil {
		return Listing{}, fmt.Errorf("parse listing id from grpc response: %w", err)
	}

	userID, err := uuid.Parse(resp.GetUserId())
	if err != nil {
		return Listing{}, fmt.Errorf("parse user id from grpc response: %w", err)
	}

	return Listing{
		ID:     listingID,
		UserID: userID,
	}, nil
}
