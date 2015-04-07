package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	"github.com/robbiet480/cec"
	"os"
	"strconv"
	"strings"
)

type Options struct {
	Host    string `short:"i" long:"ip" description:"ip to listen on" default:"127.0.0.1"`
	Port    string `short:"p" long:"port" description:"tcp port to listen on" default:"8080"`
	Adapter string `short:"a" long:"adapter" description:"cec adapter to connect to [RPI, usb, ...]"`
	Name    string `short:"n" long:"name" description:"OSD name to announce on the cec bus" default:"REST Gateway"`
	Type    string `short:"t" long:"type" description:"The device type to register as" default:"tuner"`
}

var options Options
var parser = flags.NewParser(&options, flags.Default)

var volume_level int

func main() {
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	cec.Open(options.Adapter, options.Name, options.Type)

	// Let's reset the volume level to 0
	for i := 0; i < 100; i++ {
		cec.VolumeDown()
	}
	volume_level = 0

	r := gin.Default()
	r.GET("/info", info)
	r.GET("/source", source_status)
	r.GET("/power/:device", power_status)
	r.PUT("/power/:device", power_on)
	r.DELETE("/power/:device", power_off)
	r.GET("/volume", vol_status)
	r.PUT("/volume/up", vol_up)
	r.PUT("/volume/down", vol_down)
	r.PUT("/volume/mute", vol_mute)
	r.PUT("/volume/step/:direction/:steps", vol_step)
	r.PUT("/volume/set/:level", vol_set)
	r.PUT("/key/:device/:key", key)
	r.PUT("/channel/:device/:channel", change_channel)
	r.POST("/transmit", transmit)

	r.Run(options.Host + ":" + options.Port)
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

func source_status(c *gin.Context) {
	active_devices := cec.GetActiveDevices()

	for address, active := range active_devices {
		if (active) && (cec.IsActiveSource(address)) {
			c.String(200, "INPUT HDMI "+strings.Split(cec.GetDevicePhysicalAddress(address), ".")[0])
		}
	}
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

	for i := 1; i < steps; i++ {
		if direction == "up" {
			cec.VolumeUp()
		} else if direction == "down" {
			cec.VolumeDown()
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

	if wanted_level > volume_level { // Requested level is greater then current volume level
		var final_level = volume_level - wanted_level
		for i := 1; i < final_level; i++ {
			cec.VolumeUp()
		}
	} else if wanted_level < volume_level { // Requested level is less than current volume level
		var final_level = volume_level - wanted_level
		for i := 1; i < final_level; i++ {
			cec.VolumeDown()
		}
	}

	c.String(200, volume_level)
}

func vol_status(c *gin.Context) {
	c.String(200, vol_status)
}

func transmit(c *gin.Context) {
	var commands []string
	c.Bind(&commands)

	for _, val := range commands {
		cec.Transmit(val)
	}
	c.String(204, "")
}

func vol_up(c *gin.Context) {
	cec.VolumeUp()
	c.String(204, "")
}

func vol_down(c *gin.Context) {
	cec.VolumeDown()
	c.String(204, "")
}

func vol_mute(c *gin.Context) {
	cec.Mute()
	c.String(204, "")
}

func key(c *gin.Context) {
	addr := cec.GetLogicalAddressByName(c.Params.ByName("device"))
	key := c.Params.ByName("key")

	cec.Key(addr, key)
	c.String(204, "")
}
