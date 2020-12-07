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

type NodeExporter struct {
	new *NodeExporterWin
	// nel *NodeExporterLinux
}

func NewNodeExporter() *NodeExporter {
	return &NodeExporter{
		new: &NodeExporterWin{
			Collectors: &NodeExporterWinCollectors{
				Enabled: "cpu,cs,logical_disk,net,os,service,system",
			},

			Collector: &NodeExporterWinCollector{
				Service: &NodeExporterWinService{
					ServicesWhere: "default name",
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
	}
}

func (ne *NodeExporter) WriteToFile(configs config.CombinedServices, path string) {
	pathWin := fmt.Sprintf("%sprometheus.yml", path)
	// pathLinux := fmt.Sprintf("%sdocker-compose.yml", path)

	hostname, err := os.Hostname()
	if err != nil {
		log.Panic().Err(err)
	}

	mapServiceNode := make(map[string]bool)
	for _, config := range configs {
		mapServiceNode[fmt.Sprintf(":%s", config.NodePort)] = true
	}

	// generate node_exporter
	for addr := range mapServiceNode {
		ne.new.Collector.Service.ServicesWhere = fmt.Sprintf("Name='%s'", hostname)
		ne.new.Telemetry.Addr = addr
	}

	data, err := yaml.Marshal(ne.new)
	if err != nil {
		log.Fatal().Err(err).Msg("yaml marshal failed")
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("write prometheus.yml failed")
	}

	log.Info().Str("path", path).Msgf("write %s successful", pathWin)
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
}
