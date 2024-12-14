# Sparrow

## Jupiter

Compose

```bash
docker compose -f dockers/docker-compose.yaml --profile jupiter up -d
```

Markets

```bash
jq 'map(select(.owner == "srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX"))' market-cache-all.json > market-cache.json
```

## Metrics

Compose

```bash
docker compose -f dockers/docker-compose.yaml --profile tools up -d
```

- [Prometheus](http://localhost:9090)
- [Grafana](http://localhost:3000)

## App
