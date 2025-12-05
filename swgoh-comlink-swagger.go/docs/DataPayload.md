# DataPayload

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Version** | [***Object**](.md) |  | [default to null]
**IncludePveUnits** | [***Object**](.md) | Whether to include pve units in the collections or not. | [optional] [default to null]
**DevicePlatform** | [***Object**](.md) |  | [optional] [default to Android]
**RequestSegment** | [***Object**](.md) | Mutually exclusive with the items parameter. Uses the game data segment enums. | [optional] [default to null]
**Items** | [***Object**](.md) | Mutually exclusive with the requestSegment parameter. Uses a binary bit mask converted to integer string where each collection to include in the response is represented by its digit being set to 1 to include, 0 to not include. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

