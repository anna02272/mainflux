// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package things

import (
	"context"
	"time"

	"github.com/MainfluxLabs/mainflux/auth"
	"github.com/MainfluxLabs/mainflux/pkg/apiutil"
	"github.com/MainfluxLabs/mainflux/pkg/errors"
	protomfx "github.com/MainfluxLabs/mainflux/pkg/proto"
	"github.com/MainfluxLabs/mainflux/pkg/uuid"
)

const (
	Viewer = "viewer"
	Editor = "editor"
	Admin  = "admin"
	Owner  = "owner"
)

var (
	ErrProfileAssigned = errors.New("profile currently assigned to thing(s)")
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// CreateThings adds things to the user identified by the provided key.
	CreateThings(ctx context.Context, token string, things ...Thing) ([]Thing, error)

	// UpdateThing updates the thing identified by the provided ID, that
	// belongs to the user identified by the provided token.
	UpdateThing(ctx context.Context, token string, thing Thing) error

	// UpdateThingsMetadata updates the things metadata identified by the provided IDs, that
	// belongs to the user identified by the provided token.
	UpdateThingsMetadata(ctx context.Context, token string, things ...Thing) error

	// UpdateKey updates key value of the existing thing. A non-nil error is
	// returned to indicate operation failure.
	UpdateKey(ctx context.Context, token, id, key string) error

	// ViewThing retrieves data about the thing identified with the provided
	// ID, that belongs to the user identified by the provided key.
	ViewThing(ctx context.Context, token, id string) (Thing, error)

	// ListThings retrieves data about subset of things that belongs to the
	// user identified by the provided key.
	ListThings(ctx context.Context, token string, pm apiutil.PageMetadata) (ThingsPage, error)

	// ListThingsByOrg retrieves page of things that belong to an org identified by ID.
	ListThingsByOrg(ctx context.Context, token string, orgID string, pm apiutil.PageMetadata) (ThingsPage, error)

	// ListThingsByProfile retrieves data about subset of things that are
	// connected or not connected to specified profile and belong to the user identified by
	// the provided key.
	ListThingsByProfile(ctx context.Context, token, prID string, pm apiutil.PageMetadata) (ThingsPage, error)

	// BackupThingsByGroup retrieves all things for given group ID.
	BackupThingsByGroup(ctx context.Context, token string, groupID string) (ThingsBackup, error)

	// BackupThingsByOrg retrieves all things for given org ID.
	BackupThingsByOrg(ctx context.Context, token string, orgID string) (ThingsBackup, error)

	// RemoveThings removes the things identified with the provided IDs, that
	// belongs to the user identified by the provided key.
	RemoveThings(ctx context.Context, token string, id ...string) error

	// CreateProfiles adds profiles to the user identified by the provided key.
	CreateProfiles(ctx context.Context, token string, profiles ...Profile) ([]Profile, error)

	// UpdateProfile updates the profile identified by the provided ID, that
	// belongs to the user identified by the provided key.
	UpdateProfile(ctx context.Context, token string, profile Profile) error

	// ViewProfile retrieves data about the profile identified by the provided
	// ID, that belongs to the user identified by the provided key.
	ViewProfile(ctx context.Context, token, id string) (Profile, error)

	// ListProfiles retrieves data about subset of profiles that belongs to the
	// user identified by the provided key.
	ListProfiles(ctx context.Context, token string, pm apiutil.PageMetadata) (ProfilesPage, error)

	// ListProfilesByOrg retrieves page of profiles that belong to an org identified by ID.
	ListProfilesByOrg(ctx context.Context, token string, orgID string, pm apiutil.PageMetadata) (ProfilesPage, error)

	// ViewProfileByThing retrieves data about profile that have
	// specified thing connected or not connected to it and belong to the user identified by
	// the provided key.
	ViewProfileByThing(ctx context.Context, token, thID string) (Profile, error)

	// ViewMetadataByKey retrieves metadata about the thing identified by the given key.
	ViewMetadataByKey(ctx context.Context, thingKey string) (Metadata, error)

	// RemoveProfiles removes the things identified by the provided IDs, that
	// belongs to the user identified by the provided key.
	RemoveProfiles(ctx context.Context, token string, ids ...string) error

	// GetPubConfByKey determines whether the profile can be accessed using the
	// provided key and returns thing's id if access is allowed.
	GetPubConfByKey(ctx context.Context, key string) (PubConfInfo, error)

	// GetConfigByThingID returns profile config for given thing ID.
	GetConfigByThingID(ctx context.Context, thingID string) (map[string]interface{}, error)

	// CanUserAccessThing determines whether a user has access to a thing.
	CanUserAccessThing(ctx context.Context, req UserAccessReq) error

	// CanUserAccessProfile determines whether a user has access to a profile.
	CanUserAccessProfile(ctx context.Context, req UserAccessReq) error

	// CanUserAccessGroup determines whether a user has access to a group.
	CanUserAccessGroup(ctx context.Context, req UserAccessReq) error

	// CanThingAccessGroup determines whether a given thing has access to a group with a key.
	CanThingAccessGroup(ctx context.Context, req ThingAccessReq) error

	// Identify returns thing ID for given thing key.
	Identify(ctx context.Context, key string) (string, error)

	// GetGroupIDByThingID returns a thing's group ID for given thing ID.
	GetGroupIDByThingID(ctx context.Context, thingID string) (string, error)

	// GetGroupIDByProfileID returns a profile's group ID for given profile ID.
	GetGroupIDByProfileID(ctx context.Context, profileID string) (string, error)

	// GetProfileIDByThingID returns a thing's profile ID for given thing ID.
	GetProfileIDByThingID(ctx context.Context, thingID string) (string, error)

	// Backup retrieves all things, profiles, groups, and groups memberships for all users. Only accessible by admin.
	Backup(ctx context.Context, token string) (Backup, error)

	// BackupGroupsByOrg retrieves all groups for given org ID.
	BackupGroupsByOrg(ctx context.Context, token string, orgID string) (BackupGroupsByOrg, error)

	// BackupGroupMemberships retrieves all group memberships for given group ID.
	BackupGroupMemberships(ctx context.Context, token string, groupID string) (BackupGroupMemberships, error)

	// BackupProfilesByOrg retrieves all profiles for given org ID.
	BackupProfilesByOrg(ctx context.Context, token string, orgID string) (ProfilesBackup, error)

	// BackupProfilesByGroup retrieves all profiles for given group ID.
	BackupProfilesByGroup(ctx context.Context, token string, groupID string) (ProfilesBackup, error)

	// Restore adds things, profiles, groups, and groups memberships from a backup. Only accessible by admin.
	Restore(ctx context.Context, token string, backup Backup) error

	Groups

	GroupMemberships
}

