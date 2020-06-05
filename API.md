# Trollr: A HTTP/JSON API for Troll
Trollr is a simple wrapper around the amazing Troll: A dice roll language and calculator created by Torben Mogensen.
The wrapper simply exposes and HTTP/JSON server that executes Troll, parses the results and returns it. The server
has some built in rate-limiting and pooling to prevent abuse. I created this small server to support a Discord bot
that I am working on.

## Version: 0.1.0-alpha

**Contact information:**  
Ben Doerr  
https://trollr.live  
craftsman@bendoer.me  

**License:** [MIT](http://opensource.org/licenses/MIT)

### /calc

#### GET
##### Summary:

Calculate the probabilities of dice roll.

##### Description:

Given a roll definition this will delegate the roll to Troll and return the
probabilities structured as JSON.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| d | query | The Troll roll definition. This can passed as the query parameter 'd' or in the request body. | No | string |
| d | body | The Troll roll definition. This can passed as the query parameter 'd' or in the request body. | No | string |
| c | query | What kind of cumulative probabilities you would like. One of 'ge' (default), 'gt', 'le', or 'lt'. | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | The probabilities of rolling the dice | [CalcResult](#calcresult) |
| 400 | The error will be populated in the result | [CalcResult](#calcresult) |

### /roll

#### GET
##### Summary:

Roll Dice

##### Description:

Given a roll definition this will delegate the roll to Troll and return the
results structured as JSON.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| d | query | The Troll roll definition. This can passed as the query parameter 'd' or in the request body. | No | string |
| d | body | The Troll roll definition. This can passed as the query parameter 'd' or in the request body. | No | string |
| n | query | The number of times to repeat the roll | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | The results from rolling the dice | [RollsResult](#rollsresult) |
| 400 | The error will be populated in the result | [RollsResult](#rollsresult) |

### Models


#### CalcResult

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| Average | [Probability](#probability) |  | No |
| Cumulative | string |  | No |
| Definition | string |  | No |
| Error | string |  | No |
| MeanDeviation | [Probability](#probability) |  | No |
| ProbabilitiesCum | [Probabilities](#probabilities) |  | No |
| ProbabilitiesEq | [Probabilities](#probabilities) |  | No |
| Runtime | long |  | No |
| Spread | [Probability](#probability) |  | No |

#### Probabilities

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| Probabilities | object |  |  |

#### Probability

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| Probability | number |  |  |

#### Roll

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| Roll | array |  |  |

#### Rolls

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| Rolls | array |  |  |

#### RollsResult

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| Definition | string |  | No |
| Error | string |  | No |
| NumTimes | long |  | No |
| Rolls | [Rolls](#rolls) |  | No |
| Runtime | long |  | No |