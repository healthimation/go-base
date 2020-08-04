[Back to table of contents...](../README.md)

# Create goal

This endpoint will create a goal for a user.


## Permissions

This endpoint is accessible by the user and admins. 

Goal start and expiration dates must be in the future unless you are an admin, in which case you may create a goal at any time, but the expiration date must come after the start date.

## Request

`POST goal/v1/user/:user_id/goal`

### Parameters

| Name | Type | Required | Description | Example |
|:----:|:----:|---------:|-------------|---------|
| Authorization | Header | Yes | A JWT bearer token is required. | `Bearer eyJhbGciOiJIUzI1N....` |


### Body

| Name | Type | Required | Description | Example |
|:----:|:----:|---------:|-------------|---------|
| type | string | Yes | The type of goal to create, must be one of `personal` or `point`.  New types can be added on request. | `personal` |
| category | string | Yes | A description of the goal sub-type, currently they only valid value is `weekly`, denoting this is a weekly goal. | `personal` |
| value | int | Yes | The point value of the goal to create. | `250` |
| name | string | No | The name of the goal, used only for client display. | `my goal` |
| description | string | No | The description of the goal, used only for client display | `do something awesome` |
| starts_at | string | Yes | When the goal period starts, in RFC 3339 time format.  `starts_at` may be in the future but must be before `expires_at`. | `2020-08-06T00:00:00Z` |
| expires_at | string | Yes | When the goal expires in RFC 3339 time format.  Expiration means the goal can no longer be marked as completed or skipped. | 
`2020-08-07T00:00:00Z` |


```json
{
    "type": "personal",
    "category": "weekly",
    "value": 250,
    "name": "my goal",
    "description": "do something awesome",
    "starts_at": "2020-08-06T00:00:00Z",
    "expires_at": "2020-08-07T00:00:00Z"
}
```

### Goal Types

#### Personal
Personal goals are an arbitrary goal set by the user with a name and a description and a point value.  Personal goals have the following restrictions:

* A user may have max 3 personal goals in a given (ISO calendar week)[http://myweb.ecu.edu/mccartyr/isowdcal.html].  Only the `starts_at` value is checked to determine the week, so it must fall within the desired week.
* If the goal category is `weekly` the `starts_at` and `expires_at` value must fall within the same ISO weeek.
* The point value must be one of 50, 100, or 150
* They must have a name and a description.  The service will not validate the contents of the name and description, just that the length is > 0.

#### Point
Point goals are a point value that a user wants to reach before the expiration date.  Right now the service does not check if the points were actually reached and relies on the client to check if it was.  In effect this makes point goals in the service more of a log of what the user wants to achieve.  Point goals have the following restrictions:

* You may only have one point goal with the category of `weekly` for a given a given (ISO calendar week)[http://myweb.ecu.edu/mccartyr/isowdcal.html].  Only the `starts_at` value is checked to determine the week, so it must fall within the desired week.
* If the goal category is `weekly` the `starts_at` and `expires_at` value must fall within the same ISO weeek.
* Point goal values must be one of 100, 250, 500.

### Response

There is no response body for a successful request.

| Status Code | Status | Error Code | Description |
|:-----------:|:------:|:----------:|-------------|
| 204 | Creattetd | | | 
| 400 | Bad Request | INVALID_JSON | The json body could not be parsed |
| 400 | Bad Request | INVALID_USER_ID | The user id passed was either blank or invalid |
| 400 | Bad Request | INVALID_POINTS | Point value note allowed for provided goal type and category |
| 400 | Bad Request | INVALID_GOAL | One or more of the goal options is invalid, details will be provided in the response body. |
| 401 | Unauthorized | NOT_AUTHORIZED | There was a problem with your JWT |
| 403 | Forbidden | NOT_AUTHORIZED | You are not allowed to access this endpoint |
| 409 | Conflict | MAX_GOALS_REACHED | Maximum number personal OR point goals reached for the week of of `starts_at` |
| 500 | Internal Server Error | ERROR_SERVICE | There is a problem with the service |
