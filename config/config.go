// Package config parses command-line/environment/config file arguments
// and make available to other packages.
package config

import (
	custom_log "esalert/log"
	"github.com/Akagi201/utilgo/conflag"
	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

// Opts configs
var Opts struct {
	Conf              string `long:"conf" description:"esalert config file"`
	AlertFileDir      string `long:"alerts" short:"a" required:"true" description:"A yaml file, or directory with yaml files, containing alert definitions"`
	ElasticSearchAddr string `long:"es-addr" default:"127.0.0.1:9200" description:"Address to find an elasticsearch instance on"`
	ElasticSearchUser string `long:"es-user" default:"elastic" description:"Username for the elasticsearch"`
	ElasticSearchPass string `long:"es-pass" default:"changeme" description:"Password for the elasticsearch"`
	LuaInit           string `long:"lua-init" description:"If set the given lua script file will be executed at the initialization of every lua vm"`
	LuaVMs            int    `long:"lua-vms" default:"1" description:"How many lua vms should be used. Each vm is completely independent of the other, and requests are executed on whatever vm is available at that moment. Allows lua scripts to not all be blocked on the same os thread"`
	SlackWebhook      string `long:"slack-webhook" description:"Slack webhook url, required if using any Slack actions"`
	ForceRun          string `long:"force-run" description:"If set with the name of an alert, will immediately run that alert and exit. Useful for testing changes to alert definitions"`
	LogLevel          string `long:"log-level" default:"info" description:"Adjust the log level. Valid options are: error, warn, info, debug"`
	LogDir            string `long:"log-dir" default:"os.stdout" description:"log dir, default to stdout"`
	DingDingWebhook string `long:"dingding-webhook" description:"dingding webhook url, required if using any Dingding action"`
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func init() {
	parser := flags.NewParser(&Opts, flags.Default|flags.IgnoreUnknown)

	parser.Parse()

	if Opts.Conf != "" {
		conflag.LongHyphen = true
		conflag.BoolValue = false
		args, err := conflag.ArgsFrom(Opts.Conf)
		if err != nil {
			panic(err)
		}

		parser.ParseArgs(args)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:true,
	})
	level, err := log.ParseLevel(strings.ToLower(Opts.LogLevel))
	if err == nil {
		log.SetLevel(level)
	} else {
		log.Errorf("invalid log level: %s", Opts.LogLevel)
	}
	if Opts.LogDir != "" {
	    for {
			writer, err := custom_log.NewWriter(Opts.LogDir, log.AllLevels)
			if err != nil {
				log.Errorf(err.Error())
				break
			}
			log.AddHook(writer)
			break
		}
	}

	log.Infof("esalert opts: %+v", Opts)
}
