package resolvers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"strconv"

	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"
	"github.com/jerbob92/hoppscotch-backend/models"

	"github.com/graph-gophers/graphql-go"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
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
	id := graphql.ID(r.team_invitation.Code)
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
	c := b.GetReqC(ctx)
	db := c.GetDB()

	invite := &models.TeamInvitation{}
	err := db.Model(&models.TeamInvitation{}).Where("code = ?", args.InviteID).First(invite).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("team_invite/no_invite_found")
	}
	if err != nil {
		return nil, err
	}

	currentUser, err := c.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	if invite.InviteeEmail != currentUser.Email {
		return nil, errors.New("team_invite/email_do_not_match")
	}

	return NewTeamInvitationResolver(c, invite)
}

type AcceptTeamInvitationArgs struct {
	InviteID graphql.ID
}

func (b *BaseQuery) AcceptTeamInvitation(ctx context.Context, args *AcceptTeamInvitationArgs) (*TeamMemberResolver, error) {
	c := b.GetReqC(ctx)
	db := c.GetDB()

	invite := &models.TeamInvitation{}
	err := db.Model(&models.TeamInvitation{}).Where("code = ?", args.InviteID).First(invite).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("you do not have access to this invite")
	}
	if err != nil {
		return nil, err
	}

	currentUser, err := c.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	if invite.InviteeEmail != currentUser.Email {
		return nil, errors.New("team_invite/email_do_not_match")
	}

	newTeamMember := &models.TeamMember{
		TeamID: invite.TeamID,
		UserID: currentUser.ID,
		Role:   invite.InviteeRole,
	}

	err = db.Create(newTeamMember).Error
	if err != nil {
		return nil, err
	}

	err = db.Delete(invite).Error
	if err != nil {
		return nil, err
	}

	resolver, err := NewTeamMemberResolver(c, newTeamMember)
	if err != nil {
		return nil, err
	}

	go func() {
		teamSubscriptions.EnsureChannel(invite.TeamID)

		teamSubscriptions.Subscriptions[invite.TeamID].Lock.Lock()
		defer teamSubscriptions.Subscriptions[invite.TeamID].Lock.Unlock()
		for i := range teamSubscriptions.Subscriptions[invite.TeamID].TeamMemberAdded {
			teamSubscriptions.Subscriptions[invite.TeamID].TeamMemberAdded[i] <- resolver
		}
	}()

	return resolver, nil
}

type AddTeamMemberByEmailArgs struct {
	TeamID    graphql.ID
	UserEmail string
	UserRole  models.TeamMemberRole
}

func (b *BaseQuery) AddTeamMemberByEmail(ctx context.Context, args *AddTeamMemberByEmailArgs) (*TeamMemberResolver, error) {
	// This doesn't seem to be used (anymore).
	return nil, nil
}

type CreateTeamInvitationArgs struct {
	InviteeEmail string
	InviteeRole  models.TeamMemberRole
	TeamID       graphql.ID
}

