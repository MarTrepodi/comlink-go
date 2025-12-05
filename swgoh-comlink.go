// Package swgoh_comlink provides a Go interface library for swgoh-comlink (https://github.com/swgoh-utils/swgoh-comlink)
package swgoh_comlink

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Version of the library
var Version string

// Logger interface for logging
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

// urlPortRe is a compiled regular expression for URL validation
var urlPortRe = regexp.MustCompile(`^https://\S+:(\d+)$`)

// _getPlayerPayload is a helper function to build payload for get_player functions
// allycode: player allyCode
// playerID: player game ID
// enums: boolean
// return: map
func _getPlayerPayload(allycode interface{}, playerID string, enums bool) map[string]interface{} {
	payload := map[string]interface{}{
		"payload": map[string]interface{}{},
		"enums":   enums,
	}

	// If player ID is provided use that instead of allyCode
	if playerID != "" && allycode == nil {
		payload["payload"].(map[string]interface{})["playerId"] = playerID
	} else {
		// Otherwise use allyCode to lookup player data
		var allycodeStr string
		switch v := allycode.(type) {
		case string:
			allycodeStr = v
		case int:
			allycodeStr = strconv.Itoa(v)
		case int64:
			allycodeStr = strconv.FormatInt(v, 10)
		default:
			allycodeStr = fmt.Sprintf("%v", v)
		}
		payload["payload"].(map[string]interface{})["allyCode"] = allycodeStr
	}

	return payload
}

// sanitizeURL makes sure provided URL is in the expected format and returns sanitized version
func sanitizeURL(url string) string {
	url = strings.TrimRight(url, "/")
	if strings.HasPrefix(url, "https") && !urlPortRe.MatchString(url) {
		url = fmt.Sprintf("%s:443", url)
	}
	return url
}

// SwgohComlink is the main class for swgoh-comlink interface and supported methods.
// Instances of this struct are used to query the Star Wars Galaxy of Heroes
// game servers for exposed endpoints via the swgoh-comlink proxy library
// running on the same host.
type SwgohComlink struct {
	Version      string
	URLBase      string
	StatsURLBase string
	HMAC         bool
	AccessKey    string
	SecretKey    string
	Protocol     string
	Client       *http.Client
	Logger       Logger
}

// NewSwgohComlink creates a new SwgohComlink instance
// url: The URL where swgoh-comlink is running. Defaults to 'http://localhost:3000'
// statsURL: The url of the swgoh-stats service (if used), such as 'http://localhost:3223'
// accessKey: The HMAC public key. Default to empty which indicates HMAC is not used.
// secretKey: The HMAC private key. Default to empty which indicates HMAC is not used.
// host: IP address or DNS name of server where the swgoh-comlink service is running
// port: TCP port number where the swgoh-comlink service is running [Default: 3000]
// statsPort: TCP port number of where the comlink-stats service is running [Default: 3223]
//
// Notes:
//
//	If the 'host' and 'port' parameters are provided, the 'url' and 'statsURL' parameters are ignored.
func NewSwgohComlink(url, statsURL, accessKey, secretKey, host string, port, statsPort int) *SwgohComlink {
	sc := &SwgohComlink{
		Version:      Version,
		URLBase:      sanitizeURL(url),
		StatsURLBase: sanitizeURL(statsURL),
		HMAC:         false,
		Protocol:     "http",
	}

	// host and port parameters override defaults
	if host != "" {
		sc.URLBase = sc.Protocol + fmt.Sprintf("://%s:%d", host, port)
		sc.StatsURLBase = sc.Protocol + fmt.Sprintf("://%s:%d", host, statsPort)
	}

	// Use values passed from client first, otherwise check for environment variables
	if accessKey != "" {
		sc.AccessKey = accessKey
	} else if envAccessKey := os.Getenv("ACCESS_KEY"); envAccessKey != "" {
		sc.AccessKey = envAccessKey
	}

	if secretKey != "" {
		sc.SecretKey = secretKey
	} else if envSecretKey := os.Getenv("SECRET_KEY"); envSecretKey != "" {
		sc.SecretKey = envSecretKey
	}

	if sc.AccessKey != "" && sc.SecretKey != "" {
		sc.HMAC = true
	}

	// Create HTTP client with disabled SSL verification
	sc.Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return sc
}

