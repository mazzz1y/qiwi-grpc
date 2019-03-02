# Protocol Documentation
<a name="top"/>

## Table of Contents

- [qiwi.proto](#qiwi.proto)
    - [AddAccountRequest](#protobuf.AddAccountRequest)
    - [AddAccountResponse](#protobuf.AddAccountResponse)
    - [GetAccountBalancesRequest](#protobuf.GetAccountBalancesRequest)
    - [GetAccountBalancesResponse](#protobuf.GetAccountBalancesResponse)
    - [GetAccountBalancesResponse.BalancesEntry](#protobuf.GetAccountBalancesResponse.BalancesEntry)
    - [ListAccountsRequest](#protobuf.ListAccountsRequest)
    - [ListAccountsResponse](#protobuf.ListAccountsResponse)
    - [SendMoneyToQiwiRequest](#protobuf.SendMoneyToQiwiRequest)
    - [SendMoneyToQiwiResponse](#protobuf.SendMoneyToQiwiResponse)
  
  
  
    - [Qiwi](#protobuf.Qiwi)
  

- [Scalar Value Types](#scalar-value-types)



<a name="qiwi.proto"/>
<p align="right"><a href="#top">Top</a></p>

## qiwi.proto



<a name="protobuf.AddAccountRequest"/>

### AddAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token | [string](#string) |  |  |






<a name="protobuf.AddAccountResponse"/>

### AddAccountResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contractID | [int64](#int64) |  |  |






<a name="protobuf.GetAccountBalancesRequest"/>

### GetAccountBalancesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contractID | [int64](#int64) |  |  |






<a name="protobuf.GetAccountBalancesResponse"/>

### GetAccountBalancesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| balances | [GetAccountBalancesResponse.BalancesEntry](#protobuf.GetAccountBalancesResponse.BalancesEntry) | repeated |  |






<a name="protobuf.GetAccountBalancesResponse.BalancesEntry"/>

### GetAccountBalancesResponse.BalancesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [double](#double) |  |  |






<a name="protobuf.ListAccountsRequest"/>

### ListAccountsRequest







<a name="protobuf.ListAccountsResponse"/>

### ListAccountsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contractIDs | [int64](#int64) | repeated |  |






<a name="protobuf.SendMoneyToQiwiRequest"/>

### SendMoneyToQiwiRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| senderContractID | [int64](#int64) |  |  |
| receiverContractID | [string](#string) |  |  |
| currency | [string](#string) |  |  |
| amount | [double](#double) |  |  |






<a name="protobuf.SendMoneyToQiwiResponse"/>

### SendMoneyToQiwiResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [string](#string) |  |  |





 

 

 


<a name="protobuf.Qiwi"/>

### Qiwi


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| AddAccount | [AddAccountRequest](#protobuf.AddAccountRequest) | [AddAccountResponse](#protobuf.AddAccountRequest) |  |
| ListAccounts | [ListAccountsRequest](#protobuf.ListAccountsRequest) | [ListAccountsResponse](#protobuf.ListAccountsRequest) |  |
| GetAccountBalances | [GetAccountBalancesRequest](#protobuf.GetAccountBalancesRequest) | [GetAccountBalancesResponse](#protobuf.GetAccountBalancesRequest) |  |
| SendMoneyToQiwi | [SendMoneyToQiwiRequest](#protobuf.SendMoneyToQiwiRequest) | [SendMoneyToQiwiResponse](#protobuf.SendMoneyToQiwiRequest) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

