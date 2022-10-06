package resolvers

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"strconv"

	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
	"github.com/jerbob92/hoppscotch-backend/models"

	"github.com/graph-gophers/graphql-go"
)

type TeamEnvironmentResolver struct {
	c                *graphql_context.Context
	team_environment *models.TeamEnvironment
}

func NewTeamEnvironmentResolver(c *graphql_context.Context, team_environment *models.TeamEnvironment) (*TeamEnvironmentResolver, error) {
	if team_environment == nil {
		return nil, nil
	}

	return &TeamEnvironmentResolver{c: c, team_environment: team_environment}, nil
}

func (r *TeamEnvironmentResolver) ID() (graphql.ID, error) {
	id := graphql.ID(strconv.Itoa(int(r.team_environment.ID)))
	return id, nil
}

func (r *TeamEnvironmentResolver) Variables() (string, error) {
	return r.team_environment.Variables, nil
}

func (r *TeamEnvironmentResolver) TeamID() (graphql.ID, error) {
	return graphql.ID(strconv.Itoa(int(r.team_environment.TeamID))), nil
}

func (r *TeamEnvironmentResolver) Name() (string, error) {
	return r.team_environment.Name, nil
}

type CreateTeamEnvironmentRequestArgs struct {
	Name      string
	TeamID    graphql.ID
	Variables string
}

func (b *BaseQuery) CreateTeamEnvironment(ctx context.Context, args *CreateTeamEnvironmentRequestArgs) (*TeamEnvironmentResolver, error) {
	c := b.GetReqC(ctx)
	db := c.GetDB()

	userRole, err := getUserRoleInTeam(ctx, c, args.TeamID)
	if err != nil {
		return nil, err
	}

	if userRole == nil {
		return nil, errors.New("you are not allowed to create an environment in this team")
	}

	teamID, _ := strconv.Atoi(string(args.TeamID))

	if *userRole == models.Owner || *userRole == models.Editor {
		newTeamEnvironment := &models.TeamEnvironment{
			TeamID:    uint(teamID),
			Name:      args.Name,
			Variables: args.Variables,
		}

		err := db.Save(newTeamEnvironment).Error
		if err != nil {
			return nil, err
		}

		resolver, err := NewTeamEnvironmentResolver(c, newTeamEnvironment)
		if err != nil {
			return nil, err
		}

		go func() {
			teamSubscriptions.EnsureChannel(uint(teamID))

			teamSubscriptions.Subscriptions[uint(teamID)].Lock.Lock()
			defer teamSubscriptions.Subscriptions[uint(teamID)].Lock.Unlock()
			for i := range teamSubscriptions.Subscriptions[uint(teamID)].TeamEnvironmentCreated {
				teamSubscriptions.Subscriptions[uint(teamID)].TeamEnvironmentCreated[i] <- resolver
			}
		}()

		return resolver, nil
	}

	return nil, errors.New("you are not allowed to create an environment in this team")
}

type DeleteTeamEnvironmentRequestArgs struct {
	ID graphql.ID
}

func (b *BaseQuery) DeleteTeamEnvironment(ctx context.Context, args *DeleteTeamEnvironmentRequestArgs) (bool, error) {
	c := b.GetReqC(ctx)
	db := c.GetDB()

	teamEnvironment := &models.TeamEnvironment{}
	err := db.Model(&models.TeamEnvironment{}).Where("id = ?", args.ID).First(teamEnvironment).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return false, errors.New("you do not have access to this team")
	}
	if err != nil {
		return false, err
	}

	userRole, err := getUserRoleInTeam(ctx, c, teamEnvironment.TeamID)
	if err != nil {
		return false, err
	}

	if userRole == nil {
		return false, errors.New("you are not allowed to delete an environment in this team")
	}

	if *userRole == models.Owner || *userRole == models.Editor {
		err := db.Delete(teamEnvironment).Error
		if err != nil {
			return false, err
		}

		resolver, err := NewTeamEnvironmentResolver(c, teamEnvironment)
		if err != nil {
			return false, err
		}

		go func() {
			teamSubscriptions.EnsureChannel(teamEnvironment.TeamID)

			teamSubscriptions.Subscriptions[teamEnvironment.TeamID].Lock.Lock()
			defer teamSubscriptions.Subscriptions[teamEnvironment.TeamID].Lock.Unlock()
			for i := range teamSubscriptions.Subscriptions[teamEnvironment.TeamID].TeamEnvironmentDeleted {
				teamSubscriptions.Subscriptions[teamEnvironment.TeamID].TeamEnvironmentDeleted[i] <- resolver
			}
		}()

		return true, nil
	}

	return false, errors.New("you are not allowed to create an environment in this team")
}

type UpdateTeamEnvironmentRequestArgs struct {
	ID        graphql.ID
	Name      string
	Variables string
}