func (b *BaseQuery) CreateTeamInvitation(ctx context.Context, args *CreateTeamInvitationArgs) (*TeamInvitationResolver, error) {
	c := b.GetReqC(ctx)

	userRole, err := getUserRoleInTeam(ctx, c, args.TeamID)
	if err != nil {
		return nil, err
	}

	if userRole == nil {
		return nil, errors.New("you do not have access to this team")
	}

	if *userRole != models.Owner {
		return nil, errors.New("you do not have access to this team")
	}

	currentUser, err := c.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	db := c.GetDB()

	parsedTeamID, _ := strconv.Atoi(string(args.TeamID))
	invite := &models.TeamInvitation{
		TeamID:       uint(parsedTeamID),
		UserID:       currentUser.ID,
		InviteeRole:  args.InviteeRole,
		InviteeEmail: args.InviteeEmail,
		Code:         RandString(32),
	}

	err = db.Save(invite).Error
	if err != nil {
		return nil, err
	}

	name := "A user"
	if currentUser.DisplayName != "" {
		name = currentUser.DisplayName
	} else if currentUser.Email != "" {
		name = currentUser.Email
	}

	team := &models.Team{}
	err = db.Model(&models.Team{}).Where("id = ?", args.TeamID).First(team).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, errors.New("you do not have access to this team")
	}
	if err != nil {
		return nil, err
	}

	from := fmt.Sprintf("%s <%s>", viper.GetString("smtp.from.name"), viper.GetString("smtp.from.email"))

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", invite.InviteeEmail)

	joinLink := viper.GetString("frontend_domain") + "/join-team?id=" + invite.Code

	templateVariables := struct {
		InvitingUserName string
		TeamName         string
		JoinLink         string
	}{
		InvitingUserName: name,
		TeamName:         team.Name,
		JoinLink:         joinLink,
	}

	subjectTemplate := template.New("Subject")
	subjectTemplate, err = subjectTemplate.Parse(viper.GetString("mailTemplates.teamInvite.subject"))
	if err != nil {
		return nil, err
	}

	var subject bytes.Buffer
	err = subjectTemplate.Execute(&subject, templateVariables)
	if err != nil {
		return nil, err
	}

	m.SetHeader("Subject", subject.String())

	bodyTemplate := template.New("Body")
	bodyTemplate, err = bodyTemplate.Parse(viper.GetString("mailTemplates.teamInvite.body"))
	if err != nil {
		return nil, err
	}

	var body bytes.Buffer
	err = bodyTemplate.Execute(&body, templateVariables)
	if err != nil {
		return nil, err
	}

	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(viper.GetString("smtp.host"), viper.GetInt("smtp.port"), viper.GetString("smtp.username"), viper.GetString("smtp.password"))
	if err := d.DialAndSend(m); err != nil {
		return nil, err
	}

	resolver, err := NewTeamInvitationResolver(c, invite)
	if err != nil {
		return nil, err
	}

	go func() {
		teamSubscriptions.EnsureChannel(invite.TeamID)

		teamSubscriptions.Subscriptions[invite.TeamID].Lock.Lock()
		defer teamSubscriptions.Subscriptions[invite.TeamID].Lock.Unlock()
		for i := range teamSubscriptions.Subscriptions[invite.TeamID].TeamInvitationRemoved {
			teamSubscriptions.Subscriptions[invite.TeamID].TeamInvitationAdded[i] <- resolver
		}
	}()

	return resolver, nil
}

type RevokeTeamInvitationArgs struct {
	InviteID graphql.ID
}

func (b *BaseQuery) RevokeTeamInvitation(ctx context.Context, args *RevokeTeamInvitationArgs) (bool, error) {
	c := b.GetReqC(ctx)
	db := c.GetDB()

	invite := &models.TeamInvitation{}
	err := db.Model(&models.TeamInvitation{}).Where("code = ?", args.InviteID).First(invite).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return false, errors.New("you do not have access to this invite")
	}
	if err != nil {
		return false, err
	}

	userRole, err := getUserRoleInTeam(ctx, c, invite.TeamID)
	if err != nil {
		return false, err
	}

	if userRole == nil {
		return false, errors.New("you do not have access to this invite")
	}

	if *userRole != models.Owner {
		return false, errors.New("you do not have access to this invite")
	}

	err = db.Delete(invite).Error
	if err != nil {
		return false, err
	}

	go func() {
		teamSubscriptions.EnsureChannel(invite.TeamID)

		teamSubscriptions.Subscriptions[invite.TeamID].Lock.Lock()
		defer teamSubscriptions.Subscriptions[invite.TeamID].Lock.Unlock()
		for i := range teamSubscriptions.Subscriptions[invite.TeamID].TeamInvitationRemoved {
			teamSubscriptions.Subscriptions[invite.TeamID].TeamInvitationRemoved[i] <- graphql.ID(invite.Code)
		}
	}()

	return true, nil
}