// _getGameVersion retrieves the game version from metadata
func (sc *SwgohComlink) _getGameVersion() (string, error) {
	md, err := sc.GetGameMetadata(nil, false)
	if err != nil {
		return "", err
	}
	if latestVersion, ok := md["latestGamedataVersion"].(string); ok {
		return latestVersion, nil
	}
	return "", fmt.Errorf("latestGamedataVersion not found in metadata")
}

// _post performs a POST request to the specified endpoint
func (sc *SwgohComlink) _post(urlBase, endpoint string, payload interface{}) (map[string]interface{}, error) {
	if urlBase == "" {
		urlBase = sc.URLBase
	}
	postURL := urlBase + "/" + endpoint

	reqHeaders := make(map[string]string)

	var payloadBytes []byte
	var err error

	// If access_key and secret_key are set, perform HMAC security
	if sc.HMAC {
		reqTime := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
		reqHeaders["X-Date"] = reqTime

		h := hmac.New(sha256.New, []byte(sc.SecretKey))
		h.Write([]byte(reqTime))
		h.Write([]byte("POST"))
		h.Write([]byte("/" + endpoint))

		// json dumps separators needed for compact string formatting required for compatibility with
		// comlink since it is written with javascript as the primary object model
		// ordered dicts are also required with the 'payload' key listed first for proper MD5 hash calculation
		if payload != nil {
			payloadBytes, err = json.Marshal(payload)
			if err != nil {
				return nil, &SwgohComlinkException{Message: err.Error()}
			}
		} else {
			payloadBytes = []byte("{}")
		}

		payloadHashDigest := md5.Sum(payloadBytes)
		payloadHashHex := hex.EncodeToString(payloadHashDigest[:])
		h.Write([]byte(payloadHashHex))

		hmacDigest := hex.EncodeToString(h.Sum(nil))
		reqHeaders["Authorization"] = fmt.Sprintf("HMAC-SHA256 Credential=%s,Signature=%s", sc.AccessKey, hmacDigest)
	} else {
		if payload != nil {
			payloadBytes, err = json.Marshal(payload)
			if err != nil {
				return nil, &SwgohComlinkException{Message: err.Error()}
			}
		} else {
			payloadBytes = []byte("{}")
		}
	}

	req, err := http.NewRequest("POST", postURL, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return nil, &SwgohComlinkException{Message: err.Error()}
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range reqHeaders {
		req.Header.Set(key, value)
	}

	resp, err := sc.Client.Do(req)
	if err != nil {
		return nil, &SwgohComlinkException{Message: err.Error()}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &SwgohComlinkException{Message: err.Error()}
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, &SwgohComlinkException{Message: err.Error()}
	}

	return result, nil
}

// GetUnitStats calculates unit stats using the swgoh-stats service interface to swgoh-comlink.
//
// This method communicates with an external StatCalc container/service to calculate stats for the
// given units. If the unit for which stats are being calculated is a ship, the 'requestPayload' *MUST*
// include crew members along with the ship unit.
//
// The most common use of this method is to provide an entire player roster as the 'requestPayload' argument. This
// action will calculate the stats for each character and ship in the player roster. The resulting stats are
// included in a new 'stats' key in the result object.
//
// requestPayload: map or slice of maps containing units for which to calculate stats.
// flags: Slice of strings specifying which flags to include in the request URI.
// language: String indicating the desired localized language. Default "eng_us".
//
// Returns: Input object with the calculated stats for the specified units included.
func (sc *SwgohComlink) GetUnitStats(requestPayload interface{}, flags []string, language string) (interface{}, error) {
	// Define the flags that StatCalc understands
	allowedFlags := map[string]bool{
		"gameStyle":      true,
		"calcGP":         true,
		"onlyGP":         true,
		"withoutModCalc": true,
		"percentVals":    true,
		"useMax":         true,
		"scaled":         true,
		"unscaled":       true,
		"statIDs":        true,
		"enums":          true,
		"noSpace":        true,
	}

	var queryString string
	var flagStr string

	if flags != nil {
		// Validate flags
		for _, flag := range flags {
			if !allowedFlags[flag] {
				allowedFlagsStr := "{"
				for k := range allowedFlags {
					allowedFlagsStr += k + ", "
				}
				allowedFlagsStr = strings.TrimSuffix(allowedFlagsStr, ", ") + "}"
				return nil, &SwgohComlinkValueError{
					Message: fmt.Sprintf("Invalid argument. <flags> should be a list of strings with one or more of \"%s flag values.", allowedFlagsStr),
				}
			}
		}
		flagStr = "flags=" + strings.Join(flags, ",")
	}

	var languageStr string
	if language != "" {
		languageStr = "language=" + language
	}

	if flagStr != "" || languageStr != "" {
		parts := []string{}
		if flagStr != "" {
			parts = append(parts, flagStr)
		}
		if languageStr != "" {
			parts = append(parts, languageStr)
		}
		queryString = "?" + strings.Join(parts, "&")
	}

	endpointString := "api"
	if queryString != "" {
		endpointString = "api" + queryString
	}

	// Convert map to slice if needed
	var payloadToSend interface{}
	if payloadMap, ok := requestPayload.(map[string]interface{}); ok {
		payloadToSend = []interface{}{payloadMap}
	} else {
		payloadToSend = requestPayload
	}

	result, err := sc._post(sc.StatsURLBase, endpointString, payloadToSend)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetEnums gets an object containing the game data enums
// Returns: A map containing the game data enums.
func (sc *SwgohComlink) GetEnums() (map[string]interface{}, error) {
	url := sc.URLBase + "/enums"
	resp, err := sc.Client.Get(url)
	if err != nil {
		return nil, &SwgohComlinkException{Message: err.Error()}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &SwgohComlinkException{Message: err.Error()}
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, &SwgohComlinkException{Message: err.Error()}
	}

	return result, nil
}

// Alias for non PEP usage of direct endpoint calls
func (sc *SwgohComlink) getEnums() (map[string]interface{}, error) {
	return sc.GetEnums()
}

// GetEvents gets an object containing the events game data
// enums: Boolean flag to indicate whether enum value should be converted in response. [Default is False]
// Returns: A single element map containing a list of the events game data.
func (sc *SwgohComlink) GetEvents(enums bool) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"payload": map[string]interface{}{},
		"enums":   enums,
	}
	return sc._post("", "getEvents", payload)
}

