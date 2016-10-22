package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	cec "gopkg.in/robbiet480/cec.v2"
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
	ResetVolume bool   `short:"z" long:"do-not-zero-volume" description:"Whether to reset the volume to 0 at startup"`
	StartVolume int    `short:"v" long:"initial-volume" description:"Provide an initial volume level" default:"0"`
	MaxVolume   int    `short:"c" long:"max-volume" description:"The maximum supported volume" default:"100"`
}

type Options struct {
	HTTP  HTTPOptions  `group:"HTTP Server Options"`
	CEC   CECOptions   `group:"CEC Options"`
	Audio AudioOptions `group:"Audio Options"`
}

var options Options
var parser = flags.NewParser(&options, flags.Default)

var cec_conn cec.Connection

var volume_level = options.Audio.StartVolume
var input_number int
var is_muted = false

var received_events = []interface{}{}
var message_events = []cec.LogMessage{}
var key_press_events = []cec.KeyPress{}
var command_events = []cec.Command{}
var alert_events = []cec.Alert{}
var menu_state_events = []cec.MenuState{}
var source_activated_events = []cec.SourceActivated{}

func CheckForDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input string
		if c.Params.ByName("device") != "" {
			input = c.Params.ByName("device")
		} else {
			input = options.Audio.AudioDevice
		}
		if addr := cec.GetLogicalAddressByName(input); addr == -1 {
			c.AbortWithError(404, errors.New("That device ("+input+") does not exist!"))
		} else {
			c.Set("CECAddress", addr)
			c.Next()
		}
	}
}

func main() {

	if _, err := parser.Parse(); err != nil {
		log.Fatalln("Parser error", err)
	}

	conn, err := cec.Open(options.CEC.Adapter, options.CEC.Name, options.CEC.Type)
	if err != nil {
		log.Fatalln("Error opening CEC connection", err)
	}

	if conn.PollDevice(cec.GetLogicalAddressByName(options.Audio.AudioDevice)) != true {
		var word = "a"
		if options.Audio.AudioDevice == "Audio" {
			word = "an"
		}
		log.Println("You said you had " + word + " " + options.Audio.AudioDevice + " device but one cant be found!")
		// os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(ginrus.Ginrus(log.StandardLogger(), time.RFC3339, true))
	r.Use(gin.Recovery())
	r.GET("/config", config)
	r.GET("/info", info)
	r.GET("/logs", all_logs)
	r.GET("/logs/:type", logs_for_type)
	r.GET("/input", input_status)
	r.PUT("/input/:number", input_change)
	r.GET("/power/:device", CheckForDevice(), power_status)
	r.PUT("/power/:device", CheckForDevice(), power_on)
	r.DELETE("/power/:device", CheckForDevice(), power_off)
	r.GET("/volume", CheckForDevice(), vol_status)
	r.PUT("/volume/up", CheckForDevice(), vol_up)
	r.PUT("/volume/down", CheckForDevice(), vol_down)
	r.PUT("/volume/mute", CheckForDevice(), vol_mute)
	r.GET("/volume/mute", CheckForDevice(), vol_mute_status)
	r.PUT("/volume/reset", vol_reset)
	r.PUT("/volume/step/:direction/:steps", CheckForDevice(), vol_step)
	r.PUT("/volume/set/:level", CheckForDevice(), vol_set)
	r.PUT("/volume/force/:level", vol_force_set)
	r.PUT("/key/:device/:key", CheckForDevice(), key)
	r.PUT("/multikey/:device/:key/:delay/:key2", CheckForDevice(), multi_key)
	r.PUT("/channel/:device/:channel", CheckForDevice(), change_channel)
	r.POST("/transmit", transmit)

	if options.Audio.ResetVolume != true {
		// Let's reset the volume level to 0
		log.Println("Resetting volume to 0")
		addr := cec.GetLogicalAddressByName(options.Audio.AudioDevice)
		for i := 0; i < options.Audio.MaxVolume; i++ {
			log.Println("Sending VolumeDown")
			conn.Key(addr, "VolumeDown")
			volume_level = 0
		}
		log.Println("Volume has been set to 0")
	} else {
		log.Println("Not resetting volume to 0, assuming it is already at 0")
	}

	log.Println("Getting the current active input")

	for address, active := range conn.GetActiveDevices() {
		if (active) && (conn.IsActiveSource(address)) {
			input_str := strings.Split(conn.GetDevicePhysicalAddress(address), ".")[0]
			input_atoi, _ := strconv.Atoi(input_str)
			input_number = int(input_atoi)
		}
	}

	cec_conn = *conn

	go ReceiveCallbackEvents()

	hostAndPort := net.JoinHostPort(options.HTTP.Host, options.HTTP.Port)

	log.Infof("cec-web is live and on air at %s", hostAndPort)

	r.Run(hostAndPort)
}

func ReceiveCallbackEvents() {
	for {
		event := <-cec.CallbackEvents
		received_events = append(received_events, event)
		switch event := event.(type) {
		case cec.LogMessage:
			// log.Infoln("Received LogMessage", event.Message)
			message_events = append(message_events, event)
		case cec.KeyPress:
			key_press_events = append(key_press_events, event)
			// log.Infoln("Received KeyPress", event)
		case cec.Command:
			command_events = append(command_events, event)
			// log.Infoln("Received Command", event)
		case cec.Alert:
			alert_events = append(alert_events, event)
			// log.Infoln("Received Alert", event)
		case cec.MenuState:
			menu_state_events = append(menu_state_events, event)
			// log.Infoln("Received MenuState", event)
		case cec.SourceActivated:
			source_activated_events = append(source_activated_events, event)
			// log.Infoln("Received SourceActivated", event)
		default:
			// log.Infoln("Received unknown callback event!", event)
		}
	}
}

func config(c *gin.Context) {
	c.JSON(200, options)
}

func info(c *gin.Context) {
	list := cec_conn.List()
	if list != nil {
		c.JSON(200, list)
	} else {
		c.AbortWithError(500, errors.New("Unable to get info about connected CEC devices"))
	}
}

func all_logs(c *gin.Context) {
	c.JSON(200, received_events)
}

func logs_for_type(c *gin.Context) {
	logType := c.Params.ByName("type")
	switch logType {
	case "message_events":
		c.JSON(200, message_events)
	case "key_press_events":
		c.JSON(200, key_press_events)
	case "command_events":
		c.JSON(200, command_events)
	case "alert_events":
		c.JSON(200, alert_events)
	case "menu_state_events":
		c.JSON(200, menu_state_events)
	case "source_activated_events":
		c.JSON(200, source_activated_events)
	default:
		c.AbortWithError(404, errors.New("No logs of given type ("+logType+") exist!"))
	}
}

func power_on(c *gin.Context) {
	addr := c.MustGet("CECAddress").(int)

	if err := cec_conn.PowerOn(addr); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.String(200, "on")
	}
}

