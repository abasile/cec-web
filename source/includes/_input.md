# Input

## Get the current input

```shell
curl "http://192.168.1.2:8080/input"
```

```xml
<key>.INPUT STATUS</key>
<string>GET /input</string>
```

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8

INPUT HDMI 1
```

This endpoint returns the current input label, formatted for Roomie. If none, will return nothing (with a 200 status)

### HTTP Request

`GET http://192.168.1.2:8080/config`

## Change the input

```shell
curl -X PUT "http://192.168.1.2:8080/input/1"
```

```xml
<key>.INPUT SET</key>
<string>PUT /input/1</string>
```

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8

INPUT HDMI 1
```

This endpoint will change the input and return the new input label, formatted for Roomie.

### HTTP Request

`GET http://192.168.1.2:8080/input/1`