# speedy

 > Utility for gateways/routers to collect network speed stats from everyone

The utility for Linux that captures the traffic over a network interface and stores the network usage of every device in a time-series DB. It is written in Go and uses `libpcap` to grab the packets.

The data is stored using the MAC address of every client as a key (of some sort) and storing the download and upload quantity in a second. Also stores the MAC again, the last IPv4 and IPv6 used (if available). Downsampling and cleaning is up to the database implementation or up to you (or both). It is recommended to create the structure firsrt before starting to run the utility.

## Usage

Having defined a `GOROOT` for a terminal session, you can simply do:

```bash
go get -d -v github.com/influxdata/influxdb/client/v2
go get -d -v github.com/google/gopacket
go get -d -v github.com/melchor629/speedy
go run src/github.com/melchor629/speedy/main.go #args...
```

You can always build the executable (cross-compiling is not available or at least not easily):

```bash
cd src/github.com/melchor629/speedy
go install -v ./...
$GOROOT/bin/speedy #args...
```

The utility needs a time-series database on which the data will be stored. Currently, [influxdb][1] and [timescaledb][3] are the only supported, but more it will be added.

The arguments can be seen with `speedy -help`. The available database implementations can be seen with `speedy -help db`. The available network interfaces can be seen with `speedy -help device`.

## Usage with Docker

```bash
docker container run --rm -it --cap-add NET_RAW --cap-add NET_ADMIN --network host melchor9000/speedy:alpine speedy #args...
```

It is needed to use the host network mode to allow the container to list all the network interfaces, and these two Linux capabilities to have the right permissions to capture the interfaces.

## Build the images

Being in the root of the repository, these commands will build the three tags available for this repo:

```bash
docker image build -t melchor9000/speedy:latest -f docker/latest/Dockerfile .
docker image build -t melchor9000/speedy:slim -f docker/slim/Dockerfile .
docker image build -t melchor9000/speedy:alpine-f docker/alpine/Dockerfile .
```

## Example with influxdb, chronograf and docker-compose

You can see [docker/compose/influxdb.yaml][2] for an example. To run the example, simply execute from the root of the repo:

```bash
DEVICE=YOUR_NIC_NAME docker-compose -f docker/compose/influxdb.yaml up
```

The example sets a database that everyday will downsample the data into `"monthly_gc"."downsampled"` (where `"monthly_gc"` is a _retention policy_ and `"downsampled"` is a measurement). Then, will store the downsampled versions for "a month" (aka 4 weeks). A chronograf will be avaiable in http://localhost:8888 for you to play with the data an visualize it.

## Database implementations

### influxdb

The implementation stores a measure in `measures` with the data. Is it up to you to make retention policies and continues queries, as the way you want. Inside `docker/compose/iql` there's an example of a database.

### timescaledb / postgresql

The implementation stores the data into the table passed by `-db-name` option. The `-db-url` has the following format `postgres://USER:PASSWORD@ADDRESS/DATABASE[?...extraOptions]`. See [pq][4] documentation for the full specification of the URL format. The schema and the table must be created by you, but I will let you an example down:

```sql
CREATE TABLE speedy (
  time        TIMESTAMPTZ       NOT NULL, /* This one must always be there, with that name */
  mac         MACADDR           NOT NULL,
  download    BIGINT            NOT NULL,
  upload      BIGINT            NOT NULL,
  ipv4        INET              NULL,
  ipv6        INET              NULL
);

SELECT create_hypertable('speedy', 'time');

CREATE INDEX ON speedy (mac, time DESC);
```


  [1]: https://influxdata.com
  [2]: https://github.com/melchor629/speedy/blob/master/docker/compose/influxdb.yaml
  [3]: https://timescaledb.com
  [4]: https://godoc.org/github.com/lib/pq
