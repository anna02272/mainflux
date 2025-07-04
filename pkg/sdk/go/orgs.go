// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/MainfluxLabs/mainflux/pkg/errors"
)

const orgsEndpoint = "orgs"

func (sdk mfSDK) CreateOrg(o Org, token string) (string, error) {
	data, err := json.Marshal(o)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s", sdk.authURL, orgsEndpoint)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", errors.Wrap(ErrFailedCreation, errors.New(resp.Status))
	}

	id := strings.TrimPrefix(resp.Header.Get("Location"), fmt.Sprintf("/%s/", orgsEndpoint))
	return id, nil
}

func (sdk mfSDK) GetOrg(id, token string) (Org, error) {
	url := fmt.Sprintf("%s/%s/%s", sdk.authURL, orgsEndpoint, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Org{}, err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return Org{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Org{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Org{}, errors.Wrap(ErrFailedFetch, errors.New(resp.Status))
	}

	var o Org
	if err := json.Unmarshal(body, &o); err != nil {
		return Org{}, err
	}

	return o, nil
}

func (sdk mfSDK) UpdateOrg(o Org, orgID, token string) error {
	data, err := json.Marshal(o)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s/%s", sdk.authURL, orgsEndpoint, orgID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(ErrFailedUpdate, errors.New(resp.Status))
	}

	return nil
}

func (sdk mfSDK) DeleteOrg(id, token string) error {
	url := fmt.Sprintf("%s/%s/%s", sdk.authURL, orgsEndpoint, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrap(ErrFailedRemoval, errors.New(resp.Status))
	}

	return nil
}

func (sdk mfSDK) ListOrgs(meta PageMetadata, token string) (OrgsPage, error) {
	u, err := url.Parse(sdk.authURL)
	if err != nil {
		return OrgsPage{}, err
	}
	u.Path = orgsEndpoint
	q := u.Query()
	q.Add("offset", strconv.FormatUint(meta.Offset, 10))
	if meta.Limit != 0 {
		q.Add("limit", strconv.FormatUint(meta.Limit, 10))
	}
	if meta.Name != "" {
		q.Add("name", meta.Name)
	}

	u.RawQuery = q.Encode()
	return sdk.getOrgs(token, u.String())
}

func (sdk mfSDK) getOrgs(token, url string) (OrgsPage, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return OrgsPage{}, err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return OrgsPage{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OrgsPage{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return OrgsPage{}, errors.Wrap(ErrFailedFetch, errors.New(resp.Status))
	}

	var op OrgsPage
	if err := json.Unmarshal(body, &op); err != nil {
		return OrgsPage{}, err
	}
	return op, nil
}

func (sdk mfSDK) GetMember(memberID, orgID, token string) (Member, error) {
	url := fmt.Sprintf("%s/%s/%s/%s/%s", sdk.authURL, orgsEndpoint, orgID, membersEndpoint, memberID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Member{}, err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return Member{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Member{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Member{}, errors.Wrap(ErrFailedFetch, errors.New(resp.Status))
	}

	var member Member
	if err := json.Unmarshal(body, &member); err != nil {
		return Member{}, err
	}

	return member, nil
}

func (sdk mfSDK) AssignMembers(om []OrgMember, orgID string, token string) error {
	url := fmt.Sprintf("%s/%s/%s/%s", sdk.authURL, orgsEndpoint, orgID, membersEndpoint)

	assignMemberReq := assignMembersReq{
		OrgMembers: om,
	}

	data, err := json.Marshal(assignMemberReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(ErrMemberAdd, errors.New(resp.Status))
	}

	return nil
}

func (sdk mfSDK) UnassignMembers(memberIDs []string, orgID, token string) error {
	var ids []string
	url := fmt.Sprintf("%s/%s/%s/%s", sdk.authURL, orgsEndpoint, orgID, membersEndpoint)
	ids = append(ids, memberIDs...)
	unassignMembersReq := unassignMemberReq{
		MemberIDs: ids,
	}

	data, err := json.Marshal(unassignMembersReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrap(ErrFailedRemoval, errors.New(resp.Status))
	}

	return nil
}

func (sdk mfSDK) UpdateMembers(oms []OrgMember, orgID, token string) error {
	url := fmt.Sprintf("%s/%s/%s/%s", sdk.authURL, orgsEndpoint, orgID, membersEndpoint)
	updateMembersReq := updateMemberReq{
		OrgMembers: oms,
	}

	data, err := json.Marshal(updateMembersReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(ErrFailedUpdate, errors.New(resp.Status))
	}

	return nil
}

func (sdk mfSDK) ListMembersByOrg(orgID string, pm PageMetadata, token string) (MembersPage, error) {
	url := fmt.Sprintf("%s/%s/%s/members?offset=%d&limit=%d", sdk.authURL, orgsEndpoint, orgID, pm.Offset, pm.Limit)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return MembersPage{}, err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return MembersPage{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return MembersPage{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return MembersPage{}, errors.Wrap(ErrFailedFetch, errors.New(resp.Status))
	}

	var mp MembersPage
	if err := json.Unmarshal(body, &mp); err != nil {
		return MembersPage{}, err
	}

	return mp, nil
}

func (sdk mfSDK) ListGroupsByOrg(orgID string, meta PageMetadata, token string) (GroupsPage, error) {
	apiUrl := fmt.Sprintf("%s/%s/%s/groups?offset=%d&limit=%d", sdk.thingsURL, orgsEndpoint, orgID, meta.Offset, meta.Limit)

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return GroupsPage{}, err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return GroupsPage{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GroupsPage{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return GroupsPage{}, errors.Wrap(ErrFailedFetch, errors.New(resp.Status))
	}

	var gp GroupsPage
	if err := json.Unmarshal(body, &gp); err != nil {
		return GroupsPage{}, err
	}

	return gp, nil
}
