package resolvers

import (
	"context"
	"errors"
	"strconv"
	"strings"

	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
	"github.com/jerbob92/hoppscotch-backend/models"

	"github.com/graph-gophers/graphql-go"
	"gorm.io/gorm"
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
	id := graphql.ID(strconv.Itoa(int(r.team.ID)))
	return id, nil
}

func (r *TeamResolver) EditorsCount() (int32, error) {
	db := r.c.GetDB()

	ownerCount := int64(0)

	err := db.Model(&models.TeamMember{}).Where("team_id = ? AND role = ?", r.team.ID, models.Editor).Count(&ownerCount).Error
	if err != nil {
		return 0, err
	}

	return int32(ownerCount), nil
}

type TeamMembersArgs struct {
	Cursor *graphql.ID
}

func (r *TeamResolver) Members(args *TeamMembersArgs) ([]*TeamMemberResolver, error) {
	members := []*models.TeamMember{}
	db := r.c.GetDB()

	query := db.Model(&models.TeamMember{}).Where("team_id = ?", r.team.ID)
	if args.Cursor != nil && *args.Cursor != "" {
		query.Where("id > ?", args.Cursor)
	}

	err := query.Find(&members).Error
	if err != nil {
		return nil, err
	}

	teamMemberResolves := []*TeamMemberResolver{}
	for i := range members {
		newResolver, err := NewTeamMemberResolver(r.c, members[i])
		if err != nil {
			return nil, err
		}
		teamMemberResolves = append(teamMemberResolves, newResolver)
	}

	return teamMemberResolves, nil
}

func (r *TeamResolver) MyRole(ctx context.Context) (models.TeamMemberRole, error) {
	currentUser, err := r.c.GetUser(ctx)
	if err != nil {
		return models.Viewer, err
	}

	db := r.c.GetDB()

	existingTeamMember := &models.TeamMember{}
	err = db.Where("user_id = ? AND team_id = ?", currentUser.ID, r.team.ID).First(existingTeamMember).Error
	if err != nil {
		return models.Viewer, err
	}

	return existingTeamMember.Role, nil
}

func (r *TeamResolver) Name() (string, error) {
	return r.team.Name, nil
}

func (r *TeamResolver) OwnersCount() (int32, error) {
	db := r.c.GetDB()

	ownerCount := int64(0)

	err := db.Model(&models.TeamMember{}).Where("team_id = ? AND role = ?", r.team.ID, models.Owner).Count(&ownerCount).Error
	if err != nil {
		return 0, err
	}

	return int32(ownerCount), nil
}

func (r *TeamResolver) TeamInvitations() ([]*TeamInvitationResolver, error) {
	invitations := []*models.TeamInvitation{}
	db := r.c.GetDB()
	err := db.Model(&models.TeamInvitation{}).Where("team_id = ?", r.team.ID).Find(&invitations).Error
	if err != nil {
		return nil, err
	}

	teamInvitationResolvers := []*TeamInvitationResolver{}
	for i := range invitations {
		newResolver, err := NewTeamInvitationResolver(r.c, invitations[i])
		if err != nil {
			return nil, err
		}
		teamInvitationResolvers = append(teamInvitationResolvers, newResolver)
	}

	return teamInvitationResolvers, nil
}

func (r *TeamResolver) TeamMembers() ([]*TeamMemberResolver, error) {
	members := []*models.TeamMember{}
	db := r.c.GetDB()
	err := db.Model(&models.TeamMember{}).Where("team_id = ?", r.team.ID).Find(&members).Error
	if err != nil {
		return nil, err
	}

	teamMemberResolves := []*TeamMemberResolver{}
	for i := range members {
		newResolver, err := NewTeamMemberResolver(r.c, members[i])
		if err != nil {
			return nil, err
		}
		teamMemberResolves = append(teamMemberResolves, newResolver)
	}

	return teamMemberResolves, nil
}

func (r *TeamResolver) ViewersCount() (int32, error) {
	db := r.c.GetDB()

	ownerCount := int64(0)

	err := db.Model(&models.TeamMember{}).Where("team_id = ? AND role = ?", r.team.ID, models.Viewer).Count(&ownerCount).Error
	if err != nil {
		return 0, err
	}

	return int32(ownerCount), nil
}

type MyTeamsArgs struct {
	Cursor *graphql.ID
}

