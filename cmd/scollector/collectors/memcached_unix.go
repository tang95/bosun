// +build !windows

package collectors

import (
	"strconv"
	"strings"

	"bosun.org/metadata"
	"bosun.org/opentsdb"
	"bosun.org/util"
)

func init() {
	collectors = append(collectors, &IntervalCollector{F: c_memcached_stats})
}

var memcachedMeta = map[string]MetricMeta{
	"accepting_conns": MetricMeta{
		RateType: metadata.Gauge,
		Unit:     metadata.Bool,
		Desc:     "Indicates if the memcache instance is currently accepting connections.",
	},
	"auth_cmds": MetricMeta{
		RateType: metadata.Counter,
		Unit:     metadata.Operation,
		Desc:     "The number of authentication commands handled (includes both success or failure).",
	},
	"auth_errors": MetricMeta{
		RateType: metadata.Counter,
		Unit:     metadata.Operation,
		Desc:     "The number of of failed authentications.",
	},
	"bytes_read": MetricMeta{
		Metric:   "bytes",
		TagSet:   opentsdb.TagSet{"type": "read"},
		RateType: metadata.Counter,
		Unit:     metadata.Bytes,
		Desc:     "The total number of bytes read from the network.",
	},
	"bytes_written": MetricMeta{
		Metric:   "bytes",
		TagSet:   opentsdb.TagSet{"type": "write"},
		RateType: metadata.Counter,
		Unit:     metadata.Bytes,
		Desc:     "The total number of bytes written to the network.",
	},
	"cas_badval": MetricMeta{
		RateType: metadata.Counter,
		Unit:     metadata.Operation,
		Desc:     "The number of CAS requests for which a key was found, but the CAS value did not match.",
	},
	"cas_hits": MetricMeta{
		RateType: metadata.Counter,
		Unit:     metadata.Operation,
		Desc:     "The number of successful CAS requests.",
	},
	"cas_misses": MetricMeta{
		RateType: metadata.Counter,
		Unit:     metadata.Operation,
		Desc:     "The number of CAS requests against missing keys.",
	},
	"cmd_flush": MetricMeta{
		RateType: metadata.Counter,
		Unit:     metadata.Operation,
		Desc:     "The cumulative number of flush requests.",
	},
	"cmd_set": MetricMeta{
		RateType: metadata.Counter,
		Unit:     metadata.Operation,
		Desc:     "The cumulative number of storage requests.",
	},
	"cmd_get": MetricMeta{
		RateType: metadata.Counter,
		Unit:     metadata.Operation,
		Desc:     "The cumulative number of retrieval requests.",
	},
	"conn_yields": MetricMeta{
		RateType: metadata.Counter,
		Unit:     metadata.Yield,
		Desc:     "The number of times any connection yielded to another due to hitting the memcached -R limit.",
	},
	"connection_structures": MetricMeta{
		RateType: metadata.Gauge,
		Unit:     "Connection Structures",
		Desc:     "The number of connection structures allocated by the server.",
	},
	"curr_connections": MetricMeta{
		RateType: metadata.Gauge,
		Unit:     metadata.Connection,
		Desc:     "The current number of open connections.",
	},
	"curr_items": MetricMeta{
		RateType: metadata.Gauge,
		Unit:     metadata.Item,
		Desc:     "The current number of items in the cache.",
	},
	"decr_hits": MetricMeta{
		Metric:   "commands",
		TagSet:   opentsdb.TagSet{"type": "decr", "cache": "hit"},
		RateType: metadata.Counter,
		Unit:     metadata.Operation,
		Desc:     "The total number of decr command cache hits (decr decreases a stored value by 1).",
	},
	"decr_misses": MetricMeta{
		Metric:   "commands",
		TagSet:   opentsdb.TagSet{"type": "decr", "cache": "miss"},
		RateType: metadata.Counter,
		Unit:     metadata.Operation,
		Desc:     "The total number of decr command cache misses (decr decreases a stored value by 1).",
	},
	"incr_hits": MetricMeta{
		Metric:   "commands",
		TagSet:   opentsdb.TagSet{"type": "incr", "cache": "hit"},
		RateType: metadata.Counter,
		Unit:     metadata.Operation,
		Desc:     "The total number of incr command cache hits (incr increases a stored value by 1).",
	},
	"incr_misses": MetricMeta{
		Metric:   "commands",
		TagSet:   opentsdb.TagSet{"type": "incr", "cache": "miss"},
		RateType: metadata.Counter,
		Unit:     metadata.Operation,
		Desc:     "The total number of incr command cache misses (incr increases a stored value by 1).",
	},
	"get_hits": MetricMeta{
		Metric:   "commands",
		TagSet:   opentsdb.TagSet{"type": "get", "cache": "hit"},
		RateType: metadata.Counter,
		Unit:     metadata.Operation,
		Desc:     "The total number of successful get commands (cache hits) since startup.",
	},
	"get_misses": MetricMeta{
		Metric:   "commands",
		TagSet:   opentsdb.TagSet{"type": "get", "cache": "miss"},
		RateType: metadata.Counter,
		Unit:     metadata.Operation,
		Desc:     "The total number of failed get requests because nothing was cached for this key or the cached value was too old.",
	},
	"delete_hits": MetricMeta{
		Metric:   "commands",
		TagSet:   opentsdb.TagSet{"type": "delete", "cache": "hit"},
		RateType: metadata.Counter,
		Unit:     metadata.Operation,
		Desc:     "The total number of successful delete commands (cache hits) since startup.",
	},
	"delete_misses": MetricMeta{
		Metric:   "commands",
		TagSet:   opentsdb.TagSet{"type": "delete", "cache": "miss"},
		RateType: metadata.Counter,
		Unit:     metadata.Operation,
		Desc:     "The total number of delete commands for keys not existing within the cache.",
	},
	"evictions": MetricMeta{
		RateType: metadata.Counter,
		Unit:     metadata.Eviction,
		Desc:     "The Number of objects removed from the cache to free up memory for new items because Memcached reached it's maximum memory setting (limit_maxbytes).",
	},
	"limit_maxbytes": MetricMeta{
		Metric:   "cache_limit",
		RateType: metadata.Gauge,
		Unit:     metadata.Bytes,
		Desc:     "The max allowed size of the cache.",
	},
	"bytes": MetricMeta{
		Metric:   "cache_size",
		RateType: metadata.Gauge,
		Unit:     metadata.Bytes,
		Desc:     "The current size of the cache.",
	},
	"listen_disabled_num": MetricMeta{
		Metric:   "failed_connections",
		RateType: metadata.Counter,
		Unit:     metadata.Connection,
		Desc:     "The number of denied connection attempts because memcached reached it's configured connection limit.",
	},
	"threads": MetricMeta{
		RateType: metadata.Gauge,
		Unit:     metadata.Thread,
		Desc:     "The current number of threads.",
	},
	"total_connections": MetricMeta{
		RateType: metadata.Counter,
		Unit:     metadata.Connection,
		Desc:     "The total number of successful connect attempts.",
	},
	"total_items": MetricMeta{
		RateType: metadata.Counter,
		Unit:     metadata.Item,
		Desc:     "The total number of items ever stored.",
	},
}

func c_memcached_stats() (opentsdb.MultiDataPoint, error) {
	var md opentsdb.MultiDataPoint
	const metric = "memcached."
	util.ReadCommand(func(line string) error {
		f := strings.Fields(line)
		if len(f) != 2 {
			return nil
		}
		v, err := strconv.ParseFloat(f[1], 64)
		if err != nil {
			return nil
		}
		if m, ok := memcachedMeta[f[0]]; ok {
			name := f[0]
			if m.Metric != "" {
				name = m.Metric
			}
			Add(&md, metric+name, v, m.TagSet, m.RateType, m.Unit, m.Desc)
		}
		return nil
	}, "memcached-tool", "127.0.0.1:11211", "stats")
	return md, nil
}
