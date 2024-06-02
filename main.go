package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
)

type config struct {
	Zones []struct {
		Name string
		Id   string
	}
	Sink string `yaml:"default_sink"`
}

func readConf(filename string) (*config, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	return c, err
}

// Create a sink which we can add and remove devices from.
func setup(sink string) {
	// This whole function should be cleaned up. Probably ways to do this a thousand times better.
	adapter := fmt.Sprintf("'{ factory.name=support.null-audio-sink node.name=\"%s\" node.description=\"%s\" media.class=Audio/Sink object.linger=true audio.position=[FL FR] }'", sink, sink)
	// Need a shell for some reason...
	cmd := exec.Command("bash", "-c", fmt.Sprintf("pw-cli create-node adapter %s", adapter))
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	wpctl := exec.Command("bash", "-c", fmt.Sprintf("wpctl set-default `wpctl status | grep \"\\. %s\" | cut -c10-14 | egrep -o '[0-9]*'`", sink))
	err = wpctl.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// FIXME: should be defered to so we don't end up with multiple sinks.
func tearDown(sink string) {}

func main() {
	conf, err := readConf("config/soundbox.yaml")
	if err != nil {
		log.Fatal(err)
	}

	for _, zone := range conf.Zones {
		fmt.Println(zone.Name)
	}
	fmt.Println(conf.Sink)

	setup(conf.Sink)

}