type Backup struct {
	Things           []Thing
	Profiles         []Profile
	Groups           []Group
	GroupMemberships []GroupMembership
}

type BackupGroupsByOrg struct {
	Groups []Group
}

type BackupGroupMemberships struct {
	BackupGroupMemberships []GroupMembership
}

type ProfilesBackup struct {
	Profiles []Profile
}

type ThingsBackup struct {
	Things []Thing
}

type UserAccessReq struct {
	Token  string
	ID     string
	Action string
}

type ThingAccessReq struct {
	Key string
	ID  string
}

type PubConfInfo struct {
	PublisherID   string
	ProfileConfig map[string]interface{}
}

var _ Service = (*thingsService)(nil)

type thingsService struct {
	auth             protomfx.AuthServiceClient
	users            protomfx.UsersServiceClient
	things           ThingRepository
	profiles         ProfileRepository
	groups           GroupRepository
	groupMemberships GroupMembershipsRepository
	profileCache     ProfileCache
	thingCache       ThingCache
	groupCache       GroupCache
	idProvider       uuid.IDProvider
}

// New instantiates the things service implementation.
func New(auth protomfx.AuthServiceClient, users protomfx.UsersServiceClient, things ThingRepository, profiles ProfileRepository, groups GroupRepository, groupMemberships GroupMembershipsRepository, pcache ProfileCache, tcache ThingCache, gcache GroupCache, idp uuid.IDProvider) Service {
	return &thingsService{
		auth:             auth,
		users:            users,
		things:           things,
		profiles:         profiles,
		groups:           groups,
		groupMemberships: groupMemberships,
		profileCache:     pcache,
		thingCache:       tcache,
		groupCache:       gcache,
		idProvider:       idp,
	}
}

