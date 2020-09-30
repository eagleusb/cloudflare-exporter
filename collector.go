package main

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type CloudflareCollector struct {
	// Requests
	zoneRequestTotal          *prometheus.Desc
	zoneRequestCached         *prometheus.Desc
	zoneRequestUncached       *prometheus.Desc
	zoneRequestSSLEncrypted   *prometheus.Desc
	zoneRequestSSLUnencrypted *prometheus.Desc
	zoneRequestContentType    *prometheus.Desc
	zoneRequestCountry        *prometheus.Desc
	zoneRequestHTTPStatus     *prometheus.Desc

	// Bandwidth
	zoneBandwidthTotal          *prometheus.Desc
	zoneBandwidthCached         *prometheus.Desc
	zoneBandwidthUncached       *prometheus.Desc
	zoneBandwidthSSLEncrypted   *prometheus.Desc
	zoneBandwidthSSLUnencrypted *prometheus.Desc
	zoneBandwidthContentType    *prometheus.Desc
	zoneBandwidthCountry        *prometheus.Desc
	zoneBandwidthHTTPStatus     *prometheus.Desc
}

func newCloudflareCollector() *CloudflareCollector {
	return &CloudflareCollector{
		// Requests
		zoneRequestTotal: prometheus.NewDesc("cloudflare_zone_requests_total",
			"Number of requests for zone",
			[]string{"zone"}, nil,
		),
		zoneRequestCached: prometheus.NewDesc("cloudflare_zone_requests_cached",
			"Number of cached requests for zone",
			[]string{"zone"}, nil,
		),
		zoneRequestUncached: prometheus.NewDesc("cloudflare_zone_requests_uncached",
			"Number of uncached requests for zone",
			[]string{"zone"}, nil,
		),
		zoneRequestSSLEncrypted: prometheus.NewDesc("cloudflare_zone_requests_ssl_encrypted",
			"Number of encrypted requests for zone",
			[]string{"zone"}, nil,
		),
		zoneRequestSSLUnencrypted: prometheus.NewDesc("cloudflare_zone_requests_ssl_unencrypted",
			"Number of unencypted requests for zone",
			[]string{"zone"}, nil,
		),
		zoneRequestContentType: prometheus.NewDesc("cloudflare_zone_requests_content_type",
			"Number of request for zone per content type",
			[]string{"zone", "content_type"}, nil,
		),
		zoneRequestCountry: prometheus.NewDesc("cloudflare_zone_requests_country",
			"Number of request for zone per country",
			[]string{"zone", "country"}, nil,
		),
		zoneRequestHTTPStatus: prometheus.NewDesc("cloudflare_zone_requests_status",
			"Number of request for zone per HTTP status",
			[]string{"zone", "status"}, nil,
		),

		//Bandwidth
		zoneBandwidthTotal: prometheus.NewDesc("cloudflare_zone_bandwidth_total",
			"Total bandwidth per zone in bytes",
			[]string{"zone"}, nil,
		),
		zoneBandwidthCached: prometheus.NewDesc("cloudflare_zone_bandwidth_cached",
			"Cached bandwidth per zone in bytes",
			[]string{"zone"}, nil,
		),
		zoneBandwidthUncached: prometheus.NewDesc("cloudflare_zone_bandwidth_uncached",
			"Uncached bandwidth per zone in bytes",
			[]string{"zone"}, nil,
		),
		zoneBandwidthSSLEncrypted: prometheus.NewDesc("cloudflare_zone_bandwidth_ssl_encrypted",
			"Encrypted bandwidth per zone in bytes",
			[]string{"zone"}, nil,
		),
		zoneBandwidthSSLUnencrypted: prometheus.NewDesc("cloudflare_zone_bandwidth_ssl_unencrypted",
			"Unencrypted bandwidth per zone in bytes",
			[]string{"zone"}, nil,
		),
		zoneBandwidthContentType: prometheus.NewDesc("cloudflare_zone_bandwidth_content_type",
			"Bandwidth per zone per content type",
			[]string{"zone", "content_type"}, nil,
		),
		zoneBandwidthCountry: prometheus.NewDesc("cloudflare_zone_bandwidth_country",
			"Bandwidth per country per zone",
			[]string{"zone", "country"}, nil,
		),
		zoneBandwidthHTTPStatus: prometheus.NewDesc("cloudflare_zone_bandwidth_status",
			"Bandwidth per country per status",
			[]string{"zone", "status"}, nil,
		),
	}
}

