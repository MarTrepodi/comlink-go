# {{classname}}

All URIs are relative to */*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetLeaderboard**](GacApi.md#GetLeaderboard) | **Post** /getLeaderboard | View GAC leaderboards.  Supports both the global leaderboards, as well as the leaderboards for the GAC groupings of 8.

# **GetLeaderboard**
> InlineResponse2Xx5 GetLeaderboard(ctx, body)
View GAC leaderboards.  Supports both the global leaderboards, as well as the leaderboards for the GAC groupings of 8.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**GetLeaderboardBody**](GetLeaderboardBody.md)|  | 

### Return type

[**InlineResponse2Xx5**](inline_response_2XX_5.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