func (ts *thingsService) CreateThings(ctx context.Context, token string, things ...Thing) ([]Thing, error) {
	ths := []Thing{}
	for _, thing := range things {
		ar := UserAccessReq{
			Token:  token,
			ID:     thing.GroupID,
			Action: Editor,
		}
		if err := ts.CanUserAccessGroup(ctx, ar); err != nil {
			return nil, err
		}

		prGrID, err := ts.getGroupIDByProfileID(ctx, thing.ProfileID)
		if err != nil {
			return []Thing{}, err
		}

		if prGrID != thing.GroupID {
			return nil, errors.ErrAuthorization
		}

		th, err := ts.createThing(ctx, &thing)
		if err != nil {
			return []Thing{}, err
		}
		ths = append(ths, th)
	}

	return ths, nil
}

func (ts *thingsService) createThing(ctx context.Context, thing *Thing) (Thing, error) {
	if thing.ID == "" {
		id, err := ts.idProvider.ID()
		if err != nil {
			return Thing{}, err
		}
		thing.ID = id
	}

	if thing.Key == "" {
		key, err := ts.idProvider.ID()

		if err != nil {
			return Thing{}, err
		}
		thing.Key = key
	}

	ths, err := ts.things.Save(ctx, *thing)
	if err != nil {
		return Thing{}, err
	}
	if len(ths) == 0 {
		return Thing{}, errors.ErrCreateEntity
	}

	return ths[0], nil
}

func (ts *thingsService) UpdateThing(ctx context.Context, token string, thing Thing) error {
	ar := UserAccessReq{
		Token:  token,
		ID:     thing.ID,
		Action: Editor,
	}
	if err := ts.CanUserAccessThing(ctx, ar); err != nil {
		return err
	}

	thGrID, err := ts.getGroupIDByThingID(ctx, thing.ID)
	if err != nil {
		return err
	}

	prGrID, err := ts.getGroupIDByProfileID(ctx, thing.ProfileID)
	if err != nil {
		return err
	}

	if prGrID != thGrID {
		return errors.ErrAuthorization
	}

	return ts.things.Update(ctx, thing)
}

func (ts *thingsService) UpdateThingsMetadata(ctx context.Context, token string, things ...Thing) error {
	for _, thing := range things {
		ar := UserAccessReq{
			Token:  token,
			ID:     thing.ID,
			Action: Editor,
		}

		if err := ts.CanUserAccessThing(ctx, ar); err != nil {
			return err
		}

		th, err := ts.things.RetrieveByID(ctx, thing.ID)
		if err != nil {
			return err
		}

		for k, v := range thing.Metadata {
			th.Metadata[k] = v
		}

		if err := ts.things.Update(ctx, th); err != nil {
			return err
		}
	}

	return nil
}

func (ts *thingsService) UpdateKey(ctx context.Context, token, id, key string) error {
	ar := UserAccessReq{
		Token:  token,
		ID:     id,
		Action: Editor,
	}
	if err := ts.CanUserAccessThing(ctx, ar); err != nil {
		return err
	}

	return ts.things.UpdateKey(ctx, id, key)
}

func (ts *thingsService) ViewThing(ctx context.Context, token, id string) (Thing, error) {
	ar := UserAccessReq{
		Token:  token,
		ID:     id,
		Action: Viewer,
	}
	if err := ts.CanUserAccessThing(ctx, ar); err != nil {
		return Thing{}, err
	}

	thing, err := ts.things.RetrieveByID(ctx, id)
	if err != nil {
		return Thing{}, err
	}

	return thing, nil
}

