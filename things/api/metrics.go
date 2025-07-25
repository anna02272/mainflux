// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"context"
	"time"

	"github.com/MainfluxLabs/mainflux/pkg/apiutil"
	"github.com/MainfluxLabs/mainflux/things"
	"github.com/go-kit/kit/metrics"
)

var _ things.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     things.Service
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc things.Service, counter metrics.Counter, latency metrics.Histogram) things.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) CreateThings(ctx context.Context, token string, ths ...things.Thing) (saved []things.Thing, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_things").Add(1)
		ms.latency.With("method", "create_things").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreateThings(ctx, token, ths...)
}

func (ms *metricsMiddleware) UpdateThing(ctx context.Context, token string, thing things.Thing) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_thing").Add(1)
		ms.latency.With("method", "update_thing").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateThing(ctx, token, thing)
}

func (ms *metricsMiddleware) UpdateThingsMetadata(ctx context.Context, token string, things ...things.Thing) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_things_metadata").Add(1)
		ms.latency.With("method", "update_things_metadata").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateThingsMetadata(ctx, token, things...)
}

func (ms *metricsMiddleware) UpdateKey(ctx context.Context, token, id, key string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_key").Add(1)
		ms.latency.With("method", "update_key").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateKey(ctx, token, id, key)
}

func (ms *metricsMiddleware) ViewThing(ctx context.Context, token, id string) (things.Thing, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_thing").Add(1)
		ms.latency.With("method", "view_thing").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewThing(ctx, token, id)
}

func (ms *metricsMiddleware) ViewMetadataByKey(ctx context.Context, thingKey string) (things.Metadata, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_metadata_by_key").Add(1)
		ms.latency.With("method", "view_metadata_by_key").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewMetadataByKey(ctx, thingKey)
}

func (ms *metricsMiddleware) ListThings(ctx context.Context, token string, pm apiutil.PageMetadata) (things.ThingsPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_things").Add(1)
		ms.latency.With("method", "list_things").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListThings(ctx, token, pm)
}

func (ms *metricsMiddleware) ListThingsByProfile(ctx context.Context, token, prID string, pm apiutil.PageMetadata) (things.ThingsPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_things_by_profile").Add(1)
		ms.latency.With("method", "list_things_by_profile").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListThingsByProfile(ctx, token, prID, pm)
}

func (ms *metricsMiddleware) ListThingsByOrg(ctx context.Context, token string, orgID string, pm apiutil.PageMetadata) (things.ThingsPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_things_by_org").Add(1)
		ms.latency.With("method", "list_things_by_org").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListThingsByOrg(ctx, token, orgID, pm)
}

func (ms *metricsMiddleware) RemoveThings(ctx context.Context, token string, id ...string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove_things").Add(1)
		ms.latency.With("method", "remove_things").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RemoveThings(ctx, token, id...)
}

func (ms *metricsMiddleware) CreateProfiles(ctx context.Context, token string, profiles ...things.Profile) (saved []things.Profile, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_profiles").Add(1)
		ms.latency.With("method", "create_profiles").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreateProfiles(ctx, token, profiles...)
}

func (ms *metricsMiddleware) UpdateProfile(ctx context.Context, token string, profile things.Profile) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_profile").Add(1)
		ms.latency.With("method", "update_profile").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateProfile(ctx, token, profile)
}

func (ms *metricsMiddleware) ViewProfile(ctx context.Context, token, id string) (things.Profile, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_profile").Add(1)
		ms.latency.With("method", "view_profile").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewProfile(ctx, token, id)
}

func (ms *metricsMiddleware) ListProfiles(ctx context.Context, token string, pm apiutil.PageMetadata) (things.ProfilesPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_profiles").Add(1)
		ms.latency.With("method", "list_profiles").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListProfiles(ctx, token, pm)
}

func (ms *metricsMiddleware) ListProfilesByOrg(ctx context.Context, token string, orgID string, pm apiutil.PageMetadata) (things.ProfilesPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_profiles_by_org").Add(1)
		ms.latency.With("method", "list_profiles_by_org").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListProfilesByOrg(ctx, token, orgID, pm)
}

func (ms *metricsMiddleware) ViewProfileByThing(ctx context.Context, token, thID string) (things.Profile, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_profile_by_thing").Add(1)
		ms.latency.With("method", "view_profile_by_thing").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewProfileByThing(ctx, token, thID)
}

