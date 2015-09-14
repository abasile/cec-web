cec-web
=======

A REST micro webservice to control devices via the CEC bus in HDMI. Specifically written for use with Roomie Remote via [CEC-Roomie](http://github.com/robbiet480/CEC-Roomie)

Written in Go with some help from [Gin](http://gin-gonic.github.io/gin/), [Go-Flags](https://github.com/jessevdk/go-flags) and [cec.go](https://github.com/robbiet480/cec).

Based on [chbmuc's](http://github.com/chbmuc) [cec-web](https://github.com/chbmuc/cec-web) and [cec.go](https://github.com/chbmuc/cec)

Usage
=====

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


JSON API
========

Docs for the JSON API can be found at [Github Pages](https://robbiet480.github.io/cec-web/)

Hint: Use [cec-o-matic](http://www.cec-o-matic.com/) to generate commands.