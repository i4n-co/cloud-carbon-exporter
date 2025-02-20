package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	cloudcarbonexporter "github.com/superdango/cloud-carbon-exporter"
	"github.com/superdango/cloud-carbon-exporter/internal/gcp"
	"google.golang.org/api/monitoring/v1"
)

func main() {



	for i := 0.0; i < 100; i += 10 {
		fmt.Println(int(cloudcarbonexporter.CPUWatts(i)))
	}

	ctx := context.Background()

	service, err := monitoring.NewService(ctx)
	if err != nil {
		panic(err)
	}

	type instanceInfo struct {
		duration   time.Duration
		cpuPercent float64
		cpuCores   float64
	}
	instances := make(map[string]instanceInfo)

	body, err := service.Projects.Location.Prometheus.Api.V1.QueryRange("projects/"+projectID, "global", &monitoring.QueryRangeRequest{
		Start: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		End:   time.Date(2024, 9, 31, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		Step:  "1h",
		Query: fmt.Sprintf(`compute_googleapis_com:instance_cpu_utilization{monitored_resource="gce_instance",metadata_system_region="%s"}`, region),
	}).Context(ctx).Do()
	if err != nil {
		panic(err)
	}

	queryResponse := new(gcp.PromQueryResponse)
	if err := mapstructure.Decode(body.Data, queryResponse); err != nil {
		panic(err)
	}

	for _, metric := range queryResponse.Result {
		startTS, _ := metric.ValueAt(0)
		startDate := time.Unix(int64(startTS), 0)

		endTS, _ := metric.ValueAt(len(metric.Values) - 1)
		endDate := time.Unix(int64(endTS), 0)

		avgUsage := 0.0
		i := 0
		for i = 0; i < len(metric.Values); i++ {
			_, v := metric.ValueAt(i)
			avgUsage += v
		}

		avgUsage /= float64(i + 1)

		instances[metric.Metric["instance_name"]] = instanceInfo{
			duration:   endDate.Sub(startDate),
			cpuPercent: avgUsage * 100,
			cpuCores:   1,
		}
	}

	body, err = service.Projects.Location.Prometheus.Api.V1.QueryRange("projects/"+projectID, "global", &monitoring.QueryRangeRequest{
		Start: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		End:   time.Date(2024, 9, 31, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		Step:  "5m",
		Query: fmt.Sprintf(`compute_googleapis_com:instance_cpu_reserved_cores{monitored_resource="gce_instance",metadata_system_region="%s"}`, region),
	}).Context(ctx).Do()
	if err != nil {
		panic(err)
	}

	queryResponse = new(gcp.PromQueryResponse)
	if err := mapstructure.Decode(body.Data, queryResponse); err != nil {
		panic(err)
	}

	for _, metric := range queryResponse.Result {

		_, v := metric.ValueAt(0)

		info := instances[metric.Metric["instance_name"]]
		info.cpuCores = v
		if v < 1 {
			info.cpuCores = 1
		}
		instances[metric.Metric["instance_name"]] = info
	}

	totalWh := 0.0
	for instance, info := range instances {
		wattsHour := info.cpuCores * cloudcarbonexporter.CPUWatts(info.cpuPercent) * info.duration.Hours()
		totalWh += wattsHour
		fmt.Printf("%s: cores=%.02f utilization=%0.2f%% time=%s Wh=%.0f\n", instance, info.cpuCores, info.cpuPercent, info.duration.String(), wattsHour)
	}

	fmt.Println("TOTAL kWh", totalWh/1000, "KgCO2eq=", carbonIntensity*(totalWh/1000))

	// // (W) watt/core = 1
	// (C) cores = 32
	// (S) seconds = 3600
	// (U) usage cpu = 24%
	// (L) lambda = 3
	// (X) adjuster = 2
	// W x C x S
	// --------- = (B) base 32 Wh
	//   3600

	// B +  2B * exp(U, L)
}