func (ms *metricsMiddleware) RemoveProfiles(ctx context.Context, token string, ids ...string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove_profiles").Add(1)
		ms.latency.With("method", "remove_profiles").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RemoveProfiles(ctx, token, ids...)
}

func (ms *metricsMiddleware) GetPubConfByKey(ctx context.Context, key string) (things.PubConfInfo, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "get_pub_conf_by_key").Add(1)
		ms.latency.With("method", "get_pub_conf_by_key").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.GetPubConfByKey(ctx, key)
}

func (ms *metricsMiddleware) GetConfigByThingID(ctx context.Context, thingID string) (map[string]interface{}, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "get_config_by_thing_id").Add(1)
		ms.latency.With("method", "get_config_by_thing_id").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.GetConfigByThingID(ctx, thingID)
}

func (ms *metricsMiddleware) CanUserAccessThing(ctx context.Context, req things.UserAccessReq) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "can_user_access_thing").Add(1)
		ms.latency.With("method", "can_user_access_thing").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CanUserAccessThing(ctx, req)
}

func (ms *metricsMiddleware) CanUserAccessProfile(ctx context.Context, req things.UserAccessReq) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "can_user_access_profile").Add(1)
		ms.latency.With("method", "can_user_access_profile").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CanUserAccessProfile(ctx, req)
}

func (ms *metricsMiddleware) CanUserAccessGroup(ctx context.Context, req things.UserAccessReq) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "can_user_access_group").Add(1)
		ms.latency.With("method", "can_user_access_group").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CanUserAccessGroup(ctx, req)
}

func (ms *metricsMiddleware) CanThingAccessGroup(ctx context.Context, req things.ThingAccessReq) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "can_thing_access_group").Add(1)
		ms.latency.With("method", "can_thing_access_group").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CanThingAccessGroup(ctx, req)
}

func (ms *metricsMiddleware) Identify(ctx context.Context, key string) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "identify").Add(1)
		ms.latency.With("method", "identify").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Identify(ctx, key)
}

func (ms *metricsMiddleware) GetGroupIDByThingID(ctx context.Context, thingID string) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "get_group_id_by_thing_id").Add(1)
		ms.latency.With("method", "get_group_id_by_thing_id").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.GetGroupIDByThingID(ctx, thingID)
}

func (ms *metricsMiddleware) GetGroupIDByProfileID(ctx context.Context, profileID string) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "get_group_id_by_profile_id").Add(1)
		ms.latency.With("method", "get_group_id_by_profile_id").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.GetGroupIDByProfileID(ctx, profileID)
}

func (ms *metricsMiddleware) GetProfileIDByThingID(ctx context.Context, thingID string) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "get_profile_id_by_thing_id").Add(1)
		ms.latency.With("method", "get_profile_id_by_thing_id").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.GetProfileIDByThingID(ctx, thingID)
}

func (ms *metricsMiddleware) Backup(ctx context.Context, token string) (bk things.Backup, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "backup").Add(1)
		ms.latency.With("method", "backup").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Backup(ctx, token)
}

func (ms *metricsMiddleware) BackupGroupsByOrg(ctx context.Context, token string, orgID string) (bk things.BackupGroupsByOrg, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "backup_groups_by_org").Add(1)
		ms.latency.With("method", "backup_groups_by_org").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.BackupGroupsByOrg(ctx, token, orgID)
}

func (ms *metricsMiddleware) BackupGroupMemberships(ctx context.Context, token string, groupID string) (bk things.BackupGroupMemberships, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "backup_group_memberships").Add(1)
		ms.latency.With("method", "backup_group_memberships").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.BackupGroupMemberships(ctx, token, groupID)
}

func (ms *metricsMiddleware) BackupProfilesByOrg(ctx context.Context, token string, orgID string) (pb things.ProfilesBackup, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "backup_profiles_by_org").Add(1)
		ms.latency.With("method", "backup_profiles_by_org").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.BackupProfilesByOrg(ctx, token, orgID)
}

func (ms *metricsMiddleware) BackupProfilesByGroup(ctx context.Context, token string, groupID string) (pb things.ProfilesBackup, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "backup_profiles_by_group").Add(1)
		ms.latency.With("method", "backup_profiles_by_group").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.BackupProfilesByGroup(ctx, token, groupID)
}

func (ms *metricsMiddleware) BackupThingsByGroup(ctx context.Context, token string, groupID string) (tb things.ThingsBackup, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "backup_things_by_group").Add(1)
		ms.latency.With("method", "backup_things_by_group").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.BackupThingsByGroup(ctx, token, groupID)
}

