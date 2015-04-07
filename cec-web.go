package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	"github.com/robbiet480/cec"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type HTTPOptions struct {
	Host string `short:"i" long:"ip" description:"IP address to listen on" default:"0.0.0.0"`
	Port string `short:"p" long:"port" description:"TCP port to listen on" default:"8080"`
}

type CECOptions struct {
	Adapter string `short:"a" long:"adapter" description:"CEC adapter to connect to" default:"RPI"`
	Name    string `short:"n" long:"name" description:"OSD name to announce on the CEC bus" default:"cec-web"`
	Type    string `short:"t" long:"type" description:"The device type to announce as" default:"tv" default:"recording" default:"reserved" default:"playback" default:"audio" default:"tuner"`
}

type AudioOptions struct {
	AudioDevice string `short:"d" long:"audio-device" description:"The audio device to use for volume control and status" default:"Audio" default:"TV"`
	ResetVolume bool   `long:"reset-volume" description:"Whether to reset the volume to 0 at startup" default:"true"`
}

type Options struct {
	HTTP  HTTPOptions  `group:"HTTP Server Options"`
	CEC   CECOptions   `group:"CEC Options"`
	Audio AudioOptions `group:"Audio Options"`
}

var options Options
var parser = flags.NewParser(&options, flags.Default)

var volume_level int
var input_number int
var is_muted = false

func main() {
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
	cec.Open(options.CEC.Adapter, options.CEC.Name, options.CEC.Type)

	r := gin.Default()
	r.GET("/info", info)
	r.GET("/input", input_status)
	r.PUT("/input/:number", input_change)
	r.GET("/power/:device", power_status)
	r.PUT("/power/:device", power_on)
	r.DELETE("/power/:device", power_off)
	r.GET("/volume", vol_status)
	r.PUT("/volume/up", vol_up)
	r.PUT("/volume/down", vol_down)
	r.PUT("/volume/mute", vol_mute)
	r.GET("/volume/mute", vol_mute_status)
	r.PUT("/volume/reset", vol_reset)
	r.PUT("/volume/step/:direction/:steps", vol_step)
	r.PUT("/volume/set/:level", vol_set)
	r.PUT("/volume/force/:level", vol_force_set)
	r.PUT("/key/:device/:key", key)
	r.PUT("/channel/:device/:channel", change_channel)
	r.POST("/transmit", transmit)

	if options.Audio.ResetVolume {
		// Let's reset the volume level to 0
		log.Println("Resetting volume to 0")
		for i := 0; i < 100; i++ {
			addr := cec.GetLogicalAddressByName(options.Audio.AudioDevice)
			log.Println("Sending VolumeDown")
			cec.Key(addr, "VolumeDown")
		}
	} else {
		log.Println("Not resetting volume to 0, assuming it is already at 0")
	}

	volume_level = 0

	log.Println("Volume has been set to 0")

	log.Println("Getting the current active input")

	for address, active := range cec.GetActiveDevices() {
		if (active) && (cec.IsActiveSource(address)) {
			input_str := strings.Split(cec.GetDevicePhysicalAddress(address), ".")[0]
			input_atoi, _ := strconv.Atoi(input_str)
			input_number = int(input_atoi)
		}
	}

	r.Run(options.HTTP.Host + ":" + options.HTTP.Port)
}

func info(c *gin.Context) {
	c.JSON(200, cec.List())
}

func power_on(c *gin.Context) {
	addr := cec.GetLogicalAddressByName(c.Params.ByName("device"))

	cec.PowerOn(addr)
	c.String(200, "on")
}

func power_off(c *gin.Context) {
	addr := cec.GetLogicalAddressByName(c.Params.ByName("device"))

	cec.Standby(addr)
	c.String(200, "off")
}

func power_status(c *gin.Context) {
	addr := cec.GetLogicalAddressByName(c.Params.ByName("device"))

	status := cec.GetDevicePowerStatus(addr)
	if status == "on" {
		c.String(200, "on")
	} else if status == "standby" {
		c.String(200, "off")
	} else {
		c.String(500, "invalid power state")
	}
}

func input_status(c *gin.Context) {
	c.String(200, "INPUT HDMI "+strconv.Itoa(input_number))
}

func input_change(c *gin.Context) {
	input := c.Params.ByName("number")
	cec.Transmit("3f:82:" + input + "0:00")
	input_atoi, _ := strconv.Atoi(input)
	input_number = int(input_atoi)
	c.String(200, "INPUT HDMI "+input)
}