// Alias for non PEP usage of direct endpoint calls
func (sc *SwgohComlink) getEvents(enums bool) (map[string]interface{}, error) {
	return sc.GetEvents(enums)
}

// GetGameData gets game data
// version: string (found in metadata key value 'latestGamedataVersion')
// includePveUnits: boolean [Defaults to True]
// requestSegment: integer >=0 [Defaults to 0]
// enums: boolean [Defaults to False]
// items: string [Defaults to empty] bitwise value indicating the collections to retrieve from game.
//
//	NOTE: this parameter is mutually exclusive with request_segment.
//
// devicePlatform: string [Defaults to "Android"]
// Returns: A map containing the game data.
func (sc *SwgohComlink) GetGameData(version string, includePveUnits bool, requestSegment int, enums bool, items string, devicePlatform string) (map[string]interface{}, error) {
	var gameVersion string
	var err error

	if version == "" {
		gameVersion, err = sc._getGameVersion()
		if err != nil {
			return nil, err
		}
	} else {
		gameVersion = version
	}

	payload := map[string]interface{}{
		"payload": map[string]interface{}{
			"version":         gameVersion,
			"devicePlatform":  devicePlatform,
			"includePveUnits": includePveUnits,
		},
		"enums": enums,
	}

	if items != "" {
		// presence of 'items' argument overrides the 'request_segment' argument
		// Check if items is a numeric string
		if itemsInt, err := strconv.Atoi(items); err == nil {
			payload["payload"].(map[string]interface{})["items"] = strconv.Itoa(absInt(itemsInt))
		} else {
			// Try to get from Constants
			constantValue := GetConstant(items)
			if constantValue != "" {
				payload["payload"].(map[string]interface{})["items"] = constantValue
			} else {
				payload["payload"].(map[string]interface{})["items"] = "-1"
			}
		}
	} else {
		if requestSegment < 0 || requestSegment > 4 {
			return nil, &SwgohComlinkValueError{
				Message: "Invalid argument. <request_segment> should be an integer between 0 and 4, inclusive.",
			}
		}
		payload["payload"].(map[string]interface{})["requestSegment"] = requestSegment
	}

	return sc._post("", "data", payload)
}