func (b *BaseQuery) MyTeams(ctx context.Context, args *MyTeamsArgs) ([]*TeamResolver, error) {
	c := b.GetReqC(ctx)
	currentUser, err := c.GetUser(ctx)
	if err != nil {
		c.LogErr(err)
		return nil, err
	}

	db := c.GetDB()
	teams := []*models.Team{}
	query := db.Model(&models.Team{}).Joins("JOIN `team_members` ON `team_members`.`team_id` = `teams`.`id`").Where("`team_members`.`user_id` = ?", currentUser.ID)
	if args.Cursor != nil && *args.Cursor != "" {
		query.Where("id > ?", args.Cursor)
	}

	err = query.Find(&teams).Error
	if err != nil {
		return nil, err
	}

	teamResolvers := []*TeamResolver{}
	for i := range teams {
		newResolver, err := NewTeamResolver(c, teams[i])
		if err != nil {
			return nil, err
		}
		teamResolvers = append(teamResolvers, newResolver)
	}

	return teamResolvers, nil
}

type RequestArg struct {
	RequestID graphql.ID
}

func (b *BaseQuery) Request(ctx context.Context, args *RequestArg) (*TeamRequestResolver, error) {
	c := b.GetReqC(ctx)
	db := c.GetDB()
	request := &models.TeamRequest{}
	err := db.Model(&models.TeamRequest{}).Where("id = ?", args.RequestID).First(request).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("you do not have access to this request")
	}
	if err != nil {
		return nil, err
	}

	userRole, err := getUserRoleInTeam(ctx, c, request.TeamID)
	if err != nil {
		return nil, err
	}

	if userRole == nil {
		return nil, errors.New("you do not have access to this request")
	}

	return NewTeamRequestResolver(c, request)
}

type RootCollectionsOfTeamArgs struct {
	Cursor *graphql.ID
	TeamID graphql.ID
}

func (b *BaseQuery) RootCollectionsOfTeam(ctx context.Context, args *RootCollectionsOfTeamArgs) ([]*TeamCollectionResolver, error) {
	c := b.GetReqC(ctx)
	userRole, err := getUserRoleInTeam(ctx, c, args.TeamID)
	if err != nil {
		return nil, err
	}
	if userRole == nil {
		return nil, errors.New("user not in team")
	}

	db := c.GetDB()
	teamCollections := []*models.TeamCollection{}
	query := db.Model(&models.TeamCollection{}).Where("team_id = ? AND parent_id = ?", args.TeamID, 0)
	if args.Cursor != nil && *args.Cursor != "" {
		query.Where("id > ?", args.Cursor)
	}
	err = query.Find(&teamCollections).Error
	if err != nil {
		return nil, err
	}

	teamCollectionResolvers := []*TeamCollectionResolver{}
	for i := range teamCollections {
		newResolver, err := NewTeamCollectionResolver(c, teamCollections[i])
		if err != nil {
			return nil, err
		}
		teamCollectionResolvers = append(teamCollectionResolvers, newResolver)
	}

	return teamCollectionResolvers, nil
}

type SearchForRequestArgs struct {
	Cursor     *graphql.ID
	SearchTerm string
	TeamID     graphql.ID
}

func (b *BaseQuery) SearchForRequest(ctx context.Context, args *SearchForRequestArgs) ([]*TeamRequestResolver, error) {
	c := b.GetReqC(ctx)
	userRole, err := getUserRoleInTeam(ctx, c, args.TeamID)
	if err != nil {
		return nil, err
	}
	if userRole == nil {
		return nil, errors.New("user not in team")
	}

	db := c.GetDB()
	teamRequests := []*models.TeamRequest{}
	args.SearchTerm = strings.Replace(args.SearchTerm, "%", "\\%", -1)
	args.SearchTerm = strings.Replace(args.SearchTerm, "_", "\\_", -1)
	args.SearchTerm = "%" + args.SearchTerm + "%"
	query := db.Model(&models.TeamRequest{}).Where("team_id = ? AND title LIKE ?", args.TeamID, args.SearchTerm)
	if args.Cursor != nil && *args.Cursor != "" {
		query.Where("id > ?", args.Cursor)
	}
	err = query.Find(&teamRequests).Error
	if err != nil {
		return nil, err
	}

	teamRequestResolvers := []*TeamRequestResolver{}
	for i := range teamRequests {
		newResolver, err := NewTeamRequestResolver(c, teamRequests[i])
		if err != nil {
			return nil, err
		}
		teamRequestResolvers = append(teamRequestResolvers, newResolver)
	}

	return teamRequestResolvers, nil
}

type TeamArgs struct {
	TeamID graphql.ID
}

