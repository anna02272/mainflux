// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package things

import (
	"context"

	"github.com/MainfluxLabs/mainflux/things"
	"github.com/MainfluxLabs/mainflux/things/api/http/memberships"
	"github.com/go-kit/kit/endpoint"
)

func createThingsEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createThingsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		groupID, err := svc.GetGroupIDByProfileID(ctx, req.profileID)
		if err != nil {
			return nil, err
		}

		ths := []things.Thing{}
		for _, t := range req.Things {
			th := things.Thing{
				ID:        t.ID,
				GroupID:   groupID,
				ProfileID: req.profileID,
				Name:      t.Name,
				Key:       t.Key,
				Metadata:  t.Metadata,
			}
			ths = append(ths, th)
		}

		saved, err := svc.CreateThings(ctx, req.token, ths...)
		if err != nil {
			return nil, err
		}

		res := thingsRes{
			Things:  []thingRes{},
			created: true,
		}

		for _, t := range saved {
			th := thingRes{
				ID:        t.ID,
				GroupID:   t.GroupID,
				ProfileID: t.ProfileID,
				Name:      t.Name,
				Key:       t.Key,
				Metadata:  t.Metadata,
			}
			res.Things = append(res.Things, th)
		}

		return res, nil
	}
}

func viewThingEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(resourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		thing, err := svc.ViewThing(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := viewThingRes{
			ID:        thing.ID,
			GroupID:   thing.GroupID,
			ProfileID: thing.ProfileID,
			Name:      thing.Name,
			Key:       thing.Key,
			Metadata:  thing.Metadata,
		}
		return res, nil
	}
}

func viewMetadataByKeyEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewMetadataReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		metadata, err := svc.ViewMetadataByKey(ctx, req.key)
		if err != nil {
			return nil, err
		}

		res := viewMetadataRes{
			Metadata: metadata,
		}

		return res, nil
	}
}

func listThingsEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		page, err := svc.ListThings(ctx, req.token, req.pageMetadata)
		if err != nil {
			return nil, err
		}

		return buildThingsResponse(page), nil
	}
}

func listThingsByProfileEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listByProfileReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		page, err := svc.ListThingsByProfile(ctx, req.token, req.id, req.pageMetadata)
		if err != nil {
			return nil, err
		}

		return buildThingsResponse(page), nil
	}
}

func listThingsByGroupEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listByGroupReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		page, err := svc.ListThingsByGroup(ctx, req.token, req.id, req.pageMetadata)
		if err != nil {
			return nil, err
		}

		return buildThingsResponse(page), nil
	}
}

func listThingsByOrgEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listByOrgReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		page, err := svc.ListThingsByOrg(ctx, req.token, req.id, req.pageMetadata)
		if err != nil {
			return nil, err
		}

		return buildThingsResponse(page), nil
	}
}

func backupThingsByGroupEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(backupByGroupReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		backup, err := svc.BackupThingsByGroup(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}
		return buildBackupThingsResponse(backup), nil
	}
}

func backupThingsByOrgEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(backupByOrgReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		backup, err := svc.BackupThingsByOrg(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}
		return buildBackupThingsResponse(backup), nil
	}
}

func updateThingEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateThingReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		thing := things.Thing{
			ID:        req.id,
			ProfileID: req.ProfileID,
			Name:      req.Name,
			Metadata:  req.Metadata,
		}

		if err := svc.UpdateThing(ctx, req.token, thing); err != nil {
			return nil, err
		}

		res := thingRes{ID: req.id, created: false}
		return res, nil
	}
}

func updateThingsMetadataEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateThingsMetadataReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		ths := []things.Thing{}
		for _, t := range req.Things {
			th := things.Thing{
				ID:       t.ID,
				Metadata: t.Metadata,
			}
			ths = append(ths, th)
		}

		if err := svc.UpdateThingsMetadata(ctx, req.token, ths...); err != nil {
			return nil, err
		}

		res := thingsRes{
			created: false,
		}
		return res, nil
	}
}

func updateKeyEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateKeyReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.UpdateKey(ctx, req.token, req.id, req.Key); err != nil {
			return nil, err
		}

		res := thingRes{ID: req.id, created: false}
		return res, nil
	}
}

func removeThingEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(resourceReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.RemoveThings(ctx, req.token, req.id); err != nil {
			return nil, err
		}

		return removeRes{}, nil
	}
}

func removeThingsEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(removeThingsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.RemoveThings(ctx, req.token, req.ThingIDs...); err != nil {
			return nil, err
		}

		return removeRes{}, nil
	}
}

func identifyEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(identifyReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		id, err := svc.Identify(ctx, req.Token)
		if err != nil {
			return nil, err
		}

		res := identityRes{
			ID: id,
		}

		return res, nil
	}
}

func backupEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(backupReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		backup, err := svc.Backup(ctx, req.token)
		if err != nil {
			return nil, err
		}

		return buildBackupResponse(backup), nil
	}
}

func restoreEndpoint(svc things.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(restoreReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		backup := buildBackup(req)

		if err := svc.Restore(ctx, req.token, backup); err != nil {
			return nil, err
		}

		return restoreRes{}, nil
	}
}

func buildThingsResponse(tp things.ThingsPage) ThingsPageRes {
	res := ThingsPageRes{
		pageRes: pageRes{
			Total:  tp.Total,
			Offset: tp.Offset,
			Limit:  tp.Limit,
			Order:  tp.Order,
			Dir:    tp.Dir,
			Name:   tp.Name,
		},
		Things: []thingRes{},
	}

	for _, t := range tp.Things {
		view := thingRes{
			ID:        t.ID,
			GroupID:   t.GroupID,
			ProfileID: t.ProfileID,
			Name:      t.Name,
			Key:       t.Key,
			Metadata:  t.Metadata,
		}
		res.Things = append(res.Things, view)
	}

	return res
}

func buildBackupThingsResponse(tb things.ThingsBackup) backupThingsRes {
	res := backupThingsRes{
		Things: []viewThingRes{},
	}
	for _, thing := range tb.Things {
		view := viewThingRes{
			ID:        thing.ID,
			GroupID:   thing.GroupID,
			ProfileID: thing.ProfileID,
			Name:      thing.Name,
			Key:       thing.Key,
			Metadata:  thing.Metadata,
		}
		res.Things = append(res.Things, view)
	}
	return res
}

func buildBackupResponse(backup things.Backup) backupRes {
	res := backupRes{
		Things:           []viewThingRes{},
		Profiles:         []backupProfile{},
		Groups:           []backupGroup{},
		GroupMemberships: []memberships.ViewGroupMembershipRes{},
	}

	for _, thing := range backup.Things {
		view := viewThingRes{
			ID:        thing.ID,
			GroupID:   thing.GroupID,
			ProfileID: thing.ProfileID,
			Name:      thing.Name,
			Key:       thing.Key,
			Metadata:  thing.Metadata,
		}
		res.Things = append(res.Things, view)
	}

	for _, profile := range backup.Profiles {
		view := backupProfile{
			ID:       profile.ID,
			GroupID:  profile.GroupID,
			Name:     profile.Name,
			Config:   profile.Config,
			Metadata: profile.Metadata,
		}
		res.Profiles = append(res.Profiles, view)
	}

	for _, group := range backup.Groups {
		view := backupGroup{
			ID:          group.ID,
			Name:        group.Name,
			OrgID:       group.OrgID,
			Description: group.Description,
			Metadata:    group.Metadata,
			CreatedAt:   group.CreatedAt,
			UpdatedAt:   group.UpdatedAt,
		}
		res.Groups = append(res.Groups, view)
	}

	for _, membership := range backup.GroupMemberships {
		view := memberships.ViewGroupMembershipRes{
			MemberID: membership.MemberID,
			GroupID:  membership.GroupID,
			Email:    membership.Email,
			Role:     membership.Role,
		}
		res.GroupMemberships = append(res.GroupMemberships, view)
	}

	return res
}

func buildBackup(req restoreReq) (backup things.Backup) {
	for _, thing := range req.Things {
		th := things.Thing{
			ID:        thing.ID,
			GroupID:   thing.GroupID,
			ProfileID: thing.ProfileID,
			Name:      thing.Name,
			Key:       thing.Key,
			Metadata:  thing.Metadata,
		}
		backup.Things = append(backup.Things, th)
	}

	for _, profile := range req.Profiles {
		pr := things.Profile{
			ID:       profile.ID,
			GroupID:  profile.GroupID,
			Name:     profile.Name,
			Config:   profile.Config,
			Metadata: profile.Metadata,
		}
		backup.Profiles = append(backup.Profiles, pr)
	}

	for _, group := range req.Groups {
		gr := things.Group{
			ID:          group.ID,
			Name:        group.Name,
			OrgID:       group.OrgID,
			Description: group.Description,
			Metadata:    group.Metadata,
			CreatedAt:   group.CreatedAt,
			UpdatedAt:   group.UpdatedAt,
		}
		backup.Groups = append(backup.Groups, gr)
	}

	for _, membership := range req.GroupMemberships {
		gm := things.GroupMembership{
			GroupID:  membership.GroupID,
			MemberID: membership.MemberID,
			Email:    membership.Email,
			Role:     membership.Role,
		}
		backup.GroupMemberships = append(backup.GroupMemberships, gm)
	}

	return backup
}