func (ts *thingsService) ViewMetadataByKey(ctx context.Context, thingKey string) (Metadata, error) {
	thingID, err := ts.Identify(ctx, thingKey)
	if err != nil {
		return Metadata{}, err
	}

	thing, err := ts.things.RetrieveByID(ctx, thingID)
	if err != nil {
		return Metadata{}, err
	}

	return thing.Metadata, nil
}

func (ts *thingsService) ListThings(ctx context.Context, token string, pm apiutil.PageMetadata) (ThingsPage, error) {
	if err := ts.isAdmin(ctx, token); err == nil {
		return ts.things.RetrieveAll(ctx, pm)
	}

	res, err := ts.auth.Identify(ctx, &protomfx.Token{Value: token})
	if err != nil {
		return ThingsPage{}, errors.Wrap(errors.ErrAuthentication, err)
	}

	grIDs, err := ts.getGroupIDsByMemberID(ctx, res.GetId())
	if err != nil {
		return ThingsPage{}, err
	}

	return ts.things.RetrieveByGroups(ctx, grIDs, pm)
}

func (ts *thingsService) ListThingsByOrg(ctx context.Context, token string, orgID string, pm apiutil.PageMetadata) (ThingsPage, error) {
	if err := ts.isAdmin(ctx, token); err == nil {
		if grIDs, err := ts.groups.RetrieveIDsByOrg(ctx, orgID); err == nil {
			return ts.things.RetrieveByGroups(ctx, grIDs, pm)
		}
		return ThingsPage{}, err
	}

	if err := ts.canAccessOrg(ctx, token, orgID, auth.OrgSub, Viewer); err != nil {
		return ThingsPage{}, err
	}

	user, err := ts.auth.Identify(ctx, &protomfx.Token{Value: token})
	if err != nil {
		return ThingsPage{}, errors.Wrap(errors.ErrAuthentication, err)
	}

	grIDs, err := ts.groups.RetrieveIDsByOrgMembership(ctx, orgID, user.GetId())
	if err != nil {
		return ThingsPage{}, err
	}

	return ts.things.RetrieveByGroups(ctx, grIDs, pm)
}

func (ts *thingsService) ListThingsByProfile(ctx context.Context, token, prID string, pm apiutil.PageMetadata) (ThingsPage, error) {
	ar := UserAccessReq{
		Token:  token,
		ID:     prID,
		Action: Viewer,
	}
	if err := ts.CanUserAccessProfile(ctx, ar); err != nil {
		return ThingsPage{}, err
	}

	tp, err := ts.things.RetrieveByProfile(ctx, prID, pm)
	if err != nil {
		return ThingsPage{}, err
	}

	return tp, nil
}

func (ts *thingsService) BackupThingsByGroup(ctx context.Context, token string, groupID string) (ThingsBackup, error) {
	if err := ts.canAccessGroup(ctx, token, groupID, Owner); err != nil {
		return ThingsBackup{}, err
	}

	things, err := ts.things.BackupByGroups(ctx, []string{groupID})
	if err != nil {
		return ThingsBackup{}, err
	}

	return ThingsBackup{
		Things: things,
	}, nil
}

func (ts *thingsService) BackupThingsByOrg(ctx context.Context, token string, orgID string) (ThingsBackup, error) {
	if err := ts.canAccessOrg(ctx, token, orgID, auth.OrgSub, Owner); err != nil {
		return ThingsBackup{}, err
	}

	grIDs, err := ts.groups.RetrieveIDsByOrg(ctx, orgID)
	if err != nil {
		return ThingsBackup{}, err
	}

	things, err := ts.things.BackupByGroups(ctx, grIDs)
	if err != nil {
		return ThingsBackup{}, err
	}

	return ThingsBackup{
		Things: things,
	}, nil
}

