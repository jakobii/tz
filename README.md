# tz
Timezone Conversion Utility

```sh
go install github.com/jakobii/tz@latest
```


## Examples

Convert unix timestamp to local time.

```sh
echo 1257894000000 | tz
```

Convert UTC (or any other rfc3339 timestamp) to local time.

```sh
echo '2009-11-10T23:00:00Z' | tz 
```