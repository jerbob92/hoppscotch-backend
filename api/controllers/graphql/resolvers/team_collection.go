package resolvers

import (
	"context"
	"github.com/graph-gophers/graphql-go"
	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
	"github.com/jerbob92/hoppscotch-backend/models"
	"strconv"
)

type TeamCollectionResolver struct {
	c               *graphql_context.Context
	team_collection *models.TeamCollection
}

func NewTeamCollectionResolver(c *graphql_context.Context, team_collection *models.TeamCollection) (*TeamCollectionResolver, error) {
	if team_collection == nil {
		return nil, nil
	}

	return &TeamCollectionResolver{c: c, team_collection: team_collection}, nil
}

func (r *TeamCollectionResolver) ID() (graphql.ID, error) {
	id := graphql.ID(strconv.FormatInt(r.team_collection.ID, 10))
	return id, nil
}

func (r *TeamCollectionResolver) Parent() (*TeamCollectionResolver, error) {
	return nil, nil
}

func (r *TeamCollectionResolver) Team() (*TeamResolver, error) {
	return nil, nil
}

type TeamCollectionChildrenArgs struct {
	Cursor *string
}

func (r *TeamCollectionResolver) Children(args *TeamCollectionChildrenArgs) ([]*TeamCollectionResolver, error) {
	return nil, nil
}

func (r *TeamCollectionResolver) Title() (string, error) {
	return "", nil
}

type CollectionArgs struct {
	CollectionID graphql.ID
}

func (b *BaseQuery) Collection(ctx context.Context, args *CollectionArgs) (*TeamCollectionResolver, error) {
	// @todo: implement me
	return nil, nil
}

type CollectionsOfTeamArgs struct {
	Cursor *graphql.ID
	TeamID graphql.ID
}

func (b *BaseQuery) CollectionsOfTeam(ctx context.Context, args *CollectionsOfTeamArgs) ([]*TeamCollectionResolver, error) {
	// @todo: implement me
	return nil, nil
}

type ExportCollectionsToJSONArgs struct {
	TeamID graphql.ID
}

func (b *BaseQuery) ExportCollectionsToJSON(ctx context.Context, args *ExportCollectionsToJSONArgs) (string, error) {
	// @todo: implement me
	return "nil", nil
}

type RequestsInCollectionArgs struct {
	CollectionID graphql.ID
	Cursor       *graphql.ID
}

func (b *BaseQuery) RequestsInCollection(ctx context.Context, args *RequestsInCollectionArgs) ([]*TeamRequestResolver, error) {
	// @todo: implement me
	return nil, nil
}

type CreateChildCollectionArgs struct {
	ChildTitle   string
	CollectionID graphql.ID
}

func (b *BaseQuery) CreateChildCollection(ctx context.Context, args *CreateChildCollectionArgs) (*TeamCollectionResolver, error) {
	// @todo: implement me
	return nil, nil
}

type CreateTeamRequestInput struct {
	Request string
	TeamID  graphql.ID
	Title   string
}

type CreateRequestInCollectionArgs struct {
	CollectionID graphql.ID
	Data         CreateTeamRequestInput
}

func (b *BaseQuery) CreateRequestInCollection(ctx context.Context, args *CreateRequestInCollectionArgs) (*TeamRequestResolver, error) {
	// @todo: implement me
	return nil, nil
}

type CreateRootCollectionArgs struct {
	TeamID graphql.ID
	Title  string
}

func (b *BaseQuery) CreateRootCollection(ctx context.Context, args *CreateRootCollectionArgs) (*TeamCollectionResolver, error) {
	// @todo: implement me
	return nil, nil
}

type DeleteCollectionArgs struct {
	CollectionID graphql.ID
}

func (b *BaseQuery) DeleteCollection(ctx context.Context, args *DeleteCollectionArgs) (bool, error) {
	// @todo: implement me
	return false, nil
}

type ImportCollectionFromUserFirestoreArgs struct {
	FBCollectionPath   string
	ParentCollectionID *graphql.ID
	TeamID             graphql.ID
}

func (b *BaseQuery) ImportCollectionFromUserFirestore(ctx context.Context, args *ImportCollectionFromUserFirestoreArgs) (*TeamCollectionResolver, error) {
	// @todo: implement me
	return nil, nil
}

type ImportCollectionsFromJSONArgs struct {
	JSONString         string
	ParentCollectionID *graphql.ID
	TeamID             graphql.ID
}

func (b *BaseQuery) ImportCollectionsFromJSON(ctx context.Context, args *ImportCollectionsFromJSONArgs) (bool, error) {
	// @todo: implement me
	return false, nil
}

type RenameCollectionArgs struct {
	CollectionID graphql.ID
	NewTitle     string
}

func (b *BaseQuery) RenameCollection(ctx context.Context, args *RenameCollectionArgs) (*TeamCollectionResolver, error) {
	// @todo: implement me
	return nil, nil
}

type ReplaceCollectionsWithJSONArgs struct {
	JSONString         string
	ParentCollectionID *graphql.ID
	TeamID             graphql.ID
}

func (b *BaseQuery) ReplaceCollectionsWithJSON(ctx context.Context, args *ReplaceCollectionsWithJSONArgs) (bool, error) {
	// @todo: implement me
	return false, nil
}

type SubscriptionArgs struct {
	TeamID graphql.ID
}

func (b *BaseQuery) TeamCollectionAdded(ctx context.Context, args *SubscriptionArgs) (<-chan *TeamCollectionResolver, error) {
	// @todo: implement me
	return nil, nil
}

func (b *BaseQuery) TeamCollectionRemoved(ctx context.Context, args *SubscriptionArgs) (<-chan graphql.ID, error) {
	// @todo: implement me
	return nil, nil
}

func (b *BaseQuery) TeamCollectionUpdated(ctx context.Context, args *SubscriptionArgs) (<-chan *TeamCollectionResolver, error) {
	// @todo: implement me
	return nil, nil
}

func (b *BaseQuery) TeamInvitationAdded(ctx context.Context, args *SubscriptionArgs) (<-chan *TeamInvitationResolver, error) {
	// @todo: implement me
	return nil, nil
}

func (b *BaseQuery) TeamInvitationRemoved(ctx context.Context, args *SubscriptionArgs) (<-chan graphql.ID, error) {
	// @todo: implement me
	return nil, nil
}

func (b *BaseQuery) TeamMemberAdded(ctx context.Context, args *SubscriptionArgs) (<-chan *TeamMemberResolver, error) {
	// @todo: implement me
	return nil, nil
}

func (b *BaseQuery) TeamMemberRemoved(ctx context.Context, args *SubscriptionArgs) (<-chan graphql.ID, error) {
	// @todo: implement me
	return nil, nil
}

func (b *BaseQuery) TeamMemberUpdated(ctx context.Context, args *SubscriptionArgs) (<-chan *TeamMemberResolver, error) {
	// @todo: implement me
	return nil, nil
}

func (b *BaseQuery) TeamRequestAdded(ctx context.Context, args *SubscriptionArgs) (<-chan *TeamRequestResolver, error) {
	// @todo: implement me
	return nil, nil
}

func (b *BaseQuery) TeamRequestDeleted(ctx context.Context, args *SubscriptionArgs) (<-chan graphql.ID, error) {
	// @todo: implement me
	return nil, nil
}

func (b *BaseQuery) TeamRequestUpdated(ctx context.Context, args *SubscriptionArgs) (<-chan *TeamRequestResolver, error) {
	// @todo: implement me
	return nil, nil
}