func power_off(c *gin.Context) {
	addr := c.MustGet("CECAddress").(int)

	if err := cec_conn.Standby(addr); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.String(200, "off")
	}
}

func power_status(c *gin.Context) {
	addr := c.MustGet("CECAddress").(int)

	status := cec_conn.GetDevicePowerStatus(addr)
	if status == "on" {
		c.String(200, "on")
	} else if status == "standby" {
		c.String(200, "off")
	} else {
		c.AbortWithError(500, errors.New("invalid power state"))
	}
}

func input_status(c *gin.Context) {
	c.String(200, "INPUT HDMI "+strconv.Itoa(input_number))
}

func input_change(c *gin.Context) {
	input := c.Params.ByName("number")
	if resp := cec_conn.Transmit(fmt.Sprintf("3f:82:%d0:00", input)); resp != nil {
		c.AbortWithError(500, resp)
	} else {
		input_atoi, _ := strconv.Atoi(input)
		input_number = int(input_atoi)
		c.String(200, "INPUT HDMI "+input)
	}
}

func change_channel(c *gin.Context) {
	addr := c.MustGet("CECAddress").(int)
	channel := c.Params.ByName("channel")

	for _, number := range channel {
		if resp := cec_conn.Key(addr, fmt.Sprintf("0x2%s", string(number))); resp != nil {
			c.AbortWithError(500, resp)
			break
		}
	}

	c.String(200, channel)
}

