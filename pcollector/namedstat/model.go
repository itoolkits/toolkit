// named stats collector

package namedstat

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace     = "bind"
	resolverStats = "resolver_stats"
	cacheStats    = "cache_stats"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the Bind instance query successful?",
		nil, nil,
	)
	//bootTime = prometheus.NewDesc(
	//	prometheus.BuildFQName(namespace, "", "boot_time_seconds"),
	//	"Start time of the BIND process since unix epoch in seconds.",
	//	nil, nil,
	//)
	nameServerStatistics = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "name_server_stats_total"),
		"Name Server Statistics Counters.",
		[]string{"type"}, nil,
	)
	outgoingQueries = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "outgoing_queries_total"),
		"Outgoing Queries.",
		[]string{"view", "type"}, nil,
	)
	incomingQueries = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "incoming_queries_total"),
		"Number of incoming DNS queries.",
		[]string{"type"}, nil,
	)
	incomingRequests = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "incoming_requests_total"),
		"Number of incoming DNS requests.",
		[]string{"opcode"}, nil,
	)
	resolverMetricStatsFile = map[string]*prometheus.Desc{
		"GlueFetchV4": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, resolverStats, "ipv4_ns_total"),
			"IPv4 NS address fetches.",
			[]string{"view"}, nil,
		),
		"GlueFetchV6": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, resolverStats, "ipv6_ns_total"),
			"IPv6 NS address fetches.",
			[]string{"view"}, nil,
		),
		"EDNS0Fail": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, resolverStats, "query_edns_failures_total"),
			"EDNS(0) query failures.",
			[]string{"view"}, nil,
		),
		"Mismatch": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, resolverStats, "response_mismatch_total"),
			"Number of mismatch responses received.",
			[]string{"view"}, nil,
		),
		"Retry": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, resolverStats, "query_retries_total"),
			"Number of resolver query retries.",
			[]string{"view"}, nil,
		),
		"Truncated": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, resolverStats, "response_truncated_total"),
			"Number of truncated responses received.",
			[]string{"view"}, nil,
		),
		"QueryV4": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, resolverStats, "ipv4_queries_sent_total"),
			"IPv4 queries sent.",
			[]string{"view"}, nil,
		),
		"QueryV6": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, resolverStats, "ipv6_queries_sent_total"),
			"IPv6 queries sent.",
			[]string{"view"}, nil,
		),
		"ResponseV4": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, resolverStats, "ipv4_responses_received_total"),
			"IPv4 responses received.",
			[]string{"view"}, nil,
		),
		"ResponseV6": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, resolverStats, "ipv6_responses_received_total"),
			"IPv6 responses received.",
			[]string{"view"}, nil,
		),
		"NXDOMAIN": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, resolverStats, "nxdomain_received_total"),
			"NXDOMAIN received.",
			[]string{"view"}, nil,
		),
		"SERVFAIL": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, resolverStats, "servfail_received_total"),
			"SERVFAIL received.",
			[]string{"view"}, nil,
		),
		"QryRTTnn": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, resolverStats, "queries_with_rtt_milliseconds_histogram"),
			"Frequency table on round trip times (RTTs) of queries. Each nn specifies the corresponding frequency.",
			[]string{"view"}, nil,
		),
		"QueryTimeout": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, resolverStats, "query_timeouts_total"),
			"Query timeouts.",
			[]string{"view"}, nil,
		),
	}
	socketIO = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "socket_io_total"),
		"Socket I/O statistics counters are defined per socket types.",
		[]string{"type"}, nil,
	)
	zoneMetricStats = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "zone_maintenance_total"),
		"Zone Maintenance Statistics Counters.",
		[]string{"type"}, nil,
	)
	cacheRRsetsStats = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cacheStats, "cache_rrsets"),
		"Number of RRSets in Cache database.",
		[]string{"view", "type"}, nil,
	)
	cacheStatistics = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cacheStats, "statistics"),
		"Cache Statistics.",
		[]string{"view", "type"}, nil,
	)
	cacheMetricStatsFile = map[string]*prometheus.Desc{
		"Buckets": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cacheStats, "database_buckets"),
			"cache database hash buckets",
			[]string{"view"}, nil,
		),
		"Nodes": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cacheStats, "database_nodes"),
			"cache database nodes",
			[]string{"view"}, nil,
		),
		"UseHeapHighest": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cacheStats, "use_heap_highest"),
			"cache heap highest memory in use",
			[]string{"view"}, nil,
		),
		"UseHeapMemory": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cacheStats, "use_heap_memory"),
			"cache heap memory in use",
			[]string{"view"}, nil,
		),
		"TotalHeapMemory": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cacheStats, "heap_memory_total"),
			"cache heap memory total",
			[]string{"view"}, nil,
		),
		"Hits": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cacheStats, "hits"),
			"cache hits",
			[]string{"view"}, nil,
		),
		"QueryHits": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cacheStats, "query_hits"),
			"cache hits from query",
			[]string{"view"}, nil,
		),
		"Misses": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cacheStats, "misses"),
			"cache misses",
			[]string{"view"}, nil,
		),
		"QueryMisses": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cacheStats, "query_misses"),
			"cache misses from query",
			[]string{"view"}, nil,
		),
		"DelTTL": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cacheStats, "delete_ttl"),
			"cache records deleted due to TTL expiration",
			[]string{"view"}, nil,
		),
		"DelMem": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cacheStats, "delete_memory"),
			"cache records deleted due to memory exhaustion",
			[]string{"view"}, nil,
		),
		"UseTreeHighest": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cacheStats, "use_tree_highest"),
			"cache tree highest memory in use",
			[]string{"view"}, nil,
		),
		"UseTreeMemory": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cacheStats, "use_tree_memory"),
			"cache tree memory in use",
			[]string{"view"}, nil,
		),
		"TotalTreeMemory": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cacheStats, "tree_memory_total"),
			"cache tree memory total",
			[]string{"view"}, nil,
		),
	}
	nameServerMap = map[string]string{
		"IPv4 requests received":                       "IPv4",
		"TCP requests received":                        "ReqTCP",
		"UDP queries received":                         "ReqUDP",
		"duplicate queries received":                   "QryDuplicate",
		"queries caused recursion":                     "QryRecursion",
		"queries dropped":                              "QryDropped",
		"queries resulted in NXDOMAIN":                 "QryNXDOMAIN",
		"queries resulted in SERVFAIL":                 "QrySERVFAIL",
		"queries resulted in authoritative answer":     "QryAuthAns",
		"queries resulted in non authoritative answer": "QryNoauthAns",
		"queries resulted in nxrrset":                  "QryNxrrset",
		"queries resulted in referral answer":          "QryReferral",
		"queries resulted in successful answer":        "QrySuccess",
		"requests with EDNS received":                  "ReqEdns",
		"requests with TSIG received":                  "ReqTSIG",
		"responses sent":                               "Response",
		"responses with EDNS sent":                     "RespEDNS",
		"truncated responses sent":                     "RespTruncated",
	}
	resolverStatisticsMap = map[string]string{
		"IPv4 queries sent":            "Queryv4",
		"IPv6 queries sent":            "Queryv6",
		"IPv4 responses received":      "Responsev4",
		"IPv6 responses received":      "Responsev6",
		"NXDOMAIN received":            "NXDOMAIN",
		"SERVFAIL received":            "SERVFAIL",
		"FORMERR received":             "FORMERR",
		"Mismatch responses received":  "Mismatch",
		"Other errors received":        "OtherError",
		"EDNS(0) query failures":       "EDNS0Fail",
		"IPv4 NS address fetches":      "GlueFetchv4",
		"IPv6 NS address fetches":      "GlueFetchv6",
		"queries with RTT 10-100ms":    "QryRTTnn",
		"queries with RTT 100-500ms":   "QryRTTnn",
		"queries with RTT 500-800ms":   "QryRTTnn",
		"queries with RTT 800-1600ms":  "QryRTTnn",
		"queries with RTT < 10ms":      "QryRTTnn",
		"queries with RTT > 1600ms":    "QryRTTnn",
		"query retries":                "Retry",
		"query timeouts":               "QueryTimeout",
		"truncated responses received": "Truncated",
		"bucket size":                  "BucketSize",
	}
	socketMap = map[string]string{
		"Raw sockets opened": "Raw_Open",
		"Raw sockets active": "Raw_Active",

		"TCP/IPv4 connections accepted":    "TCPv4_Accept",
		"TCP/IPv4 sockets active":          "TCPv4_Active",
		"TCP/IPv4 sockets closed":          "TCPv4_Close",
		"TCP/IPv4 sockets opened":          "TCPv4_Open",
		"TCP/IPv4 connections established": "TCPv4_Conn",

		"TCP/IPv6 socket bind failures": "TCPv6_BindFail",
		"TCP/IPv6 sockets closed":       "TCPv6_Close",
		"TCP/IPv6 sockets opened":       "TCPv6_Open",

		"UDP/IPv4 connections established": "UDPv4_Conn",
		"UDP/IPv4 send errors":             "UDPv4_SendErr",
		"UDP/IPv4 sockets active":          "UDPv4_Active",
		"UDP/IPv4 sockets closed":          "UDPv4_Close",
		"UDP/IPv4 sockets opened":          "UDPv4_Open",
		"UDP/IPv4 socket bind failures":    "UDPv4_BindFail",
	}
	zoneMap = map[string]string{
		"IPv6 notifies sent":               "NotifyOutv6",
		"IPv6 notifies received":           "NotifyInv6",
		"IPv6 SOA queries sent":            "SOAOutv6",
		"IPv6 AXFR requested":              "AXFRReqv6",
		"IPv6 IXFR requested":              "IXFRReqv6",
		"IPv4 IXFR requested":              "IXFRReqv4",
		"IPv4 SOA queries sent":            "SOAOutv4",
		"IPv4 notifies received":           "NotifyInv4",
		"IPv4 notifies sent":               "NotifyOutv4",
		"IPv4 AXFR requested":              "AXFRReqv4",
		"notifies rejected":                "NotifyRej",
		"Incoming notifies rejected":       "NotifyRej",
		"transfer requests succeeded":      "XfrSuccess",
		"Zone transfer requests succeeded": "XfrSuccess",
		"Zone transfer requests failed":    "XfrFail",
		"transfer requests failed":         "XfrFail",
	}
	cacheStatsMap = map[string]string{
		"cache database hash buckets":                    "Buckets",
		"cache database nodes":                           "Nodes",
		"cache heap highest memory in use":               "UseHeapHighest",
		"cache heap memory in use":                       "UseHeapMemory",
		"cache heap memory total":                        "TotalHeapMemory",
		"cache hits":                                     "Hits",
		"cache hits (from query)":                        "QueryHits",
		"cache misses":                                   "Misses",
		"cache misses (from query)":                      "QueryMisses",
		"cache records deleted due to TTL expiration":    "DelTTL",
		"cache records deleted due to memory exhaustion": "DelMem",
		"cache tree highest memory in use":               "UseTreeHighest",
		"cache tree memory in use":                       "UseTreeMemory",
		"cache tree memory total":                        "TotalTreeMemory",
	}
)