// Alias for non PEP usage of direct endpoint calls
func (sc *SwgohComlink) getGameData(version string, includePveUnits bool, requestSegment int, enums bool, items string, devicePlatform string) (map[string]interface{}, error) {
	return sc.GetGameData(version, includePveUnits, requestSegment, enums, items, devicePlatform)
}

// GetLocalization gets localization data from game
// id: latestLocalizationBundleVersion found in game metadata. This method will collect the latest language
//
//	version if the 'id' argument is not provided.
//
// locale: string Specify only a specific locale to retrieve [for example "ENG_US"]
// unzip: boolean [Defaults to False]
// enums: boolean [Defaults to False]
// Returns: A map containing the localization data.
func (sc *SwgohComlink) GetLocalization(id, locale string, unzip, enums bool) (map[string]interface{}, error) {
	if id == "" {
		currentGameVersion, err := sc.GetLatestGameDataVersion()
		if err != nil {
			return nil, err
		}
		if language, ok := currentGameVersion["language"].(string); ok {
			id = language
		}
	}

	if locale != "" {
		id = id + ":" + strings.ToUpper(locale)
	}

	payload := map[string]interface{}{
		"unzip": unzip,
		"enums": enums,
		"payload": map[string]interface{}{
			"id": id,
		},
	}

	return sc._post("", "localization", payload)
}

// Aliases for non PEP usage of direct endpoint calls
func (sc *SwgohComlink) getLocalization(id, locale string, unzip, enums bool) (map[string]interface{}, error) {
	return sc.GetLocalization(id, locale, unzip, enums)
}

func (sc *SwgohComlink) GetLocalizationBundle(id, locale string, unzip, enums bool) (map[string]interface{}, error) {
	return sc.GetLocalization(id, locale, unzip, enums)
}

func (sc *SwgohComlink) getLocalizationBundle(id, locale string, unzip, enums bool) (map[string]interface{}, error) {
	return sc.GetLocalization(id, locale, unzip, enums)
}

// GetGameMetadata gets the game metadata. Game metadata contains the current game and localization versions.
// clientSpecs: Optional map containing client specifications
// enums: Boolean signifying whether enums in response should be translated to text. [Default: False]
// Returns: A map containing the game metadata.
func (sc *SwgohComlink) GetGameMetadata(clientSpecs map[string]interface{}, enums bool) (map[string]interface{}, error) {
	var payload map[string]interface{}

	if clientSpecs != nil {
		payload = map[string]interface{}{
			"payload": map[string]interface{}{
				"client_specs": clientSpecs,
			},
			"enums": enums,
		}
	} else {
		payload = map[string]interface{}{}
	}

	return sc._post("", "metadata", payload)
}

