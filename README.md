cec-web
=======

A REST micro webservice to control devices via the CEC bus in HDMI. Specifically written for use with Roomie Remote via [CEC-Roomie](http://github.com/robbiet480/CEC-Roomie)

Written in Go with some help from [Gin](http://gin-gonic.github.io/gin/), [Go-Flags](https://github.com/jessevdk/go-flags) and [cec.go](https://github.com/robbiet480/cec).

Based on [chbmuc's](http://github.com/chbmuc) [cec-web](https://github.com/chbmuc/cec-web) and [cec.go](https://github.com/chbmuc/cec)

Usage
=====

    Usage:
      cec-web [OPTIONS]
    
    Application Options:
      -i, --ip=      ip to listen on (127.0.0.1)
      -p, --port=    tcp port to listen on (8080)
      -a, --adapter= cec adapter to connect to [RPI, usb, ...]
      -n, --name=    OSD name to announce on the cec bus (REST Gateway)
      -t, --type=    The device type to announce [tv, recording, reserved, tuner, playback, audio] (tuner)


JSON API
========

The app provides the following JSON based RESTful API:

## Scan CEC bus

* ``GET /info`` - Information about all the connected devices on the CEC bus

#### Resonse

    HTTP/1.1 200 OK

```json
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

## Get the current source label

* ``GET /source`` - Get the current source label, formatted for Roomie. If none, will return nothing (with a 200 status)

#### Response

    HTTP/1.1 200 OK
`INPUT HDMI 1`

## Power

* ``GET /power/:device`` - Request device power status
* ``PUT /power/:device`` - Power on device
* ``DELETE /power/:device`` - Put device in standby

``:device`` is the name of the device on the CEC bus (see ``GET /info``)

#### Responses

is powered on (PUT/GET)

    HTTP/1.1 200
`on`

is in standby/no power (GET/DELETE);

    HTTP/1.1 200
`off`

## Volume (not supported by all devices)

* ``GET /volume`` - Get the current volume
* ``PUT /volume/up`` - Increase volume
* ``PUT /volume/step/:direction/:steps`` - Move volume in direction by X steps
* ``PUT /volume/set/:level`` - Set volume to X level
* ``PUT /volume/down`` - Reduce volume
* ``PUT /volume/mute`` - Mute/unmute audio

> ``:direction`` is the direction to step the volume. Valid options are `up` or `down`
> ``:steps`` is number of steps to change the volume by
> ``:level`` is exact volume level to set

#### Responses

Volume up, Volume down, Volume mute (PUT)

    HTTP/1.1 204 No Content

Volume status (GET)
    HTTP/1.1 200
`10`

Set volume to specific level, step volume up/down (PUT)
    HTTP/1.1 200
`10`

## Remote control

* ``PUT /key/:device/:key`` - Send key press command followed by key release

> ``:device`` is the name of the device on the CEC bus (see ``GET /info``)
> ``:key`` is the name (e.g. ``down``) or the keycode in hex (e.g. ``0x00``) of a remote key

## Change the channel

* ``PUT /key/:device/:channel`` - Change the channel. Just a conveinence function instead of pressing individual buttons

> ``:device`` is the name of the device on the CEC bus (see ``GET /info``)
> ``:channel`` is the channel number (e.g. ``123``)

#### Response

    HTTP/1.1 200 OK
`123`

## Raw CEC commands

* ``POST /transmit`` - Send a list of CEC commands over the bus

data example:
```json
[
  "40:04",
  "40:64:00:48:65:6C:6C:6F:20:77:6F:72:6C:64"
]
```

Hint: Use [cec-o-matic](http://www.cec-o-matic.com/) to generate commands.
