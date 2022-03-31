package resolvers

import (
	"context"
	"errors"
	"strconv"

	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
	"github.com/jerbob92/hoppscotch-backend/models"

	"github.com/graph-gophers/graphql-go"
	"gorm.io/gorm"
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
	db := r.c.GetDB()
	existingUser := &models.User{}
	err := db.Where("id = ?", r.team_invitation.UserID).First(existingUser).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("user not found")
	}

	return NewUserResolver(r.c, existingUser)
}

func (r *TeamInvitationResolver) CreatorUid() (graphql.ID, error) {
	db := r.c.GetDB()
	existingUser := &models.User{}
	err := db.Where("id = ?", r.team_invitation.UserID).First(existingUser).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return graphql.ID(""), errors.New("user not found")
	}

	return graphql.ID(existingUser.FBUID), nil
}

func (r *TeamInvitationResolver) InviteeEmail() (graphql.ID, error) {
	return graphql.ID(r.team_invitation.InviteeEmail), nil
}

func (r *TeamInvitationResolver) InviteeRole() (models.TeamMemberRole, error) {
	return r.team_invitation.InviteeRole, nil
}

func (r *TeamInvitationResolver) Team() (*TeamResolver, error) {
	db := r.c.GetDB()
	existingTeam := &models.Team{}
	err := db.Where("id = ?", r.team_invitation.TeamID).First(existingTeam).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("team not found")
	}

	return NewTeamResolver(r.c, existingTeam)
}

func (r *TeamInvitationResolver) TeamID() (graphql.ID, error) {
	return graphql.ID(strconv.Itoa(int(r.team_invitation.TeamID))), nil
}

type TeamInvitationArgs struct {
	InviteID graphql.ID
}

func (b *BaseQuery) TeamInvitation(ctx context.Context, args *TeamInvitationArgs) (*TeamInvitationResolver, error) {
	// @todo: implement me
	// Check user/email.
	// {"errors":[{"message":"team_invite/email_do_not_match","locations":[{"line":2,"column":3}],"path":["teamInvitation"],"extensions":{"code":"INTERNAL_SERVER_ERROR"}}],"data":null}
	return nil, nil
}

type AcceptTeamInvitationArgs struct {
	InviteID graphql.ID
}

func (b *BaseQuery) AcceptTeamInvitation(ctx context.Context, args *AcceptTeamInvitationArgs) (*TeamMemberResolver, error) {
	// @todo: implement me
	// Check user/email.
	// {"errors":[{"message":"team_invite/email_do_not_match","locations":[{"line":2,"column":3}],"path":["teamInvitation"],"extensions":{"code":"INTERNAL_SERVER_ERROR"}}],"data":null}
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
