package resolvers

import (
	"context"
	"strconv"

	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
	"github.com/jerbob92/hoppscotch-backend/models"

	"github.com/graph-gophers/graphql-go"
)

type TeamInvitationResolver struct {
	c               *graphql_context.Context
	team_invitation *models.TeamInvitation
}

func NewTeamInvitationResolver(c *graphql_context.Context, team_invitation *models.TeamInvitation) (*TeamInvitationResolver, error) {
	if team_invitation == nil {
		return nil, nil
	}

	return &TeamInvitationResolver{c: c, team_invitation: team_invitation}, nil
}

func (r *TeamInvitationResolver) ID() (graphql.ID, error) {
	id := graphql.ID(strconv.Itoa(int(r.team_invitation.ID)))
	return id, nil
}

func (r *TeamInvitationResolver) Creator() (*UserResolver, error) {
	return nil, nil
}

func (r *TeamInvitationResolver) CreatorUid() (graphql.ID, error) {
	return "", nil
}

func (r *TeamInvitationResolver) InviteeEmail() (graphql.ID, error) {
	return "", nil
}

func (r *TeamInvitationResolver) InviteeRole() (models.TeamMemberRole, error) {
	return r.team_invitation.InviteeRole, nil
}

func (r *TeamInvitationResolver) Team() (*TeamResolver, error) {
	return nil, nil
}

func (r *TeamInvitationResolver) TeamID() (graphql.ID, error) {
	return "", nil
}

type TeamInvitationArgs struct {
	InviteID graphql.ID
}

func (b *BaseQuery) TeamInvitation(ctx context.Context, args *TeamInvitationArgs) (*TeamInvitationResolver, error) {
	// @todo: implement me
	return nil, nil
}

type AcceptTeamInvitationArgs struct {
	InviteID graphql.ID
}

func (b *BaseQuery) AcceptTeamInvitation(ctx context.Context, args *AcceptTeamInvitationArgs) (*TeamMemberResolver, error) {
	// @todo: implement me
	return nil, nil
}

type AddTeamMemberByEmailArgs struct {
	TeamID    graphql.ID
	UserEmail string
	UserRole  models.TeamMemberRole
}

func (b *BaseQuery) AddTeamMemberByEmail(ctx context.Context, args *AddTeamMemberByEmailArgs) (*TeamMemberResolver, error) {
	// @todo: implement me
	return nil, nil
}

type CreateTeamInvitationArgs struct {
	InviteeEmail string
	InviteeRole  models.TeamMemberRole
	TeamID       graphql.ID
}

func (b *BaseQuery) CreateTeamInvitation(ctx context.Context, args *CreateTeamInvitationArgs) (*TeamInvitationResolver, error) {
	// @todo: implement me
	return nil, nil
}

type RevokeTeamInvitationArgs struct {
	InviteID graphql.ID
}

func (b *BaseQuery) RevokeTeamInvitation(ctx context.Context, args *RevokeTeamInvitationArgs) (bool, error) {
	// @todo: implement me
	return false, nil
}
