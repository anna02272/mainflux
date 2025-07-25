// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MainfluxLabs/mainflux/pkg/errors"
)

const (
	thingsEndpoint   = "things"
	identifyEndpoint = "identify"
)

type identifyThingReq struct {
	Token string `json:"token,omitempty"`
}

type identifyThingResp struct {
	ID string `json:"id,omitempty"`
}

func (sdk mfSDK) CreateThing(t Thing, profileID, token string) (string, error) {
	things, err := sdk.CreateThings([]Thing{t}, profileID, token)
	if err != nil {
		return "", err
	}

	if len(things) < 1 {
		return "", nil
	}

	return things[0].ID, nil
}

func (sdk mfSDK) CreateThings(things []Thing, profileID, token string) ([]Thing, error) {
	data, err := json.Marshal(things)
	if err != nil {
		return []Thing{}, err
	}

	url := fmt.Sprintf("%s/%s/%s/things", sdk.thingsURL, profilesEndpoint, profileID)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return []Thing{}, err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return []Thing{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return []Thing{}, errors.Wrap(ErrFailedCreation, errors.New(resp.Status))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Thing{}, err
	}

	var ctr createThingsRes
	if err := json.Unmarshal(body, &ctr); err != nil {
		return []Thing{}, err
	}

	return ctr.Things, nil
}

func (sdk mfSDK) ListThings(pm PageMetadata, token string) (ThingsPage, error) {
	url, err := sdk.withQueryParams(sdk.thingsURL, thingsEndpoint, pm)
	if err != nil {
		return ThingsPage{}, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ThingsPage{}, err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return ThingsPage{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ThingsPage{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return ThingsPage{}, errors.Wrap(ErrFailedFetch, errors.New(resp.Status))
	}

	var tp ThingsPage
	if err := json.Unmarshal(body, &tp); err != nil {
		return ThingsPage{}, err
	}

	return tp, nil
}

func (sdk mfSDK) ListThingsByProfile(prID string, pm PageMetadata, token string) (ThingsPage, error) {
	url := fmt.Sprintf("%s/profiles/%s/things?offset=%d&limit=%d&dir=%s", sdk.thingsURL, prID, pm.Offset, pm.Limit, pm.Dir)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ThingsPage{}, err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return ThingsPage{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ThingsPage{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return ThingsPage{}, errors.Wrap(ErrFailedFetch, errors.New(resp.Status))
	}

	var tp ThingsPage
	if err := json.Unmarshal(body, &tp); err != nil {
		return ThingsPage{}, err
	}

	return tp, nil
}

func (sdk mfSDK) GetThing(id, token string) (Thing, error) {
	url := fmt.Sprintf("%s/%s/%s", sdk.thingsURL, thingsEndpoint, id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Thing{}, err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return Thing{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Thing{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Thing{}, errors.Wrap(ErrFailedFetch, errors.New(resp.Status))
	}

	var t Thing
	if err := json.Unmarshal(body, &t); err != nil {
		return Thing{}, err
	}

	return t, nil
}

func (sdk mfSDK) GetThingMetadataByKey(thingKey string) (Metadata, error) {
	url := fmt.Sprintf("%s/metadata", sdk.thingsURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Metadata{}, err
	}

	resp, err := sdk.sendThingRequest(req, thingKey, string(CTJSON))
	if err != nil {
		return Metadata{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Metadata{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Metadata{}, errors.Wrap(ErrFailedFetch, errors.New(resp.Status))
	}

	var meta Metadata
	if err := json.Unmarshal(body, &meta); err != nil {
		return Metadata{}, err
	}

	return meta, nil
}

func (sdk mfSDK) UpdateThing(t Thing, thingID, token string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s/%s", sdk.thingsURL, thingsEndpoint, thingID)

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

func (sdk mfSDK) DeleteThing(id, token string) error {
	url := fmt.Sprintf("%s/%s/%s", sdk.thingsURL, thingsEndpoint, id)

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

func (sdk mfSDK) DeleteThings(ids []string, token string) error {
	delReq := deleteThingsReq{ThingIDs: ids}
	data, err := json.Marshal(delReq)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s", sdk.thingsURL, thingsEndpoint)

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

func (sdk mfSDK) IdentifyThing(key string) (string, error) {
	idReq := identifyThingReq{Token: key}
	data, err := json.Marshal(idReq)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/%s", sdk.thingsURL, identifyEndpoint)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	resp, err := sdk.sendRequest(req, "", string(CTJSON))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.Wrap(ErrFailedFetch, errors.New(resp.Status))
	}

	var i identifyThingResp
	if err := json.Unmarshal(body, &i); err != nil {
		return "", err
	}

	return i.ID, err
}
