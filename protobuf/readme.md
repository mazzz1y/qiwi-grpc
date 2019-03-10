# Protocol Documentation
<a name="top"/>

## Table of Contents

- [qiwi.proto](#qiwi.proto)
    - [AddAccountRequest](#protobuf.AddAccountRequest)
    - [AddAccountResponse](#protobuf.AddAccountResponse)
    - [DepositCheckRequest](#protobuf.DepositCheckRequest)
    - [DepositCheckResponse](#protobuf.DepositCheckResponse)
    - [DepositRequest](#protobuf.DepositRequest)
    - [DepositResponse](#protobuf.DepositResponse)
    - [GetAccountBalancesRequest](#protobuf.GetAccountBalancesRequest)
    - [GetAccountBalancesResponse](#protobuf.GetAccountBalancesResponse)
    - [ListAccountsRequest](#protobuf.ListAccountsRequest)
    - [ListAccountsResponse](#protobuf.ListAccountsResponse)
  
  
  
    - [Qiwi](#protobuf.Qiwi)
  

- [Scalar Value Types](#scalar-value-types)



<a name="qiwi.proto"/>
<p align="right"><a href="#top">Top</a></p>

## qiwi.proto



<a name="protobuf.AddAccountRequest"/>

### AddAccountRequest
Add/Update account in DB


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token | [string](#string) |  | Qiwi API token: https://qiwi.com/api |
| contractID | [string](#string) |  | contractID, Format: 79999999999 |






<a name="protobuf.AddAccountResponse"/>

### AddAccountResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contractID | [string](#string) |  |  |
| operationLimit | [int64](#int64) |  |  |
| monthLimit | [int64](#int64) |  |  |






<a name="protobuf.DepositCheckRequest"/>

### DepositCheckRequest
Check status of deposit


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | Deposit ID |






<a name="protobuf.DepositCheckResponse"/>

### DepositCheckResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | Deposit ID |
| status | [bool](#bool) |  | Status of deposit, if true it means that money came to the wallet |
| amounts | [int64](#int64) | repeated | Amount per transaction |
| comments | [string](#string) | repeated | Comment per transaction |
| contractIDs | [string](#string) | repeated | ContractID per transaction, Format: 79999999999 |
| links | [string](#string) | repeated | Link for user-friendly payments |
| statuses | [bool](#bool) | repeated | Status per transaction |






<a name="protobuf.DepositRequest"/>

### DepositRequest
Create deposit entity and return requisites for payment


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| amount | [int64](#int64) |  | Amount(RUB) |






<a name="protobuf.DepositResponse"/>

### DepositResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | Deposit ID |
| amounts | [int64](#int64) | repeated | Array of amounts |
| comments | [string](#string) | repeated | Array of comments |
| contractIDs | [string](#string) | repeated | Array of contractIDs |
| links | [string](#string) | repeated | Array of user-friendly links |






<a name="protobuf.GetAccountBalancesRequest"/>

### GetAccountBalancesRequest
Return RUB balance for Qiwi account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contractID | [string](#string) |  | ContractID, Format: 79999999999 |






<a name="protobuf.GetAccountBalancesResponse"/>

### GetAccountBalancesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| balance | [int64](#int64) |  | RUB balance |






<a name="protobuf.ListAccountsRequest"/>

### ListAccountsRequest
Return list of account stored in DB






<a name="protobuf.ListAccountsResponse"/>

### ListAccountsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contractIDs | [string](#string) | repeated | Array of contractIDs, Format: 79999999999 |





 

 

 


<a name="protobuf.Qiwi"/>

### Qiwi


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| AddAccount | [AddAccountRequest](#protobuf.AddAccountRequest) | [AddAccountResponse](#protobuf.AddAccountRequest) |  |
| ListAccounts | [ListAccountsRequest](#protobuf.ListAccountsRequest) | [ListAccountsResponse](#protobuf.ListAccountsRequest) |  |
| GetAccountBalances | [GetAccountBalancesRequest](#protobuf.GetAccountBalancesRequest) | [GetAccountBalancesResponse](#protobuf.GetAccountBalancesRequest) |  |
| Deposit | [DepositRequest](#protobuf.DepositRequest) | [DepositResponse](#protobuf.DepositRequest) |  |
| DepositCheck | [DepositCheckRequest](#protobuf.DepositCheckRequest) | [DepositCheckResponse](#protobuf.DepositCheckRequest) |  |

 



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

