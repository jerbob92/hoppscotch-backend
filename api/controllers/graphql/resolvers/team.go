package resolvers

import (
	"context"
	"github.com/graph-gophers/graphql-go"
	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
	"github.com/jerbob92/hoppscotch-backend/models"
	"strconv"
)

type TeamResolver struct {
	c    *graphql_context.Context
	team *models.Team
}

func NewTeamResolver(c *graphql_context.Context, team *models.Team) (*TeamResolver, error) {
	if team == nil {
		return nil, nil
	}

	return &TeamResolver{c: c, team: team}, nil
}

func (r *TeamResolver) ID() (graphql.ID, error) {
	id := graphql.ID(strconv.FormatInt(r.team.ID, 10))
	return id, nil
}

func (r *TeamResolver) EditorsCount() (int32, error) {
	return 0, nil
}

type TeamMembersArgs struct {
	Cursor *graphql.ID
}

func (r *TeamResolver) Members(args *TeamMembersArgs) ([]*TeamMemberResolver, error) {
	return nil, nil
}

func (r *TeamResolver) MyRole() (models.TeamMemberRole, error) {
	return models.Owner, nil
}

func (r *TeamResolver) Name() (string, error) {
	return "", nil
}

func (r *TeamResolver) OwnersCount() (int32, error) {
	return 0, nil
}

func (r *TeamResolver) TeamInvitations() ([]*TeamInvitationResolver, error) {
	return nil, nil
}

func (r *TeamResolver) TeamMembers() ([]*TeamMemberResolver, error) {
	return nil, nil
}

func (r *TeamResolver) ViewersCount() (int32, error) {
	return 0, nil
}

type MyTeamsArgs struct {
	Cursor *graphql.ID
}

func (b *BaseQuery) MyTeams(ctx context.Context, args *MyTeamsArgs) ([]*TeamResolver, error) {
	// @todo: implement me
	return nil, nil
}

type RequestArg struct {
	RequestID graphql.ID
}

func (b *BaseQuery) Request(ctx context.Context, args *RequestArg) (*TeamRequestResolver, error) {
	// @todo: implement me
	return nil, nil
}

type RootCollectionsOfTeamArgs struct {
	Cursor *graphql.ID
	TeamID graphql.ID
}

func (b *BaseQuery) RootCollectionsOfTeam(ctx context.Context, args *RootCollectionsOfTeamArgs) ([]*TeamCollectionResolver, error) {
	// @todo: implement me
	return nil, nil
}

type SearchForRequestArgs struct {
	Cursor     *graphql.ID
	SearchTerm string
	TeamID     graphql.ID
}

func (b *BaseQuery) SearchForRequest(ctx context.Context, args *SearchForRequestArgs) ([]*TeamRequestResolver, error) {
	// @todo: implement me
	return nil, nil
}

type TeamArgs struct {
	TeamID graphql.ID
}

func (b *BaseQuery) Team(ctx context.Context, args *TeamArgs) (*TeamResolver, error) {
	// @todo: implement me
	return nil, nil
}

type CreateTeamArgs struct {
	Name string
}

func (b *BaseQuery) CreateTeam(ctx context.Context, args *CreateTeamArgs) (*TeamResolver, error) {
	// @todo: implement me
	return nil, nil
}

type DeleteTeamArgs struct {
	TeamID graphql.ID
}

func (b *BaseQuery) DeleteTeam(ctx context.Context, args *DeleteTeamArgs) (bool, error) {
	// @todo: implement me
	return false, nil
}

type LeaveTeamArgs struct {
	TeamID graphql.ID
}

func (b *BaseQuery) LeaveTeam(ctx context.Context, args *LeaveTeamArgs) (bool, error) {
	// @todo: implement me
	return false, nil
}

type RenameTeamArgs struct {
	NewName string
	TeamID  graphql.ID
}

func (b *BaseQuery) RenameTeam(ctx context.Context, args *RenameTeamArgs) (*TeamResolver, error) {
	// @todo: implement me
	return nil, nil
}