func (ts *thingsService) RemoveThings(ctx context.Context, token string, ids ...string) error {
	for _, id := range ids {
		ar := UserAccessReq{
			Token:  token,
			ID:     id,
			Action: Editor,
		}
		if err := ts.CanUserAccessThing(ctx, ar); err != nil {
			return err
		}

		if err := ts.thingCache.Remove(ctx, id); err != nil {
			return err
		}

		if err := ts.thingCache.RemoveGroup(ctx, id); err != nil {
			return err
		}
	}

	if err := ts.things.Remove(ctx, ids...); err != nil {
		return err
	}

	return nil
}

func (ts *thingsService) CreateProfiles(ctx context.Context, token string, profiles ...Profile) ([]Profile, error) {
	prs := []Profile{}
	for _, profile := range profiles {
		ar := UserAccessReq{
			Token:  token,
			ID:     profile.GroupID,
			Action: Editor,
		}
		if err := ts.CanUserAccessGroup(ctx, ar); err != nil {
			return nil, err
		}

		pr, err := ts.createProfile(ctx, &profile)
		if err != nil {
			return []Profile{}, err
		}
		prs = append(prs, pr)
	}
	return prs, nil
}

func (ts *thingsService) createProfile(ctx context.Context, profile *Profile) (Profile, error) {
	if profile.ID == "" {
		prID, err := ts.idProvider.ID()
		if err != nil {
			return Profile{}, err
		}
		profile.ID = prID
	}

	prs, err := ts.profiles.Save(ctx, *profile)
	if err != nil {
		return Profile{}, err
	}
	if len(prs) == 0 {
		return Profile{}, errors.ErrCreateEntity
	}

	return prs[0], nil
}

func (ts *thingsService) UpdateProfile(ctx context.Context, token string, profile Profile) error {
	ar := UserAccessReq{
		Token:  token,
		ID:     profile.ID,
		Action: Editor,
	}
	if err := ts.CanUserAccessProfile(ctx, ar); err != nil {
		return err
	}

	return ts.profiles.Update(ctx, profile)
}

func (ts *thingsService) ViewProfile(ctx context.Context, token, id string) (Profile, error) {
	ar := UserAccessReq{
		Token:  token,
		ID:     id,
		Action: Viewer,
	}
	if err := ts.CanUserAccessProfile(ctx, ar); err != nil {
		return Profile{}, err
	}

	profile, err := ts.profiles.RetrieveByID(ctx, id)
	if err != nil {
		return Profile{}, err
	}

	return profile, nil
}

func (ts *thingsService) ListProfiles(ctx context.Context, token string, pm apiutil.PageMetadata) (ProfilesPage, error) {
	if err := ts.isAdmin(ctx, token); err == nil {
		return ts.profiles.RetrieveAll(ctx, pm)
	}

	res, err := ts.auth.Identify(ctx, &protomfx.Token{Value: token})
	if err != nil {
		return ProfilesPage{}, errors.Wrap(errors.ErrAuthentication, err)
	}

	grIDs, err := ts.getGroupIDsByMemberID(ctx, res.GetId())
	if err != nil {
		return ProfilesPage{}, err
	}

	return ts.profiles.RetrieveByGroups(ctx, grIDs, pm)
}

func (ts *thingsService) ListProfilesByOrg(ctx context.Context, token string, orgID string, pm apiutil.PageMetadata) (ProfilesPage, error) {
	if err := ts.isAdmin(ctx, token); err == nil {
		if grIDs, err := ts.groups.RetrieveIDsByOrg(ctx, orgID); err == nil {
			return ts.profiles.RetrieveByGroups(ctx, grIDs, pm)
		}
		return ProfilesPage{}, err
	}

	if err := ts.canAccessOrg(ctx, token, orgID, auth.OrgSub, Viewer); err != nil {
		return ProfilesPage{}, err
	}

	user, err := ts.auth.Identify(ctx, &protomfx.Token{Value: token})
	if err != nil {
		return ProfilesPage{}, errors.Wrap(errors.ErrAuthentication, err)
	}

	grIDs, err := ts.groups.RetrieveIDsByOrgMembership(ctx, orgID, user.GetId())
	if err != nil {
		return ProfilesPage{}, err
	}

	return ts.profiles.RetrieveByGroups(ctx, grIDs, pm)
}

