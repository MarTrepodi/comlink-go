package swgohcomlink

// --- /metadata Request Structs (Input) ---

// GetMetaDataClientSpecs corresponds to payload.clientSpecs in the OpenAPI spec.
type GetMetaDataClientSpecs struct {
	Platform        string `json:"platform,omitempty"`
	BundleID        string `json:"bundleId,omitempty"`
	ExternalVersion string `json:"externalVersion,omitempty"`
	InternalVersion string `json:"internalVersion,omitempty"`
	Region          string `json:"region,omitempty"`
}

// GetMetaDataPayload corresponds to the top-level payload object.
type GetMetaDataPayload struct {
	ClientSpecs *GetMetaDataClientSpecs `json:"clientSpecs,omitempty"`
}

// GetMetaDataRequest is the request body for the /metadata endpoint.
type GetMetaDataRequest struct {
	Payload *GetMetaDataPayload `json:"payload,omitempty"`
	Enums   bool                `json:"enums,omitempty"` // default: false
}

// --- /metadata Response Structs (Output) ---

// ConfigEntry is a simplified model for an item in the config array.
type ConfigEntry struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

// GetMetaDataResponse is the response body for the /metadata endpoint.
type GetMetaDataResponse struct {
	Config                          []ConfigEntry `json:"config,omitempty"`
	AssetVersion                    int32         `json:"assetVersion,omitempty"`
	AssetSubpath                    string        `json:"assetSubpath,omitempty"`
	ServerTimestamp                 int64         `json:"serverTimestamp,omitempty"`
	LatestLocalizationBundleVersion string        `json:"latestLocalizationBundleVersion,omitempty"`
	LatestGamedataVersion           string        `json:"latestGamedataVersion,omitempty"`
	// A complete definition would include all ~10 fields from the OpenAPI spec.
}

// GetMetaData fetches the game metadata.
// POST /metadata
func (c *Client) GetMetaData(req *GetMetaDataRequest) (*GetMetaDataResponse, error) {
	if req == nil {
		req = &GetMetaDataRequest{}
	}

	resp := &GetMetaDataResponse{}
	err := c.request("POST", "/metadata", req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