// Aliases for non PEP usage of direct endpoint calls
func (sc *SwgohComlink) getGameMetaData(clientSpecs map[string]interface{}, enums bool) (map[string]interface{}, error) {
	return sc.GetGameMetadata(clientSpecs, enums)
}

func (sc *SwgohComlink) GetMetaData(clientSpecs map[string]interface{}, enums bool) (map[string]interface{}, error) {
	return sc.GetGameMetadata(clientSpecs, enums)
}

func (sc *SwgohComlink) GetMetadata(clientSpecs map[string]interface{}, enums bool) (map[string]interface{}, error) {
	return sc.GetGameMetadata(clientSpecs, enums)
}

func (sc *SwgohComlink) getMetaData(clientSpecs map[string]interface{}, enums bool) (map[string]interface{}, error) {
	return sc.GetGameMetadata(clientSpecs, enums)
}

// GetPlayer gets player information from game. Either allycode or playerID must be provided.
// allycode: integer or string representing player allycode
// playerID: string representing player game ID
// enums: boolean [Defaults to False]
// Returns: A map containing the player information.
func (sc *SwgohComlink) GetPlayer(allycode interface{}, playerID string, enums bool) (map[string]interface{}, error) {
	payload := _getPlayerPayload(allycode, playerID, enums)
	return sc._post("", "player", payload)
}

// Alias for non PEP usage of direct endpoint calls
func (sc *SwgohComlink) getPlayer(allycode interface{}, playerID string, enums bool) (map[string]interface{}, error) {
	return sc.GetPlayer(allycode, playerID, enums)
}

// GetPlayerArena gets player arena information from game. Either allycode or playerID must be provided.
// allycode: integer or string representing player allycode
// playerID: string representing player game ID
// playerDetailsOnly: filter results to only player details [Defaults to False]
// enums: boolean [Defaults to False]
// Returns: A map containing the player arena information.
func (sc *SwgohComlink) GetPlayerArena(allycode interface{}, playerID string, playerDetailsOnly, enums bool) (map[string]interface{}, error) {
	payload := _getPlayerPayload(allycode, playerID, enums)
	payload["payload"].(map[string]interface{})["playerDetailsOnly"] = playerDetailsOnly
	return sc._post("", "playerArena", payload)
}

// GetGuild gets guild information for a specific Guild ID.
// guildID: String ID of guild to retrieve. Guild ID can be found in the output
//
//	of the GetPlayer() call. (Required)
//
// includeRecentGuildActivityInfo: boolean [Default: False] (Optional)
// enums: Should enums in response be translated to text. [Default: False] (Optional)
// Returns: A map containing the guild information.
func (sc *SwgohComlink) GetGuild(guildID string, includeRecentGuildActivityInfo, enums bool) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"payload": map[string]interface{}{
			"guildId":                        guildID,
			"includeRecentGuildActivityInfo": includeRecentGuildActivityInfo,
		},
		"enums": enums,
	}

	guild, err := sc._post("", "guild", payload)
	if err != nil {
		return nil, err
	}

	if guildData, ok := guild["guild"]; ok {
		if guildMap, ok := guildData.(map[string]interface{}); ok {
			return guildMap, nil
		}
	}

	return guild, nil
}

// Alias for non PEP usage of direct endpoint calls
func (sc *SwgohComlink) getGuild(guildID string, includeRecentGuildActivityInfo, enums bool) (map[string]interface{}, error) {
	return sc.GetGuild(guildID, includeRecentGuildActivityInfo, enums)
}

// GetGuildsByName searches for guild by name and return match.
// name: string for guild name search
// startIndex: integer representing where in the resulting list of guild name matches
//
//	the return object should begin
//
// count: integer representing the maximum number of matches to return, [Default: 10]
// enums: Whether to translate enums in response to text, [Default: False]
// Returns: A map containing the guild search results.
func (sc *SwgohComlink) GetGuildsByName(name string, startIndex, count int, enums bool) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"payload": map[string]interface{}{
			"name":       name,
			"filterType": 4,
			"startIndex": startIndex,
			"count":      count,
		},
		"enums": enums,
	}

	return sc._post("", "getGuilds", payload)
}