func (ts *thingsService) ViewProfileByThing(ctx context.Context, token, thID string) (Profile, error) {
	ar := UserAccessReq{
		Token:  token,
		ID:     thID,
		Action: Viewer,
	}
	if err := ts.CanUserAccessThing(ctx, ar); err != nil {
		return Profile{}, err
	}

	profile, err := ts.profiles.RetrieveByThing(ctx, thID)
	if err != nil {
		return Profile{}, err
	}

	return profile, nil
}

func (ts *thingsService) RemoveProfiles(ctx context.Context, token string, ids ...string) error {
	for _, id := range ids {
		ar := UserAccessReq{
			Token:  token,
			ID:     id,
			Action: Editor,
		}

		if err := ts.CanUserAccessProfile(ctx, ar); err != nil {
			return err
		}

		if err := ts.profileCache.RemoveGroup(ctx, id); err != nil {
			return err
		}
	}

	return ts.profiles.Remove(ctx, ids...)
}

func (ts *thingsService) GetPubConfByKey(ctx context.Context, thingKey string) (PubConfInfo, error) {
	thID, err := ts.thingCache.ID(ctx, thingKey)
	if err != nil {
		id, err := ts.things.RetrieveByKey(ctx, thingKey)
		if err != nil {
			return PubConfInfo{}, err
		}
		thID = id

		if err := ts.thingCache.Save(ctx, thingKey, thID); err != nil {
			return PubConfInfo{}, err
		}
	}

	profile, err := ts.profiles.RetrieveByThing(ctx, thID)
	if err != nil {
		return PubConfInfo{}, err
	}

	return PubConfInfo{PublisherID: thID, ProfileConfig: profile.Config}, nil
}

func (ts *thingsService) GetConfigByThingID(ctx context.Context, thingID string) (map[string]interface{}, error) {
	profile, err := ts.profiles.RetrieveByThing(ctx, thingID)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return profile.Config, nil
}

func (ts *thingsService) CanUserAccessThing(ctx context.Context, req UserAccessReq) error {
	grID, err := ts.getGroupIDByThingID(ctx, req.ID)
	if err != nil {
		return err
	}

	return ts.canAccessGroup(ctx, req.Token, grID, req.Action)
}

func (ts *thingsService) CanUserAccessProfile(ctx context.Context, req UserAccessReq) error {
	grID, err := ts.getGroupIDByProfileID(ctx, req.ID)
	if err != nil {
		return err
	}

	return ts.canAccessGroup(ctx, req.Token, grID, req.Action)
}

func (ts *thingsService) CanUserAccessGroup(ctx context.Context, req UserAccessReq) error {
	return ts.canAccessGroup(ctx, req.Token, req.ID, req.Action)
}

func (ts *thingsService) CanThingAccessGroup(ctx context.Context, req ThingAccessReq) error {
	thID, err := ts.Identify(ctx, req.Key)
	if err != nil {
		return err
	}

	grID, err := ts.getGroupIDByThingID(ctx, thID)
	if err != nil {
		return err
	}

	if grID != req.ID {
		return errors.ErrAuthorization
	}

	return nil
}

func (ts *thingsService) Identify(ctx context.Context, key string) (string, error) {
	id, err := ts.thingCache.ID(ctx, key)
	if err == nil {
		return id, nil
	}

	id, err = ts.things.RetrieveByKey(ctx, key)
	if err != nil {
		return "", err
	}

	if err := ts.thingCache.Save(ctx, key, id); err != nil {
		return "", err
	}
	return id, nil
}

func (ts *thingsService) GetGroupIDByThingID(ctx context.Context, thingID string) (string, error) {
	return ts.getGroupIDByThingID(ctx, thingID)
}

func (ts *thingsService) GetGroupIDByProfileID(ctx context.Context, profileID string) (string, error) {
	return ts.getGroupIDByProfileID(ctx, profileID)
}

