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
	id := graphql.ID(strconv.Itoa(int(r.team_request.ID)))
	return id, nil
}

func (r *TeamRequestResolver) Collection() (*TeamCollectionResolver, error) {
	db := r.c.GetDB()
	collection := &models.TeamCollection{}
	err := db.Model(&models.TeamCollection{}).Where("id = ?", r.team_request.TeamCollectionID).First(collection).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("you do not have access to this collection")
	}
	if err != nil {
		return nil, err
	}
	return NewTeamCollectionResolver(r.c, collection)
}

func (r *TeamRequestResolver) CollectionID() (graphql.ID, error) {
	return graphql.ID(strconv.Itoa(int(r.team_request.TeamCollectionID))), nil
}

func (r *TeamRequestResolver) Request() (string, error) {
	return r.team_request.Request, nil
}

func (r *TeamRequestResolver) Team() (*TeamResolver, error) {
	db := r.c.GetDB()
	team := &models.Team{}
	err := db.Model(&models.Team{}).Where("id = ?", r.team_request.TeamID).First(team).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("you do not have access to this team")
	}
	if err != nil {
		return nil, err
	}
	return NewTeamResolver(r.c, team)
}

func (r *TeamRequestResolver) TeamID() (graphql.ID, error) {
	return graphql.ID(strconv.Itoa(int(r.team_request.TeamID))), nil
}

func (r *TeamRequestResolver) Title() (string, error) {
	return r.team_request.Title, nil
}

type DeleteRequestArgs struct {
	RequestID graphql.ID
}

func (b *BaseQuery) DeleteRequest(ctx context.Context, args *DeleteRequestArgs) (bool, error) {
	c := b.GetReqC(ctx)
	db := c.GetDB()
	request := &models.TeamRequest{}
	err := db.Model(&models.TeamRequest{}).Where("id = ?", args.RequestID).First(request).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return false, errors.New("you do not have access to this request")
	}
	if err != nil {
		return false, err
	}

	userRole, err := getUserRoleInTeam(ctx, c, request.TeamID)
	if err != nil {
		return false, err
	}

	if userRole == nil {
		return false, errors.New("you do not have access to delete to this request")
	}

	if *userRole == models.Owner || *userRole == models.Editor {
		err := db.Delete(request).Error
		if err != nil {
			return false, err
		}

		go bus.Publish("team:"+strconv.Itoa(int(request.TeamID))+":requests:deleted", graphql.ID(strconv.Itoa(int(request.ID))))

		return true, nil
	}

	return false, errors.New("you are not allowed to delete a request in this team")
}

type MoveRequestArgs struct {
	DestCollID graphql.ID
	RequestID  graphql.ID
}

func (b *BaseQuery) MoveRequest(ctx context.Context, args *MoveRequestArgs) (*TeamRequestResolver, error) {
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
		return nil, errors.New("you do not have move to this request")
	}

	collection := &models.TeamCollection{}
	err = db.Model(&models.TeamCollection{}).Where("id = ?", args.DestCollID).First(collection).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("you do not have access to this collection")
	}
	if err != nil {
		return nil, err
	}

	targetUserRole, err := getUserRoleInTeam(ctx, c, collection.TeamID)
	if err != nil {
		return nil, err
	}

	if targetUserRole == nil {
		return nil, errors.New("you do not have access to move to this request")
	}

	if (*userRole == models.Owner || *userRole == models.Editor) && (*targetUserRole == models.Owner || *targetUserRole == models.Editor) {
		teamChanged := false
		oldTeamID := request.TeamID
		newTeamID := collection.TeamID
		if collection.TeamID != request.TeamID {
			teamChanged = true
		}

		request.TeamCollectionID = collection.ID
		request.TeamID = collection.TeamID
		err := db.Save(request).Error
		if err != nil {
			return nil, err
		}

		resolver, err := NewTeamRequestResolver(c, request)
		if err != nil {
			return nil, err
		}

		if teamChanged {
			go bus.Publish("team:"+strconv.Itoa(int(oldTeamID))+":requests:deleted", graphql.ID(strconv.Itoa(int(request.ID))))
			go bus.Publish("team:"+strconv.Itoa(int(newTeamID))+":requests:added", graphql.ID(strconv.Itoa(int(request.ID))))
		} else {
			go bus.Publish("team:"+strconv.Itoa(int(newTeamID))+":requests:updated", resolver)
		}

		return resolver, nil
	}

	return nil, errors.New("you are not allowed to delete a request in this team")
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
		return nil, errors.New("you do not have access to update to this request")
	}

	if *userRole == models.Owner || *userRole == models.Editor {
		if args.Data.Title != nil {
			request.Title = *args.Data.Title
		}
		if args.Data.Request != nil {
			request.Request = *args.Data.Request
		}
		err := db.Save(request).Error
		if err != nil {
			return nil, err
		}

		requestResolver, err := NewTeamRequestResolver(c, request)
		if err != nil {
			return nil, err
		}

		go bus.Publish("team:"+strconv.Itoa(int(request.TeamID))+":requests:updated", requestResolver)

		return requestResolver, nil
	}

	return nil, errors.New("you are not allowed to update a request in this team")
}
