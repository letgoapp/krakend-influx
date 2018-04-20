## Local dev

Start the influx container

```
$ docker run -p 8086:8086 \
	  -e INFLUXDB_DB=krakend \
	  -e INFLUXDB_USER=letgo -e INFLUXDB_USER_PASSWORD=pas5w0rd \
	  -e INFLUXDB_ADMIN_USER=admin -e INFLUXDB_ADMIN_PASSWORD=supersecretpassword \
	  -it --name=influx \
	  influxdb
```

and in a new terminal, open the CLI

```
$ docker exec -it influx influx
```

## Grafana dashboard

```
$ docker run \
  -d \
  -p 3000:3000 \
  --name=grafana \
  grafana/grafana
```

and import the dashboard `grafana-dashboard.json`