func vol_step(c *gin.Context) {
	steps_str := c.Params.ByName("steps")
	steps_atoi, _ := strconv.Atoi(steps_str)
	steps := int(steps_atoi)
	direction := c.Params.ByName("direction")
	addr := c.MustGet("CECAddress").(int)

	if direction != "up" && direction != "down" {
		c.AbortWithError(400, errors.New("Invalid direction. Valid directions are up or down."))
	}

	for i := 0; i < steps; i++ {
		if direction == "up" {
			if resp := cec_conn.Key(addr, "VolumeUp"); resp == nil {
				volume_level = volume_level + steps
			} else {
				c.AbortWithError(500, resp)
				break
			}
		} else if direction == "down" {
			if resp := cec_conn.Key(addr, "VolumeDown"); resp == nil {
				volume_level = volume_level - steps
			} else {
				c.AbortWithError(500, resp)
				break
			}
		}
	}

}

func vol_set(c *gin.Context) {
	level_str := c.Params.ByName("level")
	level_atoi, _ := strconv.Atoi(level_str)
	wanted_level := int(level_atoi)

	log.Println("Wanted_level is " + strconv.Itoa(wanted_level) + " and volume_level is " + strconv.Itoa(volume_level))

	addr := c.MustGet("CECAddress").(int)

	if wanted_level > options.Audio.MaxVolume {
		c.AbortWithError(400, errors.New("The maximum volume level is "+strconv.Itoa(options.Audio.MaxVolume)))
	} else if wanted_level > volume_level { // Requested level is greater then current volume level
		log.Println("FIRST")
		var final_level = wanted_level - volume_level
		log.Println("Final_level is " + strconv.Itoa(final_level))
		for i := 0; i < final_level; i++ {
			if resp := cec_conn.Key(addr, "VolumeUp"); resp != nil {
				c.AbortWithError(500, resp)
				break
			}
		}
	} else if wanted_level < volume_level { // Requested level is less than current volume level
		log.Println("SECOND")
		var final_level = volume_level - wanted_level
		log.Println("Final_level is " + strconv.Itoa(final_level))
		for i := 0; i < final_level; i++ {
			if resp := cec_conn.Key(addr, "VolumeDown"); resp != nil {
				c.AbortWithError(500, resp)
				break
			}
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
	if volume_level == options.Audio.MaxVolume {
		c.AbortWithError(400, errors.New("Volume already at maximum"))
	} else {
		addr := c.MustGet("CECAddress").(int)
		if resp := cec_conn.Key(addr, "VolumeUp"); resp != nil {
			c.AbortWithError(500, resp)
		}
		volume_level = volume_level + 1
		c.String(204, "")
	}
}

func vol_down(c *gin.Context) {
	if volume_level == 0 {
		c.AbortWithError(400, errors.New("Volume is already at minimum"))
	} else {
		addr := c.MustGet("CECAddress").(int)
		if resp := cec_conn.Key(addr, "VolumeDown"); resp != nil {
			c.AbortWithError(500, resp)
		}
		volume_level = volume_level - 1
		c.String(204, "")
	}
}

func vol_mute(c *gin.Context) {
	cec_conn.Mute()
	is_muted = true
	c.String(204, "")
}

func vol_mute_status(c *gin.Context) {
	c.String(200, strconv.FormatBool(is_muted))
}

func vol_reset(c *gin.Context) {
	addr := c.MustGet("CECAddress").(int)
	for i := 0; i < options.Audio.MaxVolume; i++ {
		log.Println("Sending VolumeDown")
		if resp := cec_conn.Key(addr, "VolumeDown"); resp != nil {
			c.AbortWithError(500, resp)
			break
		}
	}
	volume_level = 0
	c.String(200, strconv.Itoa(volume_level))
}

func key(c *gin.Context) {
	addr := c.MustGet("CECAddress").(int)
	key := c.Params.ByName("key")

	if resp := cec_conn.Key(addr, key); resp != nil {
		c.AbortWithError(500, resp)
	}
	c.String(204, "")
}

func multi_key(c *gin.Context) {
	addr := c.MustGet("CECAddress").(int)
	key := c.Params.ByName("key")
	key_two := c.Params.ByName("key2")
	delay, _ := strconv.Atoi(c.Params.ByName("delay"))

	if resp := cec_conn.Key(addr, key); resp != nil {
		c.AbortWithError(500, resp)
	}
	time.Sleep(time.Duration(delay) * time.Millisecond)
	if resp := cec_conn.Key(addr, key_two); resp != nil {
		c.AbortWithError(500, resp)
	}
	c.String(204, "")
}

func transmit(c *gin.Context) {
	var commands []string
	c.Bind(&commands)

	for _, val := range commands {
		if resp := cec_conn.Transmit(val); resp != nil {
			c.AbortWithError(500, resp)
		}
	}
	c.String(204, "")
}
