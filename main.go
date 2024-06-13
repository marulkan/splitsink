package main

import (
	"encoding/json"
	"fmt"
	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type config struct {
	Zones []struct {
		Name string
		Id   string
	}
	DefaultSink string `yaml:"default_sink"`
	Sinks       []string
	Webserver   struct {
		Ip   string
		Port string
	}
}

type Sink struct {
	Name  string
	Zones []struct {
		Name string
		Id   string
	}
}

type page struct {
	Zones []struct {
		Name string
		Id   string
	}
	Sinks []string
}

func readConf() (*config, error) {
	cfgFile, err := xdg.ConfigFile("splitsink/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	buf, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	c := &config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", cfgFile, err)
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
		log.Println(err)
	}
}

func setDefaultSink(sink string) {
	wpctl := exec.Command("bash", "-c", fmt.Sprintf("wpctl set-default `wpctl status | grep \"\\. %s\" | cut -c10-14 | egrep -o '[0-9]*'`", sink))
	err := wpctl.Run()
	if err != nil {
		log.Println(err)
	}
}

func tearDown(sink string) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("pw-cli destroy `wpctl status | grep \"\\. %s\" | cut -c10-14 | egrep -o '[0-9]*'`", sink))
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

func addZoneToSink(zone string, sink string) {
	fmt.Printf("addZoneToSink zone: %s, sink: %s", zone, sink)
	cmd := exec.Command("bash", "-c", fmt.Sprintf("pw-link \"%s:monitor_FL\" %s:playback_FL", sink, zone))
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	cmd = exec.Command("bash", "-c", fmt.Sprintf("pw-link \"%s:monitor_FR\" %s:playback_FR", sink, zone))
	err = cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

func removeZoneFromSink(zone string, sink string) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("pw-link -d \"%s:monitor_FL\" %s:playback_FL", sink, zone))
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	cmd = exec.Command("bash", "-c", fmt.Sprintf("pw-link -d \"%s:monitor_FR\" %s:playback_FR", sink, zone))
	err = cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

func listZonesInSinks(conf *config) []Sink {
	var r []Sink
	for _, sink := range conf.Sinks {
		var zones []struct {
			Name string
			Id   string
		}
		out, err := exec.Command("bash", "-c", fmt.Sprintf("pw-link -l \"%s:monitor_FL\" | grep \"|->\" | awk '{ print $2 }' | awk -F':' '{ print $1 }'", sink)).Output()
		//out, err := exec.Command("bash", "-c", "cat out.file").Output()
		if err != nil {
			log.Println(err)
		}
		links := strings.Split(string(out), "\n")
		for _, zone := range conf.Zones {
			for _, link := range links {
				if zone.Id == link {
					zones = append(zones, struct {
						Name string
						Id   string
					}{Name: zone.Name, Id: zone.Id})
				}
			}
		}
		r = append(r, Sink{Name: sink, Zones: zones})
	}
	return r
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
	mux.HandleFunc("GET /api/list", func(w http.ResponseWriter, r *http.Request) {
		s, _ := json.Marshal(listZonesInSinks(conf))
		fmt.Fprint(w, string(s))
	})
	mux.HandleFunc("GET /api/sinks", func(w http.ResponseWriter, r *http.Request) {
		s, _ := json.Marshal(conf.Sinks)
		fmt.Fprint(w, string(s))
	})
	mux.HandleFunc("GET /api/zones", func(w http.ResponseWriter, r *http.Request) {
		z, _ := json.Marshal(conf.Zones)
		fmt.Fprint(w, string(z))
	})
	mux.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		// s, _ := json.Marshal(conf.Sinks)
		p := &page{Sinks: conf.Sinks, Zones: conf.Zones}
		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Println(err)
		}
		t.Execute(w, p)
	})
	url := fmt.Sprintf("%s:%s", conf.Webserver.Ip, conf.Webserver.Port)
	http.ListenAndServe(url, mux)
}

func main() {
	conf, err := readConf()
	if err != nil {
		log.Fatal(err)
	}

	for _, zone := range conf.Zones {
		fmt.Println(zone.Name)
	}

	for _, sink := range conf.Sinks {
		defer tearDown(sink)
		createSink(sink)
	}
	setDefaultSink(conf.DefaultSink)

	webserver(conf)
}
