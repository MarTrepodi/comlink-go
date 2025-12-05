# {{classname}}

All URIs are relative to */*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetPlayer**](PlayerApi.md#GetPlayer) | **Post** /player | Get a player profile
[**GetPlayerArenaProfile**](PlayerApi.md#GetPlayerArenaProfile) | **Post** /playerArena | Get a player&#x27;s arena profile

# **GetPlayer**
> InlineResponse2Xx8 GetPlayer(ctx, body)
Get a player profile

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**PlayerBody**](PlayerBody.md)|  | 

### Return type

[**InlineResponse2Xx8**](inline_response_2XX_8.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetPlayerArenaProfile**
> InlineResponse2Xx9 GetPlayerArenaProfile(ctx, body)
Get a player's arena profile

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**PlayerArenaBody**](PlayerArenaBody.md)|  | 

### Return type

[**InlineResponse2Xx9**](inline_response_2XX_9.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

