# Core

## Scan the CEC bus

```shell
curl http://192.168.1.2:8080/info
```

```xml
<key>Info</key>
<string>GET /info</string>
```

```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{
  "Playback":{
    "OSDName":"REST Gateway",
    "Vendor":"Panasonic",
    "LogicalAddress":4,
    "ActiveSource":false,
    "PowerStatus":"on",
    "PhysicalAddress":"f.f.f.f"
  },
  "TV":{
    "OSDName":"TV",
    "Vendor":"Panasonic",
    "LogicalAddress":0,
    "ActiveSource":false,
    "PowerStatus":"standby",
    "PhysicalAddress":"0.0.0.0"
  }
}
```

This endpoint returns information about all the connected devices on the CEC bus

### HTTP Request

`GET http://192.168.1.2:8080/info`

## Get the configuration

```shell
curl http://192.168.1.2:8080/config
```

```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{
    "HTTP": {
        "Host": "0.0.0.0",
        "Port": "8080"
    },
    "CEC": {
        "Adapter": "RPI",
        "Name": "cec-web",
        "Type": "tuner"
    },
    "Audio": {
        "AudioDevice": "TV",
        "ResetVolume": true,
        "StartVolume": 0,
        "MaxVolume": 100
    }
}
```

This endpoint retrieves the configuration parameters used to start cec-web.

### HTTP Request

`GET http://192.168.1.2:8080/config`