# Power

## Get the power status for a device

```shell
curl "http://192.168.1.2:8080/power/TV"
```

```xml
<key>.POWER STATUS</key>
<string>GET /power/TV</string>
```

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8

on
```

This endpoint returns the power status for a given device

### HTTP Request

`GET http://192.168.1.2:8080/power/<device>`

### URL Parameters

| Parameter | Description                                                  |
|-----------|--------------------------------------------------------------|
| device    | A friendly name for a device on the CEC bus                  |

## Turn a device on

```shell
curl -X PUT "http://192.168.1.2:8080/power/TV"
```

```xml
<key>POWER ON</key>
<string>PUT /power/TV</string>
```

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8

on
```

This endpoint turns the given device on

### HTTP Request

`PUT http://192.168.1.2:8080/power/<device>`

### URL Parameters

| Parameter | Description                                                  |
|-----------|--------------------------------------------------------------|
| device    | A friendly name for a device on the CEC bus                  |

## Turn a device off

```shell
curl -X DELETE "http://192.168.1.2:8080/power/TV"
```

```xml
<key>POWER OFF</key>
<string>DELETE /power/TV</string>
```

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8

off
```

This endpoint turns the given device off

### HTTP Request

`DELETE http://192.168.1.2:8080/power/<device>`

### URL Parameters

| Parameter | Description                                                  |
|-----------|--------------------------------------------------------------|
| device    | A friendly name for a device on the CEC bus                  |