func (ts *thingsService) GetProfileIDByThingID(ctx context.Context, thingID string) (string, error) {
	th, err := ts.things.RetrieveByID(ctx, thingID)
	if err != nil {
		return "", err
	}

	return th.ProfileID, nil
}

func (ts *thingsService) Backup(ctx context.Context, token string) (Backup, error) {
	if err := ts.isAdmin(ctx, token); err != nil {
		return Backup{}, err
	}

	groups, err := ts.groups.BackupAll(ctx)
	if err != nil {
		return Backup{}, err
	}

	groupMemberships, err := ts.groupMemberships.BackupAll(ctx)
	if err != nil {
		return Backup{}, err
	}

	things, err := ts.things.BackupAll(ctx)
	if err != nil {
		return Backup{}, err
	}

	profiles, err := ts.profiles.BackupAll(ctx)
	if err != nil {
		return Backup{}, err
	}

	return Backup{
		Things:           things,
		Profiles:         profiles,
		Groups:           groups,
		GroupMemberships: groupMemberships,
	}, nil
}

func (ts *thingsService) BackupGroupsByOrg(ctx context.Context, token string, orgID string) (BackupGroupsByOrg, error) {
	if err := ts.canAccessOrg(ctx, token, orgID, auth.OrgSub, Owner); err != nil {
		return BackupGroupsByOrg{}, err
	}

	groups, err := ts.groups.BackupByOrg(ctx, orgID)
	if err != nil {
		return BackupGroupsByOrg{}, err
	}

	return BackupGroupsByOrg{
		Groups: groups,
	}, nil
}

func (ts *thingsService) BackupGroupMemberships(ctx context.Context, token string, groupID string) (BackupGroupMemberships, error) {
	ar := UserAccessReq{
		Token:  token,
		ID:     groupID,
		Action: Owner,
	}
	if err := ts.CanUserAccessGroup(ctx, ar); err != nil {
		return BackupGroupMemberships{}, err
	}

	groupMemberships, err := ts.groupMemberships.BackupByGroup(ctx, groupID)
	if err != nil {
		return BackupGroupMemberships{}, err
	}

	var memberIDs []string
	for _, gm := range groupMemberships {
		memberIDs = append(memberIDs, gm.MemberID)
	}

	usersResp, err := ts.users.GetUsersByIDs(ctx, &protomfx.UsersByIDsReq{Ids: memberIDs})
	if err != nil {
		return BackupGroupMemberships{}, err
	}

	emailMap := make(map[string]string)
	for _, user := range usersResp.Users {
		emailMap[user.Id] = user.Email
	}

	for i := range groupMemberships {
		groupMemberships[i].Email = emailMap[groupMemberships[i].MemberID]
	}

	return BackupGroupMemberships{
		BackupGroupMemberships: groupMemberships,
	}, nil
}

func (ts *thingsService) BackupProfilesByOrg(ctx context.Context, token string, orgID string) (ProfilesBackup, error) {
	if err := ts.canAccessOrg(ctx, token, orgID, auth.OrgSub, Owner); err != nil {
		return ProfilesBackup{}, err
	}

	grIDs, err := ts.groups.RetrieveIDsByOrg(ctx, orgID)
	if err != nil {
		return ProfilesBackup{}, err
	}

	profiles, err := ts.profiles.BackupByGroups(ctx, grIDs)
	if err != nil {
		return ProfilesBackup{}, err
	}

	return ProfilesBackup{
		Profiles: profiles,
	}, nil
}

func (ts *thingsService) BackupProfilesByGroup(ctx context.Context, token string, groupID string) (ProfilesBackup, error) {
	if err := ts.canAccessGroup(ctx, token, groupID, Owner); err != nil {
		return ProfilesBackup{}, err
	}

	profiles, err := ts.profiles.BackupByGroups(ctx, []string{groupID})
	if err != nil {
		return ProfilesBackup{}, err
	}

	return ProfilesBackup{
		Profiles: profiles,
	}, nil
}