func (ms *metricsMiddleware) BackupThingsByOrg(ctx context.Context, token string, orgID string) (tb things.ThingsBackup, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "backup_things_by_org").Add(1)
		ms.latency.With("method", "backup_things_by_org").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.BackupThingsByOrg(ctx, token, orgID)
}

func (ms *metricsMiddleware) Restore(ctx context.Context, token string, backup things.Backup) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "restore").Add(1)
		ms.latency.With("method", "restore").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Restore(ctx, token, backup)
}

func (ms *metricsMiddleware) CreateGroups(ctx context.Context, token string, grs ...things.Group) ([]things.Group, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_groups").Add(1)
		ms.latency.With("method", "create_groups").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreateGroups(ctx, token, grs...)
}

func (ms *metricsMiddleware) UpdateGroup(ctx context.Context, token string, g things.Group) (things.Group, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_group").Add(1)
		ms.latency.With("method", "update_group").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateGroup(ctx, token, g)
}

func (ms *metricsMiddleware) ViewGroup(ctx context.Context, token, id string) (things.Group, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_group").Add(1)
		ms.latency.With("method", "view_group").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewGroup(ctx, token, id)
}

func (ms *metricsMiddleware) ListGroups(ctx context.Context, token string, pm apiutil.PageMetadata) (things.GroupPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_groups").Add(1)
		ms.latency.With("method", "list_groups").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListGroups(ctx, token, pm)
}

func (ms *metricsMiddleware) ListGroupsByOrg(ctx context.Context, token, orgID string, pm apiutil.PageMetadata) (things.GroupPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_groups_by_org").Add(1)
		ms.latency.With("method", "list_groups_by_org").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListGroupsByOrg(ctx, token, orgID, pm)
}

func (ms *metricsMiddleware) ListThingsByGroup(ctx context.Context, token, groupID string, pm apiutil.PageMetadata) (tp things.ThingsPage, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_things_by_group").Add(1)
		ms.latency.With("method", "list_things_by_group").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListThingsByGroup(ctx, token, groupID, pm)
}

func (ms *metricsMiddleware) ViewGroupByThing(ctx context.Context, token, thingID string) (things.Group, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_group_by_thing").Add(1)
		ms.latency.With("method", "view_group_by_thing").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewGroupByThing(ctx, token, thingID)
}

func (ms *metricsMiddleware) RemoveGroups(ctx context.Context, token string, ids ...string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove_groups").Add(1)
		ms.latency.With("method", "remove_groups").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RemoveGroups(ctx, token, ids...)
}

func (ms *metricsMiddleware) ListProfilesByGroup(ctx context.Context, token, groupID string, pm apiutil.PageMetadata) (things.ProfilesPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_profiles_by_group").Add(1)
		ms.latency.With("method", "list_profiles_by_group").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListProfilesByGroup(ctx, token, groupID, pm)
}

func (ms *metricsMiddleware) ViewGroupByProfile(ctx context.Context, token, profileID string) (things.Group, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_group_by_profile").Add(1)
		ms.latency.With("method", "view_group_by_profile").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewGroupByProfile(ctx, token, profileID)
}

func (ms *metricsMiddleware) CreateGroupMemberships(ctx context.Context, token string, gms ...things.GroupMembership) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_group_memberships").Add(1)
		ms.latency.With("method", "create_group_memberships").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreateGroupMemberships(ctx, token, gms...)
}

func (ms *metricsMiddleware) ListGroupMemberships(ctx context.Context, token, groupID string, pm apiutil.PageMetadata) (things.GroupMembershipsPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_group_memberships").Add(1)
		ms.latency.With("method", "list_group_memberships").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListGroupMemberships(ctx, token, groupID, pm)
}

func (ms *metricsMiddleware) UpdateGroupMemberships(ctx context.Context, token string, gms ...things.GroupMembership) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_group_memberships").Add(1)
		ms.latency.With("method", "update_group_memberships").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateGroupMemberships(ctx, token, gms...)
}

func (ms *metricsMiddleware) RemoveGroupMemberships(ctx context.Context, token, groupID string, memberIDs ...string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove_group_memberships").Add(1)
		ms.latency.With("method", "remove_group_memberships").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RemoveGroupMemberships(ctx, token, groupID, memberIDs...)
}
