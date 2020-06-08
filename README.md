Trollr: A HTTP/JSON API for Troll
========================================================================================================================

```
     _____         _ _          .-------.    ______
    |_   _|       | | |        /   o   /|   /\     \
      | |_ __ ___ | | |_ __   /_______/o|  /o \  o  \
      | | '__/ _ \| | | '__|  | o     | | /   o\_____\
      | | | | (_) | | | |     |   o   |o/ \o   /o    /
      | |_|  \___/|_|_|_|     |     o |/   \ o/  o  /
      \_/ trollr.live/api     '-------'     \/____o/

```

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/bendoerr/trollr/Build,%20Test,%20&%20Lint?label=Build%2C%20Test%2C%20Lint&logo=github&style=for-the-badge)
![License: MIT](https://img.shields.io/github/license/bendoerr/trollr?color=blue&logo=open-source-initiative&logoColor=white&style=for-the-badge)
![Swagger Validator](https://img.shields.io/swagger/valid/3.0?logo=swagger&specUrl=https%3A%2F%2Fraw.githubusercontent.com%2Fbendoerr%2Ftrollr%2Fmaster%2Fstatic%2Fswagger.json&style=for-the-badge&logoColor=white)
![Go Reportcard](https://goreportcard.com/badge/github.com/bendoerr/trollr?style=for-the-badge)
![Website](https://img.shields.io/website?down_message=offline&label=trollr.live&logo=semantic%20web&style=for-the-badge&up_message=online&url=https%3A%2F%2Ftrollr.live)

Trollr is a simple wrapper around the amazing [Troll: A dice roll language and calculator][troll_homepage] created by 
[Torben Mogensen][torben_homepage]. Trollr has no affiliation with Troll or Torben Mogensen but the author of Trollr
greatly appreciates his work. The wrapper simply exposes and HTTP/JSON server that executes Troll, parses the results 
and returns it. The server has some built in rate-limiting and pooling to prevent abuse. I created this small server to
support a Discord bot that I am working on.

[troll_homepage]: http://hjemmesider.diku.dk/~torbenm/Troll/
[torben_homepage]: http://hjemmesider.diku.dk/~torbenm/

Using the Trollr API
------------------------------------------------------------------------------------------------------------------------

### Swagger

Trollr has a Swagger UI that can be used for examples and testing at [trollr.live/api/swagger][trollr_swagger]. In
addition, you can find the latest Swagger JSON here in this repository or served from the API at
[trollr.live/api/swagger.json][trollr_swagger_json].

[trollr_swagger]: https://trollr.live/api/swagger
[trollr_swagger_json]: https://trollr.live/api/swagger.json

### The API

#### Roll - POST trollr.live/api/roll

Given a roll definition this will delegate the roll to Troll and return the
results structured as JSON.

###### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| d | query | The Troll roll definition. This can passed as the query parameter 'd' or in the request body. | No | string |
| d | body | The Troll roll definition. This can passed as the query parameter 'd' or in the request body. | No | string |
| n | query | The number of times to repeat the roll | No | integer |

###### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | The results from rolling the dice | [RollsResult](#rollsresult) |
| 400 | The error will be populated in the result | [RollsResult](#rollsresult) |

###### RollsResult

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| Definition | string |  | No |
| Error | string |  | No |
| NumTimes | long |  | No |
| Rolls | [Rolls](#rolls) |  | No |
| Runtime | long |  | No |

###### Roll

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| Roll | array[string] |  |  |

###### Rolls

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| Rolls | array[[Roll](#roll)] |  |  |

##### Examples

###### Simple roll definition in the query parameter

```http request
POST /api/roll?d=2d6 HTTP/1.1
Accept: */*
Host: trollr.live



HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Notice-Message: The 'Trollr' API is a simple server that wraps the amazing 'Troll' program. This API is not associated with the author of the 'Troll' program.
Notice-Troll-Author: Torben Mogensen <torbenm@di.ku.dk>
Notice-Troll-Manual: http://hjemmesider.diku.dk/~torbenm/Troll/manual.pdf
Notice-Troll-Url: http://hjemmesider.diku.dk/~torbenm/Troll/

{
    "definition": "2d6",
    "num_times": 1,
    "rolls": [
        [
            "3",
            "4"
        ]
    ],
    "runtime": 1
}
```

###### Complex roll definition in the request body

```http request
POST /api/roll HTTP/1.1
Accept: application/json, */*
Host: trollr.live

\ Savage Worlds
N:=8;
max { sum (accumulate x:=d6 while x=6),
      sum (accumulate y:=d N while y=N) }

HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Notice-Message: The 'Trollr' API is a simple server that wraps the amazing 'Troll' program. This API is not associated with the author of the 'Troll' program.
Notice-Troll-Author: Torben Mogensen <torbenm@di.ku.dk>
Notice-Troll-Manual: http://hjemmesider.diku.dk/~torbenm/Troll/manual.pdf
Notice-Troll-Url: http://hjemmesider.diku.dk/~torbenm/Troll/

{
    "definition": "\\ Savage Worlds\r\nN:=8;\r\nmax { sum (accumulate x:=d6 while x=6),\r\n      sum (accumulate y:=d N while y=N) }",
    "num_times": 1,
    "rolls": [
        [
            "4"
        ]
    ],
    "runtime": 1014
}
```


#### Calc - POST trollr.live/api/calc

Given a roll definition this will delegate the roll to Troll and return the
probabilities structured as JSON.

###### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| d | query | The Troll roll definition. This can passed as the query parameter 'd' or in the request body. | No | string |
| d | body | The Troll roll definition. This can passed as the query parameter 'd' or in the request body. | No | string |
| c | query | What kind of cumulative probabilities you would like. One of 'ge' (default), 'gt', 'le', or 'lt'. | No | string |

###### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | The probabilities of rolling the dice | [CalcResult](#calcresult) |
| 400 | The error will be populated in the result | [CalcResult](#calcresult) |

###### CalcResult

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

###### Probabilities

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| Probabilities | map[string][Probability](#probability) |  |  |

###### Probability

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| Probability | number |  |  |


##### Example

```http request
POST /calc?d=sum+2d6 HTTP/1.1
Accept: */*
Host: localhost:6789



HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Notice-Message: The 'Trollr' API is a simple server that wraps the amazing 'Troll' program. This API is not associated with the author of the 'Troll' program.
Notice-Troll-Author: Torben Mogensen <torbenm@di.ku.dk>
Notice-Troll-Manual: http://hjemmesider.diku.dk/~torbenm/Troll/manual.pdf
Notice-Troll-Url: http://hjemmesider.diku.dk/~torbenm/Troll/

{
    "average": 7,
    "cumulative": "ge",
    "definition": "sum 2d6",
    "mean_deviation": 1.94444444444,
    "probabilities_eq": {
        "2": 2.77777777778,
        "3": 5.55555555556,
        "4": 8.33333333333,
        "5": 11.1111111111,
        "6": 13.8888888889,
        "7": 16.6666666667,
        "8": 13.8888888889,
        "9": 11.1111111111,
        "10": 8.33333333333,
        "11": 5.55555555556,
        "12": 2.77777777778
    },
    "probabilities_cum": {
        "2": 100,
        "3": 97.2222222222,
        "4": 91.6666666667,
        "5": 83.3333333333,
        "6": 72.2222222222,
        "7": 58.3333333333,
        "8": 41.6666666667,
        "9": 27.7777777778,
        "10": 16.6666666667,
        "11": 8.33333333333,
        "12": 2.77777777778
    },
    "runtime": 1012,
    "spread": 2.4152294577
}
```

## Roll Definitions

Troll is compatible with the usual `ndx` dice language but extends it much further providing support for any kind of
rolling mechanic out there. 

For the best reference of creating roll definitions check out the [Troll User Manual][troll_manual]. Also helpful is the
[Troll Quick Reference][troll_quick]. And finally you can review roll defintions that the community has submitted on the
[Troll Web Interface][troll_web].

Here is a quick reference of the Troll Roll Language.

| Syntax | Description |
| ---: | :--- |
| `dn`    or    `Dn` | roll one dn (a die labeled 1 - n) |
| `mdn`    or    `mDn` | roll m dn |
| `zn`    or    `Zn` | roll one zn (a die labeled 0 - n) |
| `mzn`    or    `mZn` | roll m zn |
| `+`, `-`, `*`, `/`, `mod` | arithmetic on single values |
| `sgn` | sign of number (as -1, 0 or 1) |
| `sum` | add up values in collection |
| `count` | count values in collection |
| `U`    or    `@` | union of collections |
| `{e1,...,en}` | union of e1,...,en |
| `min`, `max` | minimum or maximum value in collection |
| `minimal`, `maximal` | all minimum or maximum values in collection |
| `median` | the median value in a collection |
| `least n`, `largest n` | n least or n largest values in collection |
| `m # e` | m samples of e |
| `..` | range of values |
| `choose` | choose value from collection |
| `e pick n` | pick (without replacement) n values from collection e |
| `<`, `<=`, `>`, `>=` , `=`, `=/=` | filters: Keep values from 2nd argument that compare to 1st argument |
| `drop` | elements found in 1st argument and not in 2nd |
| `keep` | elements found in 1st argument that are also found in 2nd |
| `--` | multiset difference |
| `different` | remove duplicates |
| `if-then-else` | conditional. Any non-empty is considered true |
| `?p` | return 1 with probability p and {} otherwise |
| `&` | substitute for logical and |
| `!` | substitute for logical not |
| `x := e1; e2` | bind x to value of e1 in e2. |
| `foreach x in e1 do e2` | evaluate e2 for each value in e1 and union the results. |
| `repeat x := e1 while/until e2` | repeatedly evaluate e1 while or until e2 becomes true (non-empty). Return last value |
| `accumulate x := e1 while/until e2` | repeatedly evaluate e1 while or until e2 becomes true (non-empty). Return union of all values |
| `function` | define function |
| `compositional` | define compositional function |
| `call` | call function |
| `[e1,e2]` | Pair of e1,e2 |
| `%1` | First component of pair |
| `%2` | Second component of pair |
| `~` | x~v returns the value of x if x is defined and otherwise returns v |

The [Troll Quick Reference][troll_quick] also lists some other operators as well as precedence. 

#### Examples

##### Savage Worlds (Aces/Exploding and Wild Card dice)

```text
\ Savage Worlds
N:=8;

max { sum (accumulate x:=d6 while x=6),
      sum (accumulate y:=d N while y=N) }
```

##### Blades in the Dark
```text
\ Blades in the Dark (Naive)
N:=2;
if N=0 then
  min(2d6)
else
  max(N d6)
```

```text
\ Blades in the Dark (Outcome)
N := 2;

R := if N=0 then 2d6 else N d6;
P := if N=0 then min(R) else max(R);
C := if N=0 then count(minimal(R)) else count(maximal(R));

if 6=P then
  if 1<C then "6+ Critical success!"
  else "6 Full success!"
else
  if 3<P then "4-5 Partial success!"
  else "1-3 Bad outcome!"
```

##### Shadowrun 5e

For a complete set of Shadowrun 5e functions including combat, edge and opposed rolls see [this gist][shadowrun_gist].

```text
\ Example: Roll a normal skill roll with a 12 dice pool and limit of 5
\call fSimpleLimit(12,5)

\ Example: Roll the same pool but using edge to "Push the Limit" and an edge rating of 3
\call fSimpleEdgePush(12,3)

\-- Simple Pool Test without Limit

function fSimple(pool)=
  count 4 < pool#(d6)


\-- Simple Pool Test with Limit

function fSimpleLimit(pool,limit)=
  min{call fSimple(pool), limit}


\-- Simple Pool Test with Push the Limit
\--   Adds edge rating number of dice to the dice pool, ignores any limits and
\--   applies the rule of sixes.

function fSimpleEdgePush(pool,edge)=
  count 4 < (pool + edge)#(accumulate x:=d6 while x=6)


\-- Simple Pool Test with Second Chances
\--   Reroll all dice that did not hit on the first roll.

function fSimpleEdgeReroll(pool, limit)=
  firstRoll:=call fSimple(pool);
  secondRoll:=call fSimple(max{pool - firstRoll, 0});
  min{firstRoll+secondRoll, limit}
```

[troll_manual]: http://hjemmesider.diku.dk/~torbenm/Troll/manual.pdf
[troll_quick]: http://hjemmesider.diku.dk/~torbenm/Troll/quickRef.html
[troll_web]: https://topps.diku.dk/torbenm/troll.msp
[shadowrun_gist]: https://gist.github.com/bendoerr/c68f8b4d5e15f10ba0be918ed110ca9a