func (b *BaseQuery) UpdateTeamEnvironment(ctx context.Context, args *UpdateTeamEnvironmentRequestArgs) (*TeamEnvironmentResolver, error) {
	c := b.GetReqC(ctx)
	db := c.GetDB()

	teamEnvironment := &models.TeamEnvironment{}
	err := db.Model(&models.TeamEnvironment{}).Where("id = ?", args.ID).First(teamEnvironment).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("you do not have access to this team")
	}
	if err != nil {
		return nil, err
	}

	userRole, err := getUserRoleInTeam(ctx, c, teamEnvironment.TeamID)
	if err != nil {
		return nil, err
	}

	if userRole == nil {
		return nil, errors.New("you are not allowed to update an environment in this team")
	}

	if *userRole == models.Owner || *userRole == models.Editor {
		teamEnvironment.Name = args.Name
		teamEnvironment.Variables = args.Variables

		err := db.Save(teamEnvironment).Error
		if err != nil {
			return nil, err
		}

		resolver, err := NewTeamEnvironmentResolver(c, teamEnvironment)
		if err != nil {
			return nil, err
		}

		go func() {
			teamSubscriptions.EnsureChannel(teamEnvironment.TeamID)

			teamSubscriptions.Subscriptions[teamEnvironment.TeamID].Lock.Lock()
			defer teamSubscriptions.Subscriptions[teamEnvironment.TeamID].Lock.Unlock()
			for i := range teamSubscriptions.Subscriptions[teamEnvironment.TeamID].TeamEnvironmentCreated {
				teamSubscriptions.Subscriptions[teamEnvironment.TeamID].TeamEnvironmentCreated[i] <- resolver
			}
		}()

		return resolver, nil
	}

	return nil, errors.New("you are not allowed to update an environment in this team")
}

type DeleteAllVariablesFromTeamEnvironmentRequestArgs struct {
	ID graphql.ID
}

func (b *BaseQuery) DeleteAllVariablesFromTeamEnvironment(ctx context.Context, args *DeleteAllVariablesFromTeamEnvironmentRequestArgs) (*TeamEnvironmentResolver, error) {
	c := b.GetReqC(ctx)
	db := c.GetDB()

	teamEnvironment := &models.TeamEnvironment{}
	err := db.Model(&models.TeamEnvironment{}).Where("id = ?", args.ID).First(teamEnvironment).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("you do not have access to this team")
	}
	if err != nil {
		return nil, err
	}

	userRole, err := getUserRoleInTeam(ctx, c, teamEnvironment.TeamID)
	if err != nil {
		return nil, err
	}

	if userRole == nil {
		return nil, errors.New("you are not allowed to update an environment in this team")
	}

	if *userRole == models.Owner || *userRole == models.Editor {
		teamEnvironment.Variables = ""

		err := db.Save(teamEnvironment).Error
		if err != nil {
			return nil, err
		}

		resolver, err := NewTeamEnvironmentResolver(c, teamEnvironment)
		if err != nil {
			return nil, err
		}

		go func() {
			teamSubscriptions.EnsureChannel(teamEnvironment.TeamID)

			teamSubscriptions.Subscriptions[teamEnvironment.TeamID].Lock.Lock()
			defer teamSubscriptions.Subscriptions[teamEnvironment.TeamID].Lock.Unlock()
			for i := range teamSubscriptions.Subscriptions[teamEnvironment.TeamID].TeamEnvironmentCreated {
				teamSubscriptions.Subscriptions[teamEnvironment.TeamID].TeamEnvironmentCreated[i] <- resolver
			}
		}()

		return resolver, nil
	}

	return nil, errors.New("you are not allowed to update an environment in this team")
}

type CreateDuplicateEnvironmentRequestArgs struct {
	ID graphql.ID
}

func (b *BaseQuery) CreateDuplicateEnvironment(ctx context.Context, args *CreateDuplicateEnvironmentRequestArgs) (*TeamEnvironmentResolver, error) {
	c := b.GetReqC(ctx)
	db := c.GetDB()

	teamEnvironment := &models.TeamEnvironment{}
	err := db.Model(&models.TeamEnvironment{}).Where("id = ?", args.ID).First(teamEnvironment).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("you do not have access to this team")
	}
	if err != nil {
		return nil, err
	}

	userRole, err := getUserRoleInTeam(ctx, c, teamEnvironment.TeamID)
	if err != nil {
		return nil, err
	}

	if userRole == nil {
		return nil, errors.New("you are not allowed to duplicate an environment in this team")
	}

	if *userRole == models.Owner || *userRole == models.Editor {
		newTeamEnvironment := &models.TeamEnvironment{
			TeamID:    teamEnvironment.TeamID,
			Name:      teamEnvironment.Name,
			Variables: teamEnvironment.Variables,
		}

		err := db.Save(newTeamEnvironment).Error
		if err != nil {
			return nil, err
		}

		resolver, err := NewTeamEnvironmentResolver(c, newTeamEnvironment)
		if err != nil {
			return nil, err
		}

		go func() {
			teamSubscriptions.EnsureChannel(teamEnvironment.TeamID)

			teamSubscriptions.Subscriptions[teamEnvironment.TeamID].Lock.Lock()
			defer teamSubscriptions.Subscriptions[teamEnvironment.TeamID].Lock.Unlock()
			for i := range teamSubscriptions.Subscriptions[teamEnvironment.TeamID].TeamEnvironmentCreated {
				teamSubscriptions.Subscriptions[teamEnvironment.TeamID].TeamEnvironmentCreated[i] <- resolver
			}
		}()

		return resolver, nil
	}

	return nil, errors.New("you are not allowed to duplicate an environment in this team")
}
