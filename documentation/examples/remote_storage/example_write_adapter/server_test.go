// Copyright 2024 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"

	"github.com/prometheus/prometheus/prompb"
	writev2 "github.com/prometheus/prometheus/prompb/io/prometheus/write/v2"
)

func TestPrintV1(t *testing.T) {
	req := &prompb.WriteRequest{
		Timeseries: []prompb.TimeSeries{
			{
				Labels: []prompb.Label{
					{Name: "__name__", Value: "test_metric"},
					{Name: "label1", Value: "value1"},
				},
				Samples: []prompb.Sample{
					{Value: 1.0, Timestamp: 1000},
				},
				Exemplars: []prompb.Exemplar{
					{
						Labels:    []prompb.Label{{Name: "trace_id", Value: "abc123"}},
						Value:     1.5,
						Timestamp: 1500,
					},
				},
				Histograms: []prompb.Histogram{
					{
						Count:          &prompb.Histogram_CountInt{CountInt: 10},
						Sum:            100.0,
						Schema:         1,
						ZeroThreshold:  0.001,
						ZeroCount:      &prompb.Histogram_ZeroCountInt{ZeroCountInt: 2},
						PositiveSpans:  []prompb.BucketSpan{{Offset: 0, Length: 2}},
						PositiveDeltas: []int64{1, 1},
					},
				},
			},
		},
	}

	printV1(req)
}

func TestPrintV2(t *testing.T) {
	req := &writev2.Request{
		Symbols: []string{"", "__name__", "test_metric", "trace_id", "abc123"},
		Timeseries: []writev2.TimeSeries{
			{
				LabelsRefs: []uint32{1, 2},
				Samples: []writev2.Sample{
					{Value: 1.0, Timestamp: 1000},
				},
				Exemplars: []writev2.Exemplar{
					{
						LabelsRefs: []uint32{3, 4},
						Value:      1.5,
						Timestamp:  1500,
					},
				},
				Histograms: []writev2.Histogram{
					{
						Count:          &writev2.Histogram_CountInt{CountInt: 10},
						Sum:            100.0,
						Schema:         1,
						ZeroThreshold:  0.001,
						ZeroCount:      &writev2.Histogram_ZeroCountInt{ZeroCountInt: 2},
						PositiveSpans:  []writev2.BucketSpan{{Offset: 0, Length: 2}},
						PositiveDeltas: []int64{1, 1},
					},
				},
			},
		},
	}

	if err := printV2(req); err != nil {
		t.Fatalf("printV2 returned unexpected error: %v", err)
	}
}

func TestPrintV2Errors(t *testing.T) {
	req := &writev2.Request{
		Symbols: []string{"", "__name__", "test_metric"},
		Timeseries: []writev2.TimeSeries{
			{
				LabelsRefs: []uint32{1, 100},
				Samples: []writev2.Sample{
					{Value: 1.0, Timestamp: 1000},
				},
			},
		},
	}

	if err := printV2(req); err == nil {
		t.Fatal("printV2 should return error for invalid label reference")
	}
}
