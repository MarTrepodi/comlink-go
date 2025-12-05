# {{classname}}

All URIs are relative to */*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetGuild**](GuildApi.md#GetGuild) | **Post** /guild | Get information about a guild
[**GetGuildLeaderboard**](GuildApi.md#GetGuildLeaderboard) | **Post** /getGuildLeaderboard | Get the guild leaderboards
[**GetGuilds**](GuildApi.md#GetGuilds) | **Post** /getGuilds | Search guilds

# **GetGuild**
> InlineResponse2Xx2 GetGuild(ctx, body)
Get information about a guild

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**GuildBody**](GuildBody.md)|  | 

### Return type

[**InlineResponse2Xx2**](inline_response_2XX_2.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetGuildLeaderboard**
> InlineResponse2Xx3 GetGuildLeaderboard(ctx, body)
Get the guild leaderboards

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**GetGuildLeaderboardBody**](GetGuildLeaderboardBody.md)|  | 

### Return type

[**InlineResponse2Xx3**](inline_response_2XX_3.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetGuilds**
> InlineResponse2Xx4 GetGuilds(ctx, body)
Search guilds

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**GetGuildsBody**](GetGuildsBody.md)|  | 

### Return type

[**InlineResponse2Xx4**](inline_response_2XX_4.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

