## <a name="resource-rate_limit">RateLimit</a>

Stability: `prototype`

Get your current rate limit status

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **rate:limit** | *integer* |  | `42` |
| **rate:remaining** | *integer* |  | `42` |
| **rate:reset** | *integer* |  | `42` |
| **resources:core:limit** | *integer* |  | `42` |
| **resources:core:remaining** | *integer* |  | `42` |
| **resources:core:reset** | *integer* |  | `42` |
| **resources:search:limit** | *integer* |  | `42` |
| **resources:search:remaining** | *integer* |  | `42` |
| **resources:search:reset** | *integer* |  | `42` |

### RateLimit Info



```
GET /rate_limit
```


#### Curl Example

```bash
$ curl -n https://api.github.com/rate_limit
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "resources": {
    "core": {
      "limit": 42,
      "remaining": 42,
      "reset": 42
    },
    "search": {
      "limit": 42,
      "remaining": 42,
      "reset": 42
    }
  },
  "rate": {
    "limit": 42,
    "remaining": 42,
    "reset": 42
  }
}
```