// Alias for non PEP usage of direct endpoint calls
func (sc *SwgohComlink) getGuildByName(name string, startIndex, count int, enums bool) (map[string]interface{}, error) {
	return sc.GetGuildsByName(name, startIndex, count, enums)
}

// GetGuildsByCriteria searches for guild by guild criteria and return matches.
// searchCriteria: Map containing search criteria
// startIndex: integer representing where in the resulting list of guild name matches the return object
//
//	should begin
//
// count: integer representing the maximum number of matches to return
// enums: Whether to translate enum values to text [Default: False]
// Returns: A map containing the guild search results.
func (sc *SwgohComlink) GetGuildsByCriteria(searchCriteria map[string]interface{}, startIndex, count int, enums bool) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"payload": map[string]interface{}{
			"searchCriteria": searchCriteria,
			"filterType":     5,
			"startIndex":     startIndex,
			"count":          count,
		},
		"enums": enums,
	}

	return sc._post("", "getGuilds", payload)
}

// Alias for non PEP usage of direct endpoint calls
func (sc *SwgohComlink) getGuildByCriteria(searchCriteria map[string]interface{}, startIndex, count int, enums bool) (map[string]interface{}, error) {
	return sc.GetGuildsByCriteria(searchCriteria, startIndex, count, enums)
}

// GetLeaderboard retrieves Grand Arena Championship leaderboard information.
// leaderboardType (int): Type 4 is for scanning gac brackets, and only returns results while an event is active.
//
//	When type 4 is indicated, the "league" and "division" arguments must also be provided.
//	Type 6 is for the global leaderboards for the league + divisions.
//	When type 6 is indicated, the "eventInstanceID" and "groupID" must also be provided.
//
// league (int|string): Enum values 20, 40, 60, 80, and 100 correspond to carbonite, bronzium, chromium, aurodium,
//
//	and kyber respectively. Also accepts string values for each league.
//
// division (int|string): Enum values 5, 10, 15, 20, and 25 correspond to divisions 5 through 1 respectively.
//
//	Also accepts string or int values for each division.
//
// eventInstanceID (string): When leaderboardType 4 is indicated, a combination of the event Id and
//
//	the instance ID separated by ':'
//	Example: CHAMPIONSHIPS_GRAND_ARENA_GA2_EVENT_SEASON_36:O1675202400000
//
// groupID (string): When leaderboardType 4 is indicated, must start with the same eventInstanceId, followed
//
//	by the league and bracketId, separated by ':'. The number at the end is the bracketId, and
//	goes from 0 to N, where N is the last group of 8 players.
//	Example: CHAMPIONSHIPS_GRAND_ARENA_GA2_EVENT_SEASON_36:O1675202400000:CARBONITE:10431
//
// enums (bool): Whether to translate enum values to text [Default: False]
// Returns: A map containing the leaderboard data.
func (sc *SwgohComlink) GetLeaderboard(leaderboardType int, league interface{}, division interface{}, eventInstanceID, groupID string, enums bool) (map[string]interface{}, error) {
	leagues := map[string]int{
		"kyber":     100,
		"aurodium":  80,
		"chromium":  60,
		"bronzium":  40,
		"carbonite": 20,
	}

	divisions := map[string]int{
		"1": 25,
		"2": 20,
		"3": 15,
		"4": 10,
		"5": 5,
	}

	// Translate parameters if needed
	var leagueInt int
	if leagueStr, ok := league.(string); ok {
		leagueInt = leagues[strings.ToLower(leagueStr)]
	} else if leagueVal, ok := league.(int); ok {
		leagueInt = leagueVal
	}

	var divisionInt int
	if divisionStr, ok := division.(string); ok {
		divisionInt = divisions[strings.ToLower(divisionStr)]
	} else if divisionVal, ok := division.(int); ok {
		// Check if it's a single digit
		if divisionVal >= 1 && divisionVal <= 5 {
			divisionInt = divisions[strconv.Itoa(divisionVal)]
		} else {
			divisionInt = divisionVal
		}
	}

	payload := map[string]interface{}{
		"payload": map[string]interface{}{
			"leaderboardType": leaderboardType,
		},
		"enums": enums,
	}

	if leaderboardType == 4 {
		payload["payload"].(map[string]interface{})["eventInstanceId"] = eventInstanceID
		payload["payload"].(map[string]interface{})["groupId"] = groupID
	} else if leaderboardType == 6 {
		payload["payload"].(map[string]interface{})["league"] = leagueInt
		payload["payload"].(map[string]interface{})["division"] = divisionInt
	}

	return sc._post("", "getLeaderboard", payload)
}