func (b *BaseQuery) Team(ctx context.Context, args *TeamArgs) (*TeamResolver, error) {
	c := b.GetReqC(ctx)
	currentUser, err := c.GetUser(ctx)
	if err != nil {
		c.LogErr(err)
		return nil, err
	}

	db := c.GetDB()
	membership := &models.TeamMember{}
	err = db.Model(&models.TeamMember{}).Where("user_id = ? AND team_id = ?", currentUser.ID, args.TeamID).Preload("Team").First(membership).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("you do not have access to this team")
	}
	if err != nil {
		return nil, err
	}

	return NewTeamResolver(c, &membership.Team)
}

type CreateTeamArgs struct {
	Name string
}

func (b *BaseQuery) CreateTeam(ctx context.Context, args *CreateTeamArgs) (*TeamResolver, error) {
	c := b.GetReqC(ctx)
	currentUser, err := c.GetUser(ctx)
	if err != nil {
		c.LogErr(err)
		return nil, err
	}

	db := c.GetDB()
	newTeam := &models.Team{
		Name: args.Name,
	}

	err = db.Create(newTeam).Error
	if err != nil {
		return nil, err
	}

	newTeamMember := &models.TeamMember{
		TeamID: newTeam.ID,
		UserID: currentUser.ID,
		Role:   models.Owner,
	}

	err = db.Create(newTeamMember).Error
	if err != nil {
		return nil, err
	}

	return NewTeamResolver(c, newTeam)
}

type DeleteTeamArgs struct {
	TeamID graphql.ID
}

func (b *BaseQuery) DeleteTeam(ctx context.Context, args *DeleteTeamArgs) (bool, error) {
	c := b.GetReqC(ctx)

	currentUser, err := c.GetUser(ctx)
	if err != nil {
		c.LogErr(err)
		return false, err
	}

	db := c.GetDB()
	existingTeamMember := &models.TeamMember{}
	err = db.Where("user_id = ? AND team_id = ?", currentUser.ID, args.TeamID).First(existingTeamMember).Error
	if err != nil {
		return false, err
	}

	if existingTeamMember.Role != models.Owner {
		return false, errors.New("no access to delete")
	}

	// Cleanup related records.
	err = db.Delete(&models.TeamCollection{}, "team_id = ?", args.TeamID).Error
	if err != nil {
		return false, err
	}
	err = db.Delete(&models.TeamInvitation{}, "team_id = ?", args.TeamID).Error
	if err != nil {
		return false, err
	}
	err = db.Delete(&models.TeamMember{}, "team_id = ?", args.TeamID).Error
	if err != nil {
		return false, err
	}
	err = db.Delete(&models.TeamRequest{}, "team_id = ?", args.TeamID).Error
	if err != nil {
		return false, err
	}

	err = db.Delete(&models.Team{}, "id = ?", args.TeamID).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

type LeaveTeamArgs struct {
	TeamID graphql.ID
}

func (b *BaseQuery) LeaveTeam(ctx context.Context, args *LeaveTeamArgs) (bool, error) {
	c := b.GetReqC(ctx)

	currentUser, err := c.GetUser(ctx)
	if err != nil {
		c.LogErr(err)
		return false, err
	}

	db := c.GetDB()
	existingTeamMember := &models.TeamMember{}
	err = db.Where("user_id = ? AND team_id = ?", currentUser.ID, args.TeamID).First(existingTeamMember).Error
	if err != nil {
		return false, err
	}

	err = db.Delete(&models.TeamMember{}, "user_id = ? AND team_id = ?", currentUser.ID, args.TeamID).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

type RenameTeamArgs struct {
	NewName string
	TeamID  graphql.ID
}

func (b *BaseQuery) RenameTeam(ctx context.Context, args *RenameTeamArgs) (*TeamResolver, error) {
	c := b.GetReqC(ctx)

	currentUser, err := c.GetUser(ctx)
	if err != nil {
		c.LogErr(err)
		return nil, err
	}

	db := c.GetDB()
	existingTeamMember := &models.TeamMember{}
	err = db.Where("user_id = ? AND team_id = ?", currentUser.ID, args.TeamID).Preload("Team").First(existingTeamMember).Error
	if err != nil {
		return nil, err
	}

	if existingTeamMember.Role != models.Owner {
		return nil, errors.New("no access to rename")
	}

	existingTeamMember.Team.Name = args.NewName

	err = db.Save(&existingTeamMember.Team).Error
	if err != nil {
		return nil, err
	}

	return NewTeamResolver(c, &existingTeamMember.Team)
}