func (collector *CloudflareCollector) Describe(ch chan<- *prometheus.Desc) {
	// Requests
	ch <- collector.zoneRequestTotal
	ch <- collector.zoneRequestCached
	ch <- collector.zoneRequestUncached
	ch <- collector.zoneRequestSSLEncrypted
	ch <- collector.zoneRequestSSLUnencrypted
	ch <- collector.zoneRequestContentType
	ch <- collector.zoneRequestCountry
	ch <- collector.zoneRequestHTTPStatus
	// Bandwidth
	ch <- collector.zoneBandwidthTotal
	ch <- collector.zoneBandwidthCached
	ch <- collector.zoneBandwidthUncached
	ch <- collector.zoneBandwidthSSLEncrypted
	ch <- collector.zoneBandwidthSSLUnencrypted
	ch <- collector.zoneBandwidthContentType
	ch <- collector.zoneBandwidthCountry
	ch <- collector.zoneBandwidthHTTPStatus
}

func (collector *CloudflareCollector) Collect(ch chan<- prometheus.Metric) {

	var wg sync.WaitGroup
	zones := fetchZones()
	for _, z := range zones {
		wg.Add(1)
		go func(ID string, name string) {

			zt := fetchZoneTotals(ID)

			//Write latest value for each metric in the prometheus metric channel.
			//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
			ch <- prometheus.MustNewConstMetric(collector.zoneRequestTotal, prometheus.CounterValue, float64(zt.Requests.All), name)
			ch <- prometheus.MustNewConstMetric(collector.zoneRequestCached, prometheus.CounterValue, float64(zt.Requests.Cached), name)
			ch <- prometheus.MustNewConstMetric(collector.zoneRequestUncached, prometheus.CounterValue, float64(zt.Requests.Uncached), name)
			ch <- prometheus.MustNewConstMetric(collector.zoneRequestSSLEncrypted, prometheus.CounterValue, float64(zt.Requests.SSL.Encrypted), name)
			ch <- prometheus.MustNewConstMetric(collector.zoneRequestSSLUnencrypted, prometheus.CounterValue, float64(zt.Requests.SSL.Unencrypted), name)

			for ct, value := range zt.Requests.ContentType {
				ch <- prometheus.MustNewConstMetric(collector.zoneRequestContentType, prometheus.CounterValue, float64(value), name, ct)
			}

			for country, value := range zt.Requests.Country {
				ch <- prometheus.MustNewConstMetric(collector.zoneRequestCountry, prometheus.CounterValue, float64(value), name, country)
			}

			for status, value := range zt.Requests.HTTPStatus {
				ch <- prometheus.MustNewConstMetric(collector.zoneRequestHTTPStatus, prometheus.CounterValue, float64(value), name, status)
			}

			ch <- prometheus.MustNewConstMetric(collector.zoneBandwidthTotal, prometheus.CounterValue, float64(zt.Bandwidth.All), name)
			ch <- prometheus.MustNewConstMetric(collector.zoneBandwidthCached, prometheus.CounterValue, float64(zt.Bandwidth.Cached), name)
			ch <- prometheus.MustNewConstMetric(collector.zoneBandwidthUncached, prometheus.CounterValue, float64(zt.Bandwidth.Uncached), name)
			ch <- prometheus.MustNewConstMetric(collector.zoneBandwidthSSLEncrypted, prometheus.CounterValue, float64(zt.Bandwidth.SSL.Encrypted), name)
			ch <- prometheus.MustNewConstMetric(collector.zoneBandwidthSSLUnencrypted, prometheus.CounterValue, float64(zt.Bandwidth.SSL.Unencrypted), name)

			defer wg.Done()

		}(z.ID, z.Name)
	}
	wg.Wait()

}