# {{classname}}

All URIs are relative to */*

Method | HTTP request | Description
------------- | ------------- | -------------
[**EnumsGet**](GameDataApi.md#EnumsGet) | **Get** /enums | Get an object containing all of the game data enums
[**GetGameData**](GameDataApi.md#GetGameData) | **Post** /data | Get the game data
[**GetLocalizationBundle**](GameDataApi.md#GetLocalizationBundle) | **Post** /localization | Get the localization bundle
[**GetMetaData**](GameDataApi.md#GetMetaData) | **Post** /metadata | Get the game metadata
[**GetNameSpaces**](GameDataApi.md#GetNameSpaces) | **Post** /getNameSpaces | Get the name spaces for segmented content
[**GetSegmentedContent**](GameDataApi.md#GetSegmentedContent) | **Post** /getSegmentedContent | Get the segmented content details

# **EnumsGet**
> EnumsGet(ctx, )
Get an object containing all of the game data enums

### Required Parameters
This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetGameData**
> InlineResponse2Xx1 GetGameData(ctx, body)
Get the game data

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**DataBody**](DataBody.md)|  | 

### Return type

[**InlineResponse2Xx1**](inline_response_2XX_1.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetLocalizationBundle**
> Object GetLocalizationBundle(ctx, body)
Get the localization bundle

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**LocalizationBody**](LocalizationBody.md)|  | 

### Return type

**Object**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetMetaData**
> InlineResponse2Xx6 GetMetaData(ctx, optional)
Get the game metadata

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GameDataApiGetMetaDataOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GameDataApiGetMetaDataOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of MetadataBody**](MetadataBody.md)|  | 

### Return type

[**InlineResponse2Xx6**](inline_response_2XX_6.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetNameSpaces**
> InlineResponse2Xx7 GetNameSpaces(ctx, optional)
Get the name spaces for segmented content

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GameDataApiGetNameSpacesOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GameDataApiGetNameSpacesOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of GetNameSpacesBody**](GetNameSpacesBody.md)|  | 

### Return type

[**InlineResponse2Xx7**](inline_response_2XX_7.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetSegmentedContent**
> InlineResponse2Xx10 GetSegmentedContent(ctx, optional)
Get the segmented content details

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GameDataApiGetSegmentedContentOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GameDataApiGetSegmentedContentOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of GetSegmentedContentBody**](GetSegmentedContentBody.md)|  | 

### Return type

[**InlineResponse2Xx10**](inline_response_2XX_10.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

