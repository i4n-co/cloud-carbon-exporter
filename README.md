# Cloud Carbon Exporter

Your cloud energy draw and carbon emissions in realtime.

The exporter enables your Cloud team to adhere in [Carbon Driven Development](#) principles.

## How it works

This exporter will discover all resources running in a specified project or account and estimate the energy ⚡ (watt) used by them and calculate the associated CO₂eq emissions ☁️ based on their location.

### Estimated Watts

Each resource discovered by the exporter is embellished with additional data from specific apis or cloud monitoring. Those vitality signals are used by a calculation model to precisly estimate the current energy usage (CPU load, Storage used, requests/seconds, etc.)

### Estimated CO₂eq/second

Once the watt estimation is complete, we match the resource's location with a carbon coefficient to calculate the CO₂ equivalent (CO₂eq).

The model is based on public [data shared by cloud providers.](https://github.com/GoogleCloudPlatform/region-carbon-info).

### Demo

You can find a demo grafana instance on : [https://demo.carbondriven.dev](https://demo.carbondriven.dev/public-dashboards/04a3c6d5961c4463b91a3333d488e584)

## Install

You can download the official Docker Image on the [Github Package Registry](https://github.com/superdango/cloud-carbon-exporter/pkgs/container/cloud-carbon-exporter)

```
$ docker pull ghcr.io/superdango/cloud-carbon-exporter:latest
```

## Configuration

The Cloud Carbon Exporter can work on Google Cloud Platform and Amazon Web Service (more to come)

### Google Cloud Platform

The exporter uses GCP Application Default Credentials:

- GOOGLE_APPLICATION_CREDENTIALS environment variable
- `gcloud auth application-default` login command
- The attached service account, returned by the metadata server (inside GCP)

```
$ docker run -p 2922 ghcr.io/superdango/cloud-carbon-exporter:latest \
        -cloud.provider=gcp \
        -gcp.projectid=myproject
```

### Amazon Web Services

Configure the exported via:

- Environment Variables (AWS_SECRET_ACCESS_KEY, AWS_ACCESS_KEY_ID, AWS_SESSION_TOKEN)
- Shared Configuration
- Shared Credentials files.

```
$ docker run -p 2922 ghcr.io/superdango/cloud-carbon-exporter:latest \
        -cloud.provider=aws
```

### Deployment

Cloud Carbon Exporter can easily run on serverless platform like GCP Cloud Run or AWS Lambda.

### Usage

```
Usage of ./cloud-carbon-exporter:
  -cloud.provider string
        cloud provider type (gcp, aws)
  -demo.enabled
        return fictive demo data
  -gcp.projectid string
        gcp project to export data from
  -listen string
        addr to listen to (default "0.0.0.0:2922")
  -log.format string
        log format (text, json) (default "text")
  -log.level string
        log severity (debug, info, warn, error) (default "info")
```

## Development

    go build \
        -o exporter \
        github.com/superdango/cloud-carbon-exporter/cmd && \
        ./exporter -cloud.provider=aws -log.level=debug

## Licence

This software is provided as is, without waranty under [AGPL 3.0 licence](https://www.gnu.org/licenses/agpl-3.0.en.html)

## ⭐ Sponsor

[dangofish.com](dangofish.com) - Tools and Services for Cloud Carbon Developers.
