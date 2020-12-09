package export

import (
	"fmt"
	"io/ioutil"
	"main/src/config"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type NodeExporterWinTelemetry struct {
	Addr        string `yaml:"addr"`
	Path        string `yaml:"path"`
	MaxRequests int    `yaml:"max-requests"`
}

type NodeExporterWinScrape struct {
	TimeoutMargin float32 `yaml:"timeout-margin"`
}

type NodeExporterWinLog struct {
	Level string `yaml:"level"`
}

type NodeExporterWinService struct {
	ServicesWhere string `yaml:"services-where"`
}

type NodeExporterWinCollector struct {
	Service *NodeExporterWinService `yaml:"service"`
}

type NodeExporterWinCollectors struct {
	Enabled string `yaml:"enabled"`
}

type NodeExporterWin struct {
	Collectors *NodeExporterWinCollectors `yaml:"collectors"`
	Collector  *NodeExporterWinCollector  `yaml:"collector"`
	Log        *NodeExporterWinLog        `yaml:"log"`
	Scrape     *NodeExporterWinScrape     `yaml:"scrape"`
	Telemetry  *NodeExporterWinTelemetry  `yaml:"telemetry"`
}

type NodeExporterLinuxConfig struct {
	Image         string   `yaml:"image"`
	ContainerName string   `yaml:"container_name"`
	Command       string   `yaml:"command"`
	Volumes       []string `yaml:"volumes"`
	Hostname      string   `yaml:"hostname"`
	Restart       string   `yaml:"restart"`
	Ports         []string `yaml:"ports"`
}

type NodeExporterLinuxServices struct {
	NodeExporterLinuxConfig *NodeExporterLinuxConfig `yaml:"node-exporter"`
}

type NodeExporterLinux struct {
	Version  string                     `yaml:"version"`
	Services *NodeExporterLinuxServices `yaml:"services"`
}

type NodeExporter struct {
	new *NodeExporterWin
	nel *NodeExporterLinux
}

func NewNodeExporter() *NodeExporter {
	return &NodeExporter{
		new: &NodeExporterWin{
			Collectors: &NodeExporterWinCollectors{
				Enabled: "cpu,cs,logical_disk,net,os,service,system",
			},

			Collector: &NodeExporterWinCollector{
				Service: &NodeExporterWinService{
					ServicesWhere: "default windows hostname",
				},
			},

			Log: &NodeExporterWinLog{
				Level: "debug",
			},

			Scrape: &NodeExporterWinScrape{
				TimeoutMargin: 0.5,
			},

			Telemetry: &NodeExporterWinTelemetry{
				Addr:        ":9200",
				Path:        "/metrics",
				MaxRequests: 5,
			},
		},

		nel: &NodeExporterLinux{
			Version: "3",

			Services: &NodeExporterLinuxServices{
				NodeExporterLinuxConfig: &NodeExporterLinuxConfig{
					Image:         "quay.io/prometheus/node-exporter",
					ContainerName: "node-exporter",
					Command:       "--web.listen-address=:9200",
					Volumes:       []string{"/:/host:ro"},
					Hostname:      "default linux hostname",
					Restart:       "always",
					Ports:         make([]string, 0),
				},
			},
		},
	}
}

func (ne *NodeExporter) WriteToFile(configs config.CombinedServices, path string) {
	pathWin := fmt.Sprintf("%swindows_config.yml", path)
	pathLinux := fmt.Sprintf("%sdocker-compose.yml", path)

	hostname, err := os.Hostname()
	if err != nil {
		log.Panic().Err(err)
	}

	mapServiceNode := make(map[string]bool)
	for _, config := range configs {
		mapServiceNode[config.NodePort] = true
	}

	// generate node_exporter
	for addr := range mapServiceNode {
		ne.new.Collector.Service.ServicesWhere = fmt.Sprintf("Name='%s'", hostname)
		ne.new.Telemetry.Addr = fmt.Sprintf(":%s", addr)
		ne.nel.Services.NodeExporterLinuxConfig.Hostname = hostname
		ne.nel.Services.NodeExporterLinuxConfig.Ports = append(ne.nel.Services.NodeExporterLinuxConfig.Ports, fmt.Sprintf("%s:%s", addr, addr))
		ne.nel.Services.NodeExporterLinuxConfig.Command = fmt.Sprintf("--web.listen-address=:%s", addr)
	}

	// write windows_config.yml
	dataWin, errWin := yaml.Marshal(ne.new)
	if err != nil {
		log.Fatal().Err(errWin).Msg("yaml marshal failed")
	}

	errWin = ioutil.WriteFile(pathWin, dataWin, 0644)
	if errWin != nil {
		log.Fatal().Err(errWin).Msg("write windows_config.yml failed")
	}

	log.Info().Str("path", pathWin).Msgf("write %s successful", pathWin)

	// write docker-compose.yml
	dataLinux, errLinux := yaml.Marshal(ne.nel)
	if errLinux != nil {
		log.Fatal().Err(errLinux).Msg("yaml marshal failed")
	}

	errLinux = ioutil.WriteFile(pathLinux, dataLinux, 0644)
	if errLinux != nil {
		log.Fatal().Err(errLinux).Msg("write docker-compose.yml failed")
	}

	log.Info().Str("path", pathLinux).Msgf("write %s successful", pathLinux)
}

// test
func (ne *NodeExporter) UnmarshalToStruct() {
	data, err := ioutil.ReadFile("../config/node_exporter/windows_config.yml")
	if err != nil {
		log.Fatal().Err(err).Msg("read windows_config.yml failed")
	}

	err = yaml.Unmarshal(data, ne.new)
	if err != nil {
		log.Fatal().Err(err).Msg("unmarshal yaml failed")
	}

	log.Info().Interface("windows_config.yaml", ne.new).Msg("unmarshal success")

	data, err = ioutil.ReadFile("../config/node_exporter/docker-compose.yml")
	if err != nil {
		log.Fatal().Err(err).Msg("read docker-compose.yml failed")
	}

	err = yaml.Unmarshal(data, ne.nel)
	if err != nil {
		log.Fatal().Err(err).Msg("unmarshal yaml failed")
	}

	log.Info().Interface("docker_compose.yaml", ne.new).Msg("unmarshal success")
}
