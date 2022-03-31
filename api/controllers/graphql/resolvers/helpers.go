package resolvers

import (
	"context"

	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
	"github.com/jerbob92/hoppscotch-backend/models"

	"gorm.io/gorm"
)

func getUserRoleInTeam(ctx context.Context, c *graphql_context.Context, teamID interface{}) (*models.TeamMemberRole, error) {
	currentUser, err := c.GetUser(ctx)
	if err != nil {
		c.LogErr(err)
		return nil, err
	}

	db := c.GetDB()

	existingTeamMember := &models.TeamMember{}
	err = db.Where("user_id = ? AND team_id = ?", currentUser.ID, teamID).Preload("Team").First(existingTeamMember).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &existingTeamMember.Role, nil
}
