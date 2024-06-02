package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type config struct {
	Zones []struct {
		Name string
		Id   string
	}
	DefaultSink string `yaml:"default_sink"`
	Sinks       []string
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
func createSink(sink string) {
	// This whole function should be cleaned up. Probably ways to do this a thousand times better.
	adapter := fmt.Sprintf("'{ factory.name=support.null-audio-sink node.name=\"%s\" node.description=\"%s\" media.class=Audio/Sink object.linger=true audio.position=[FL FR] }'", sink, sink)
	// Need a shell for some reason...
	cmd := exec.Command("bash", "-c", fmt.Sprintf("pw-cli create-node adapter %s", adapter))
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func setDefaultSink(sink string) {
	wpctl := exec.Command("bash", "-c", fmt.Sprintf("wpctl set-default `wpctl status | grep \"\\. %s\" | cut -c10-14 | egrep -o '[0-9]*'`", sink))
	err := wpctl.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func tearDown(sink string) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("pw-cli destroy `wpctl status | grep \"\\. %s\" | cut -c10-14 | egrep -o '[0-9]*'`", sink))
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func addZoneToSink(zone string, sink string) {
	fmt.Printf("addZoneToSink zone: %s, sink: %s", zone, sink)
	cmd := exec.Command("bash", "-c", fmt.Sprintf("pw-link \"%s:monitor_FL\" %s:playback_FL", zone, sink))
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command("bash", "-c", fmt.Sprintf("pw-link \"%s:monitor_FR\" %s:playback_FR", zone, sink))
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func removeZoneFromSink(zone string, sink string) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("pw-link -d \"%s:monitor_FL\" %s:playback_FL", zone, sink))
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command("bash", "-c", fmt.Sprintf("pw-link -d \"%s:monitor_FR\" %s:playback_FR", zone, sink))
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func webserver(conf *config) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/remove/{sink}/{zone}", func(w http.ResponseWriter, r *http.Request) {
		for _, sink := range conf.Sinks {
			if sink == r.PathValue("sink") {
				for _, zone := range conf.Zones {
					if zone.Name == r.PathValue("zone") {
						removeZoneFromSink(zone.Id, sink)
					}
				}
			}
		}
	})
	mux.HandleFunc("GET /api/add/{sink}/{zone}", func(w http.ResponseWriter, r *http.Request) {
		for _, sink := range conf.Sinks {
			if sink == r.PathValue("sink") {
				for _, zone := range conf.Zones {
					if zone.Name == r.PathValue("zone") {
						addZoneToSink(zone.Id, sink)
					}
				}
			}
		}
	})
	mux.HandleFunc("GET /api/sinks", func(w http.ResponseWriter, r *http.Request) {
		s, _ := json.Marshal(conf.Sinks)
		fmt.Fprint(w, string(s))
	})
	mux.HandleFunc("GET /api/zones", func(w http.ResponseWriter, r *http.Request) {
		z, _ := json.Marshal(conf.Zones)
		fmt.Fprint(w, string(z))
	})
	http.ListenAndServe("localhost:8090", mux)
}

func main() {
	conf, err := readConf("config/soundbox.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// for _, zone := range conf.Zones {
	// 	fmt.Println(zone.Name)
	// }

	for _, sink := range conf.Sinks {
		defer tearDown(sink)
		createSink(sink)
	}
	setDefaultSink(conf.DefaultSink)

	webserver(conf)

}
