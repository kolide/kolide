package service

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/kolide/fleet/server/kolide"
	"github.com/pkg/errors"
)

// ApplyPackSpecs sends the list of Packs to be applied (upserted) to the
// Fleet instance.
func (c *Client) ApplyPackSpecs(specs []*kolide.PackSpec) error {
	req := applyPackSpecsRequest{Specs: specs}
	response, err := c.AuthenticatedDo("POST", "/api/v1/kolide/spec/packs", req)
	if err != nil {
		return errors.Wrap(err, "POST /api/v1/kolide/spec/packs")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.Errorf("apply pack spec got HTTP %d, expected 200", response.StatusCode)
	}

	var responseBody applyPackSpecsResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return errors.Wrap(err, "decode apply pack spec response")
	}

	if responseBody.Err != nil {
		return errors.Errorf("apply pack spec: %s", responseBody.Err)
	}

	return nil
}

// GetPackSpecs retrieves the list of all Packs.
func (c *Client) GetPackSpecs(specs []*kolide.PackSpec) ([]*kolide.PackSpec, error) {
	req := applyPackSpecsRequest{Specs: specs}
	response, err := c.AuthenticatedDo("GET", "/api/v1/kolide/spec/packs", req)
	if err != nil {
		return nil, errors.Wrap(err, "GET /api/v1/kolide/spec/packs")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.Errorf("get pack spec got HTTP %d, expected 200", response.StatusCode)
	}

	var responseBody getPackSpecsResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, errors.Wrap(err, "decode get pack spec response")
	}

	if responseBody.Err != nil {
		return nil, errors.Errorf("get pack spec: %s", responseBody.Err)
	}

	return responseBody.Specs, nil
}

// DeletePack deletes the pack with the matching name.
func (c *Client) DeletePack(name string) error {
	verb, path := "DELETE", "/api/v1/kolide/packs/"+url.QueryEscape(name)
	response, err := c.AuthenticatedDo(verb, path, nil)
	if err != nil {
		return errors.Wrapf(err, "%s %s", verb, path)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.Errorf("get pack spec got HTTP %d, expected 200", response.StatusCode)
	}

	var responseBody deletePackResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return errors.Wrap(err, "decode get pack spec response")
	}

	if responseBody.Err != nil {
		return errors.Errorf("get pack spec: %s", responseBody.Err)
	}

	return nil
}
