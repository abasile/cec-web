# Remote Control

## Send key code to CEC bus

```shell
curl -X PUT "http://192.168.1.2:8080/key/TV/PreviousChannel"
```

```xml
<key>PREVIOUS CHANNEL</key>
<string>PUT /key/TV/PreviousChannel</string>
```

```http
HTTP/1.1 204 No Content
Content-Type: text/plain; charset=utf-8
```

This endpoint will send the given key code to the specified device and then release the key

### HTTP Request

`PUT http://192.168.1.2:8080/key/<device>/<key>`

### URL Parameters

| Parameter | Description                                                  |
|-----------|--------------------------------------------------------------|
| device    | A friendly name for a device on the CEC bus                  |
| key       | A key name (e.g. `down`) or the keycode in hex (e.g. `0x00`) |

## Send multiple keys to CEC bus

```shell
curl -X PUT "http://192.168.1.2:8080/multikey/TV/VolumeUp/0/VolumeDown"
```

```xml
<key>MUTE OFF</key>
<string>PUT /multikey/TV/VolumeUp/0/VolumeDown</string>
```

```http
HTTP/1.1 204 No Content
Content-Type: text/plain; charset=utf-8
```

This endpoint will send the given key codes to the specified device with a delay between the first and second code

### HTTP Request

`PUT http://192.168.1.2:8080/mutlikey/<device>/<key1>/<delay>/<key2>`

### URL Parameters

| Parameter | Description                                                  |
|-----------|--------------------------------------------------------------|
| device    | A friendly name for a device on the CEC bus                  |
| key1      | A key name (e.g. `down`) or the keycode in hex (e.g. `0x00`) |
| delay     | The delay between the two key presses, in milliseconds       |
| key2      | A key name (e.g. `down`) or the keycode in hex (e.g. `0x00`) |

## Change the channel

```shell
curl -X PUT "http://192.168.1.2:8080/channel/TV/53"
```

```xml
<key>.CHANNEL SET</key>
<string>PUT /channel/TV/53</string>
```

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8

53
```

This endpoint will split the given channel number and send each digit as an individual key press

### HTTP Request

`PUT http://192.168.1.2:8080/channel/<device>/<channel>`

### URL Parameters

| Parameter | Description                                  |
|-----------|----------------------------------------------|
| device    | A friendly name for a device on the CEC bus  |
| channel   | A channel number                             |

## Transmit raw CEC commands

```shell
curl -X POST -d '["3f:82:10:00"]' -A "Content-type: application/json" http://192.168.1.2:8080/transmit
```

```xml
<key>INPUT HDMI 1</key>
<string>POST /transmit
  Content-type: application/json

  ["3f:82:10:00"]
</string>
```

```http
HTTP/1.1 204 No Content
Content-Type: text/plain; charset=utf-8
```

This endpoint will transmit the given commands to the CEC bus directly. Commands must be in JSON array format. 

### HTTP Request

`POST http://192.168.1.2:8080/transmit`