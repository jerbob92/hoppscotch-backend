package resolvers

import (
	"context"
	"github.com/graph-gophers/graphql-go"
	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
	"github.com/jerbob92/hoppscotch-backend/models"
	"strconv"
)

type TeamMemberResolver struct {
	c           *graphql_context.Context
	team_member *models.TeamMember
}

func NewTeamMemberResolver(c *graphql_context.Context, team_member *models.TeamMember) (*TeamMemberResolver, error) {
	if team_member == nil {
		return nil, nil
	}

	return &TeamMemberResolver{c: c, team_member: team_member}, nil
}

func (r *TeamMemberResolver) MembershipID() (graphql.ID, error) {
	id := graphql.ID(strconv.FormatInt(r.team_member.ID, 10))
	return id, nil
}

func (r *TeamMemberResolver) Role() (models.TeamMemberRole, error) {
	return r.team_member.Role, nil
}

func (r *TeamMemberResolver) User() (*UserResolver, error) {
	return nil, nil
}

type RemoveTeamMemberArgs struct {
	TeamID  graphql.ID
	UserUID graphql.ID
}

func (b *BaseQuery) RemoveTeamMember(ctx context.Context, args *RemoveTeamMemberArgs) (bool, error) {
	// @todo: implement me
	return false, nil
}

type UpdateTeamMemberRoleArgs struct {
	NewRole models.TeamMemberRole
	TeamID  graphql.ID
	UserUID graphql.ID
}

func (b *BaseQuery) UpdateTeamMemberRole(ctx context.Context, args *UpdateTeamMemberRoleArgs) (*TeamMemberResolver, error) {
	// @todo: implement me
	return nil, nil
}