func change_channel(c *gin.Context) {
	addr := cec.GetLogicalAddressByName(c.Params.ByName("device"))
	channel := c.Params.ByName("channel")

	for _, number := range channel {
		cec.Key(addr, "0x2"+string(number))
	}

	c.String(200, channel)
}

func vol_step(c *gin.Context) {
	steps_str := c.Params.ByName("steps")
	steps_atoi, _ := strconv.Atoi(steps_str)
	steps := int(steps_atoi)
	direction := c.Params.ByName("direction")

	for i := 0; i < steps; i++ {
		if direction == "up" {
			addr := cec.GetLogicalAddressByName(options.Audio.AudioDevice)
			cec.Key(addr, "VolumeUp")
			volume_level = volume_level + steps
		} else if direction == "down" {
			addr := cec.GetLogicalAddressByName(options.Audio.AudioDevice)
			cec.Key(addr, "VolumeDown")
			volume_level = volume_level - steps
		} else {
			c.String(400, "Invalid direction. Valid directions are up or down.")
		}
	}

	c.String(204, "")
}

func vol_set(c *gin.Context) {
	level_str := c.Params.ByName("level")
	level_atoi, _ := strconv.Atoi(level_str)
	wanted_level := int(level_atoi)

	log.Println("Wanted_level is " + strconv.Itoa(wanted_level) + " and volume_level is " + strconv.Itoa(volume_level))

	if wanted_level > 100 {
		c.String(400, "The maximum volume level is 100")
	} else if wanted_level > volume_level { // Requested level is greater then current volume level
		log.Println("FIRST")
		var final_level = wanted_level - volume_level
		log.Println("Final_level is " + strconv.Itoa(final_level))
		for i := 0; i < final_level; i++ {
			addr := cec.GetLogicalAddressByName(options.Audio.AudioDevice)
			cec.Key(addr, "VolumeUp")
		}
	} else if wanted_level < volume_level { // Requested level is less than current volume level
		log.Println("SECOND")
		var final_level = volume_level - wanted_level
		log.Println("Final_level is " + strconv.Itoa(final_level))
		for i := 0; i < final_level; i++ {
			addr := cec.GetLogicalAddressByName(options.Audio.AudioDevice)
			cec.Key(addr, "VolumeDown")
		}
	}

	volume_level = wanted_level

	c.String(200, strconv.Itoa(volume_level))
}

func vol_force_set(c *gin.Context) {
	level_str := c.Params.ByName("level")
	level_atoi, _ := strconv.Atoi(level_str)
	volume_level = int(level_atoi)
	c.String(204, "")
}

func vol_status(c *gin.Context) {
	if is_muted == true {
		c.String(200, "muted")
	} else {
		c.String(200, strconv.Itoa(volume_level))
	}
}

func vol_up(c *gin.Context) {
	if volume_level == 100 {
		c.String(400, "Volume already at maximum")
	} else {
		addr := cec.GetLogicalAddressByName(options.Audio.AudioDevice)
		cec.Key(addr, "VolumeUp")
		volume_level = volume_level + 1
		c.String(204, "")
	}
}

func vol_down(c *gin.Context) {
	if volume_level == 0 {
		c.String(400, "Volume is already at minimum")
	} else {
		addr := cec.GetLogicalAddressByName(options.Audio.AudioDevice)
		cec.Key(addr, "VolumeDown")
		volume_level = volume_level - 1
		c.String(204, "")
	}
}

func vol_mute(c *gin.Context) {
	cec.Mute()
	is_muted = true
	c.String(204, "")
}

func vol_mute_status(c *gin.Context) {
	c.String(200, strconv.FormatBool(is_muted))
}

func vol_reset(c *gin.Context) {
	for i := 0; i < 100; i++ {
		log.Println("Sending VolumeDown")
		addr := cec.GetLogicalAddressByName(options.Audio.AudioDevice)

		cec.Key(addr, "VolumeDown")
	}
	volume_level = 0
	c.String(200, strconv.Itoa(volume_level))
}

func key(c *gin.Context) {
	addr := cec.GetLogicalAddressByName(c.Params.ByName("device"))
	key := c.Params.ByName("key")

	cec.Key(addr, key)
	c.String(204, "")
}

func transmit(c *gin.Context) {
	var commands []string
	c.Bind(&commands)

	for _, val := range commands {
		cec.Transmit(val)
	}
	c.String(204, "")
}
