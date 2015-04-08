# Volume

## Get the current volume level

```shell
curl "http://192.168.1.2:8080/volume"
```

```xml
<key>.VOLUME STATUS</key>
<string>GET /volume</string>
```

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8

10
```

This endpoint returns the current volume level.

### HTTP Request

`GET http://192.168.1.2:8080/volume`

## Set the volume level

```shell
curl -X PUT "http://192.168.1.2:8080/volume/set/10"
```

```xml
<key>.VOLUME SET</key>
<string>PUT /volume/set/10</string>
```

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8

10
```

This endpoint will set the volume to the level given and return the new volume level

### HTTP Request

`PUT http://192.168.1.2:8080/volume/set/<level>`

### URL Parameters

| Parameter | Description                              |
|-----------|------------------------------------------|
| level     | The number of steps to change the volume |


## Increment/decrement the volume level

```shell
curl -X PUT "http://192.168.1.2:8080/volume/up"
```

```xml
<key>VOLUME UP</key>
<string>PUT /volume/up</string>
```

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8

11
```

This endpoint will change the volume level by 1 in the direction given and return the new volume level

### HTTP Request

`PUT http://192.168.1.2:8080/volume/<direction>`

### URL Parameters

| Parameter | Description                        | Options  |
|-----------|------------------------------------|----------|
| direction | The direction to change the volume | up, down |

## Step the volume level

```shell
curl -X PUT "http://192.168.1.2:8080/volume/step/up/10"
```

```xml
<key>Volume +10</key>
<string>PUT /volume/step/up/10</string>
```

```
20
```

This endpoint will change the volume level by the number of steps given in the direction given and return the new volume level

### HTTP Request

`PUT http://192.168.1.2:8080/volume/step/<direction>/<steps>`

### URL Parameters

| Parameter | Description                              | Options  |
|-----------|------------------------------------------|----------|
| direction | The direction to change the volume       | up, down |
| steps     | The number of steps to change the volume |          |

## Get the mute status

```shell
curl "http://192.168.1.2:8080/volume/mute"
```

```xml
<key>.MUTE STATUS</key>
<string>GET /volume/mute</string>
```

```
true
```

This endpoint will return the mute status as a boolean

### HTTP Request

`GET http://192.168.1.2:8080/volume/mute`

## Toggle the mute status

```shell
curl -X PUT "http://192.168.1.2:8080/volume/mute"
```

```xml
<key>MUTE TOGGLE</key>
<string>PUT /volume/mute</string>
```

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8

true
```

This endpoint will toggle the mute status and return the new mute status as a boolean

### HTTP Request

`PUT http://192.168.1.2:8080/volume/mute`

## Reset Volume

```shell
curl -X PUT "http://192.168.1.2:8080/volume/reset"
```

```xml
<key>Volume Reset</key>
<string>PUT /volume/reset</string>
```

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8

11
```

This endpoint will step the volume down to the MinVolume as specified in the configuration and return the new volume level (which should be MinVolume)

### HTTP Request

`PUT http://192.168.1.2:8080/volume/reset`

## Set cec-web's volume level

```shell
curl -X PUT "http://192.168.1.2:8080/volume/force/10"
```

```xml
<key>Volume Force</key>
<string>PUT /volume/force/10</string>
```

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8

10
```

This endpoint will set the internally tracked volume level to the level given and return the new volume level. Note that this will _not_ send any messages to the CEC bus and is meant as a conveinence function to sync cec-web's knowledge of volume with reality.

### HTTP Request

`PUT http://192.168.1.2:8080/volume/force/<level>`

### URL Parameters

| Parameter | Description                              |
|-----------|------------------------------------------|
| level     | The volume level                         |