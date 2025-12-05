package swgohcomlink

import "errors"

// --- /player Request Structs (Input) ---

// GetPlayerPayload requires exactly one of AllyCode or PlayerID.
type GetPlayerPayload struct {
	AllyCode          string `json:"allyCode,omitempty"`          // pattern: ^[1-9]{9}$
	PlayerID          string `json:"playerId,omitempty"`          // pattern: ^[A-Za-z0-9\-_]{22}$
	PlayerDetailsOnly bool   `json:"playerDetailsOnly,omitempty"` // default: false
}

// GetPlayerRequest is the request body for the /player endpoint.
type GetPlayerRequest struct {
	Payload GetPlayerPayload `json:"payload"`
	Enums   bool             `json:"enums,omitempty"` // default: false
}

// Validate ensures that exactly one of AllyCode or PlayerID is set, as required by the spec.
func (r *GetPlayerRequest) Validate() error {
	hasAllyCode := r.Payload.AllyCode != ""
	hasPlayerID := r.Payload.PlayerID != ""

	if hasAllyCode && hasPlayerID {
		return errors.New("must include either allyCode or playerId, but not both")
	}
	if !hasAllyCode && !hasPlayerID {
		return errors.New("must include either allyCode or playerId")
	}
	return nil
}

// --- /player Response Structs (Output) ---

// PlayerProfileStat is a simplified model for a player's profile statistic.
type PlayerProfileStat struct {
	NameKey string `json:"nameKey"`
	Value   string `json:"value"`
	// A complete definition would include fields like StatID
}

// Unit is a simplified model for a player's roster unit.
type Unit struct {
	DefID  string `json:"defId"`
	Rarity int32  `json:"rarity"`
	// A complete definition would include fields like RelicTier, GearLevel, etc.
}

// GetPlayerResponse is the response body for the /player endpoint.
type GetPlayerResponse struct {
	Name        string              `json:"name"`
	Level       int32               `json:"level"`
	AllyCode    int64               `json:"allyCode"`
	PlayerID    string              `json:"playerId"`
	RosterUnit  []Unit              `json:"rosterUnit"`
	ProfileStat []PlayerProfileStat `json:"profileStat"`
	GuildID     string              `json:"guildId"`
	GuildName   string              `json:"guildName"`
	// A complete definition would include all ~30 fields from the OpenAPI spec.
}

// GetPlayer fetches a player profile by Ally Code or Player ID.
// POST /player
func (c *Client) GetPlayer(req *GetPlayerRequest) (*GetPlayerResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &GetPlayerResponse{}
	err := c.request("POST", "/player", req.Payload, resp) // Note: request expects payload as the body
	if err != nil {
		return nil, err
	}

	return resp, nil
}
