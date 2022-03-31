package resolvers

import (
	"context"
	"github.com/graph-gophers/graphql-go"
	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
	"github.com/jerbob92/hoppscotch-backend/models"
	"strconv"
)

type TeamRequestResolver struct {
	c            *graphql_context.Context
	team_request *models.TeamRequest
}

func NewTeamRequestResolver(c *graphql_context.Context, team_request *models.TeamRequest) (*TeamRequestResolver, error) {
	if team_request == nil {
		return nil, nil
	}

	return &TeamRequestResolver{c: c, team_request: team_request}, nil
}

func (r *TeamRequestResolver) ID() (graphql.ID, error) {
	id := graphql.ID(strconv.FormatInt(r.team_request.ID, 10))
	return id, nil
}

func (r *TeamRequestResolver) Collection() (*TeamCollectionResolver, error) {
	return nil, nil
}

func (r *TeamRequestResolver) CollectionID() (graphql.ID, error) {
	return graphql.ID(""), nil
}

func (r *TeamRequestResolver) Request() (string, error) {
	return r.team_request.Request, nil
}

func (r *TeamRequestResolver) Team() (*TeamResolver, error) {
	return nil, nil
}

func (r *TeamRequestResolver) TeamID() (graphql.ID, error) {
	return graphql.ID(""), nil
}

func (r *TeamRequestResolver) Title() (string, error) {
	return r.team_request.Title, nil
}

type DeleteRequestArgs struct {
	RequestID graphql.ID
}

func (b *BaseQuery) DeleteRequest(ctx context.Context, args *DeleteRequestArgs) (bool, error) {
	// @todo: implement me
	return false, nil
}

type MoveRequestArgs struct {
	DestCollID graphql.ID
	RequestID  graphql.ID
}

func (b *BaseQuery) MoveRequest(ctx context.Context, args *MoveRequestArgs) (*TeamRequestResolver, error) {
	// @todo: implement me
	return nil, nil
}

type UpdateTeamRequestInput struct {
	Request *string
	Title   *string
}

type UpdateRequestArgs struct {
	Data      UpdateTeamRequestInput
	RequestID graphql.ID
}

func (b *BaseQuery) UpdateRequest(ctx context.Context, args *UpdateRequestArgs) (*TeamRequestResolver, error) {
	// @todo: implement me
	return nil, nil
}