func (ts *thingsService) Restore(ctx context.Context, token string, backup Backup) error {
	if err := ts.isAdmin(ctx, token); err != nil {
		return err
	}

	for _, group := range backup.Groups {
		if _, err := ts.groups.Save(ctx, group); err != nil {
			return err
		}
	}

	if _, err := ts.things.Save(ctx, backup.Things...); err != nil {
		return err
	}

	if _, err := ts.profiles.Save(ctx, backup.Profiles...); err != nil {
		return err
	}

	for _, g := range backup.GroupMemberships {
		gm := GroupMembership{
			MemberID: g.MemberID,
			GroupID:  g.GroupID,
			Role:     g.Role,
		}

		if err := ts.groupMemberships.Save(ctx, gm); err != nil {
			return err
		}
	}

	return nil
}

func getTimestmap() time.Time {
	return time.Now().UTC().Round(time.Millisecond)
}

func (ts *thingsService) ListThingsByGroup(ctx context.Context, token string, groupID string, pm apiutil.PageMetadata) (ThingsPage, error) {
	ar := UserAccessReq{
		Token:  token,
		ID:     groupID,
		Action: Viewer,
	}
	if err := ts.CanUserAccessGroup(ctx, ar); err != nil {
		return ThingsPage{}, err
	}

	return ts.things.RetrieveByGroups(ctx, []string{groupID}, pm)
}

func (ts *thingsService) ListProfilesByGroup(ctx context.Context, token, groupID string, pm apiutil.PageMetadata) (ProfilesPage, error) {
	ar := UserAccessReq{
		Token:  token,
		ID:     groupID,
		Action: Viewer,
	}
	if err := ts.CanUserAccessGroup(ctx, ar); err != nil {
		return ProfilesPage{}, err
	}

	return ts.profiles.RetrieveByGroups(ctx, []string{groupID}, pm)
}

func (ts *thingsService) isAdmin(ctx context.Context, token string) error {
	req := &protomfx.AuthorizeReq{
		Token:   token,
		Subject: auth.RootSub,
	}

	if _, err := ts.auth.Authorize(ctx, req); err != nil {
		return errors.Wrap(errors.ErrAuthorization, err)
	}

	return nil
}

func (ts *thingsService) canAccessOrg(ctx context.Context, token, orgID, subject, action string) error {
	req := &protomfx.AuthorizeReq{
		Token:   token,
		Object:  orgID,
		Subject: subject,
		Action:  action,
	}

	if _, err := ts.auth.Authorize(ctx, req); err != nil {
		return errors.Wrap(errors.ErrAuthorization, err)
	}

	return nil
}

func (ts *thingsService) getGroupIDByThingID(ctx context.Context, thID string) (string, error) {
	grID, err := ts.thingCache.ViewGroup(ctx, thID)
	if err != nil {
		th, err := ts.things.RetrieveByID(ctx, thID)
		if err != nil {
			return "", err
		}
		grID = th.GroupID

		if err := ts.thingCache.SaveGroup(ctx, th.ID, th.GroupID); err != nil {
			return "", err
		}
	}

	return grID, nil
}

func (ts *thingsService) getGroupIDByProfileID(ctx context.Context, prID string) (string, error) {
	grID, err := ts.profileCache.ViewGroup(ctx, prID)
	if err != nil {
		pr, err := ts.profiles.RetrieveByID(ctx, prID)
		if err != nil {
			return "", err
		}
		grID = pr.GroupID

		if err := ts.profileCache.SaveGroup(ctx, pr.ID, pr.GroupID); err != nil {
			return "", err
		}
	}

	return grID, nil
}

func (ts *thingsService) getGroupIDsByMemberID(ctx context.Context, memberID string) ([]string, error) {
	grIDs, err := ts.groupCache.RetrieveGroupIDsByMember(ctx, memberID)
	if err != nil {
		grIDs, err = ts.groupMemberships.RetrieveGroupIDsByMember(ctx, memberID)
		if err != nil {
			return []string{}, err
		}
	}
	return grIDs, nil
}
