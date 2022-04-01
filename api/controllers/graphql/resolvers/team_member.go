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
	id := graphql.ID(strconv.Itoa(int(r.team_member.ID)))
	return id, nil
}

func (r *TeamMemberResolver) Role() (models.TeamMemberRole, error) {
	return r.team_member.Role, nil
}

func (r *TeamMemberResolver) User() (*UserResolver, error) {
	db := r.c.GetDB()
	existingUser := &models.User{}
	err := db.Where("id = ?", r.team_member.UserID).First(existingUser).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("user not found")
	}

	return NewUserResolver(r.c, existingUser)
}

type RemoveTeamMemberArgs struct {
	TeamID  graphql.ID
	UserUID graphql.ID
}

func (b *BaseQuery) RemoveTeamMember(ctx context.Context, args *RemoveTeamMemberArgs) (bool, error) {
	c := b.GetReqC(ctx)
	db := c.GetDB()

	userRole, err := getUserRoleInTeam(ctx, c, args.TeamID)
	if err != nil {
		return false, err
	}

	if userRole == nil {
		return false, errors.New("you do not have access to this team")
	}

	if *userRole == models.Owner {
		existingUser := &models.User{}
		err = db.Where("fb_uid = ?", args.UserUID).First(existingUser).Error
		if err != nil {
			return false, err
		}

		teamMember := &models.TeamMember{}
		err := db.Model(&models.TeamMember{}).Where("team_id = ? AND user_id = ?", args.TeamID, existingUser.ID).First(teamMember).Error
		if err != nil {
			return false, err
		}

		err = db.Delete(teamMember).Error
		if err != nil {
			return false, err
		}

		go func() {
			teamSubscriptions.EnsureChannel(teamMember.TeamID)

			teamSubscriptions.Subscriptions[teamMember.TeamID].Lock.Lock()
			defer teamSubscriptions.Subscriptions[teamMember.TeamID].Lock.Unlock()
			for i := range teamSubscriptions.Subscriptions[teamMember.TeamID].TeamMemberRemoved {
				teamSubscriptions.Subscriptions[teamMember.TeamID].TeamMemberRemoved[i] <- graphql.ID(existingUser.FBUID)
			}
		}()

		return true, nil
	}

	return false, errors.New("you do not have access to remove a team member on this team")
}

type UpdateTeamMemberRoleArgs struct {
	NewRole models.TeamMemberRole
	TeamID  graphql.ID
	UserUID graphql.ID
}

func (b *BaseQuery) UpdateTeamMemberRole(ctx context.Context, args *UpdateTeamMemberRoleArgs) (*TeamMemberResolver, error) {
	c := b.GetReqC(ctx)
	db := c.GetDB()

	userRole, err := getUserRoleInTeam(ctx, c, args.TeamID)
	if err != nil {
		return nil, err
	}

	if userRole == nil {
		return nil, errors.New("you do not have access to this team")
	}

	if *userRole == models.Owner {
		existingUser := &models.User{}
		err = db.Where("fb_uid = ?", args.UserUID).First(existingUser).Error
		if err != nil {
			return nil, err
		}

		teamMember := &models.TeamMember{}
		err := db.Model(&models.TeamMember{}).Where("team_id = ? AND user_id = ?", args.TeamID, existingUser.ID).First(teamMember).Error
		if err != nil {
			return nil, err
		}

		teamMember.Role = args.NewRole
		err = db.Save(teamMember).Error
		if err != nil {
			return nil, err
		}

		resolver, err := NewTeamMemberResolver(c, teamMember)
		if err != nil {
			return nil, err
		}

		go func() {
			teamSubscriptions.EnsureChannel(teamMember.TeamID)

			teamSubscriptions.Subscriptions[teamMember.TeamID].Lock.Lock()
			defer teamSubscriptions.Subscriptions[teamMember.TeamID].Lock.Unlock()
			for i := range teamSubscriptions.Subscriptions[teamMember.TeamID].TeamMemberUpdated {
				teamSubscriptions.Subscriptions[teamMember.TeamID].TeamMemberUpdated[i] <- resolver
			}
		}()

		return resolver, nil
	}

	return nil, errors.New("you do not have access to update a team member's role on this team")
}
