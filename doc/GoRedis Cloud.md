GoRedis Cloud
======

1 instance
shards 0-99

2 instances
shards 0-49,50-99

4 instances
shards 0-24,25-49,50-74,75-99

*HOWTO:*

	SHARD SEL 100422
	SET user:100422:name latermoon
	SHARD SEL 300000
	SET user:300000:name peter

*Server*

	goredis-proxy: handle all connections
	goredis-server: goredis service
	goredis-daemon: start goredis-server process