// Aliases for non PEP usage of direct endpoint calls
func (sc *SwgohComlink) getLeaderboard(leaderboardType int, league interface{}, division interface{}, eventInstanceID, groupID string, enums bool) (map[string]interface{}, error) {
	return sc.GetLeaderboard(leaderboardType, league, division, eventInstanceID, groupID, enums)
}

func (sc *SwgohComlink) GetGacLeaderboard(leaderboardType int, league interface{}, division interface{}, eventInstanceID, groupID string, enums bool) (map[string]interface{}, error) {
	return sc.GetLeaderboard(leaderboardType, league, division, eventInstanceID, groupID, enums)
}

func (sc *SwgohComlink) getGacLeaderboard(leaderboardType int, league interface{}, division interface{}, eventInstanceID, groupID string, enums bool) (map[string]interface{}, error) {
	return sc.GetLeaderboard(leaderboardType, league, division, eventInstanceID, groupID, enums)
}

// GetGuildLeaderboard fetches the guild leaderboard data for given leaderboard ID.
//
// This function interacts with an external API to retrieve the leaderboard
// data for the supplied guild leaderboard ID. The user can specify the
// number of entries to fetch.
//
// leaderboardID: slice containing one leaderboard ID map for the data to be fetched.
// count: int, optional - The number of leaderboard entries to retrieve. Defaults to 200.
// enums: bool, optional - Whether or not to translate enum values. Defaults to False.
// Returns: A map containing the guild leaderboard data.
func (sc *SwgohComlink) GetGuildLeaderboard(leaderboardID []interface{}, count int, enums bool) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"payload": map[string]interface{}{
			"leaderboardId": leaderboardID,
			"count":         count,
		},
		"enums": enums,
	}

	return sc._post("", "getGuildLeaderboard", payload)
}

// Alias for non PEP usage of direct endpoint calls
func (sc *SwgohComlink) getGuildLeaderboard(leaderboardID []interface{}, count int, enums bool) (map[string]interface{}, error) {
	return sc.GetGuildLeaderboard(leaderboardID, count, enums)
}

// GetNameSpaces fetches namespaces based on the specified filter criteria.
// *** (PLACEHOLDER) - Actual use is unknown at this time ***
// onlyCompatible: Determines whether to fetch only compatible namespaces.
// enums: Specifies whether enum types should be included in the response.
// Returns: A map containing the information about the retrieved namespaces.
func (sc *SwgohComlink) GetNameSpaces(onlyCompatible, enums bool) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"payload": map[string]interface{}{
			"onlyCompatible": onlyCompatible,
		},
		"enums": enums,
	}

	return sc._post("", "getNameSpaces", payload)
}

// GetLatestGameDataVersion is a helper method to retrieve the latest game data version
func (sc *SwgohComlink) GetLatestGameDataVersion() (map[string]interface{}, error) {
	return sc.GetGameMetadata(nil, false)
}

// absInt returns the absolute value of an integer
func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
