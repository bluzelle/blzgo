package bluzelle

import (
	"encoding/json"
)

type VersionResponseApplicationVersion struct {
	Version string `json:"version"`
}

type VersionResponse struct {
	ApplicationVersion *VersionResponseApplicationVersion `json:"application_version"`
}

func (ctx *Client) Version() (string, error) {
	body, err := ctx.APIQuery("/node_info")
	if err != nil {
		return "", err
	}

	res := &VersionResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return "", err
	}
	return res.ApplicationVersion.Version, nil
}
