# Getting Started

## Requirements

* Go installed
* A device capable of running [libcec](http://libcec.pulse-eight.com/) such as a Raspberry Pi or the [Pulse-Eight USB-CEC Adapter](https://www.pulse-eight.com/p/104/usb-hdmi-cec-adapter)

## Installation
1. Clone this repository
2. Run `go build cec-web.go` in the resulting folder
3. Run the resulting binary: `./cec-web`

## Command Line Options

```
Usage:
  cec-web [OPTIONS]

HTTP Server Options:
  -i, --ip=                 IP address to listen on (0.0.0.0)
  -p, --port=               TCP port to listen on (8080)
  -r, --announce            Whether to announce the server location via Avahi/Bonjour/Zeroconf (true)

CEC Options:
  -a, --adapter=            CEC adapter to connect to (RPI)
  -n, --name=               OSD name to announce on the CEC bus (cec-web)
  -t, --type=               The device type to announce as (tv, recording, reserved, playback, audio, tuner)

Audio Options:
  -d, --audio-device=       The audio device to use for volume control and status (Audio, TV)
  -z, --do-not-zero-volume  Whether to reset the volume to 0 at startup
  -v, --initial-volume=     Provide an initial volume level (0)
  -c, --max-volume=         The maximum supported volume (100)

Help Options:
  -h, --help                Show this help message

```