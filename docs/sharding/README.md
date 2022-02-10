# Instance Sharding

Simply said, with sharding you can split up receiving events across multiple instances (shards). If you want to read more about that, please read the official [Discord docs](https://discord.com/developers/docs/topics/gateway#sharding).

## Setup

First of all, you need a shared host for the MySQL Database, Minio or S3 Storage as well as for the Redis instance. All shards will connect to these instances to share a common data persistency and state. The shards itself can be basically hosted anywhere as long as you have access to the instances of the services mentioned before.

To enable sharding, you need to set your `total` shard amount in the config. As specified in the Discord documentation, you can also use different amounts across shards to switch instances or to dynamically allocate more instances.
```yml
discord:
  sharding:
    total: 5
```

Now, you can apply shard IDs either directly and statically for each instance via config or reserve IDs automatically using the shared state instance.

If you want to go full manual, just set the shard ID via the config.
```yml
discord:
  sharding:
    id: 2
    total: 5
```

If you want to let `dgrs` take care of the distribution of share IDs, just set `autoid` to `true`.
```yml
discord:
  sharding:
    autoid: true
    total: 5
```

> **Attention:** Manually set IDs will not be registered in the `dgrs` state and will be picked by shards with `autoid` enabled. It is generally not recomendet to mix both configuration variants.

Also, keep in mind that some scheduled functions are only executed on the instance with the shard ID `0`. So, you must ensure that there is **only one single** instance with shard ID `0`. Otherwise, unexpected behavior might occur. 

## Example

If you want to play around with an example, take a look at the `docker-compose.yml` in this directory. It will set up all services required by shinpuru as well as three replicas of shinpuru using `autoid` to distribute shard IDs.

Just apply the `docker-compose.yml` to a Docker swarm stack.
```
docker swarm deploy -c docs/sharding/docker-compose.yml shinpuru
```

## How does `autoid` work?

When `autoid` is enabled in the config, shinpuru will reserve a shard ID from the shared Redis state using [`dgrs`](https://github.com/zekroTJA/dgrs). This reservation lasts for one minute. Therefore, a heartbeat ticker is started which refreshes the reservation every 45 seconds. When the instance "dies", the reservation is released. Also, when shinpuru shuts down, the reservation is released as well.

IDs are reserved consecutively starting with 0. When an ID is released in between two reserved consecutive IDs, it will be used for the next reservation.

Example:
```
reserve → 0
reserve → 1
reserve → 2
release 1
reserve → 1
reserve → 3
```

This system does not take the total amount of shards in respective to be able to scale up dynamically. This also means, if you reserve more shard IDs than specified in total, shinpuru will not start.