/*
 * Copyright (c) 2024. Devtron Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import "testing"

func TestCompareLimitsRequests(t *testing.T) {
	requests := "requests"
	limits := "limits"
	resources := "resources"
	cpu := "cpu"
	memory := "memory"
	type args struct {
		dat map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "empty base object",
			args:    args{dat: nil},
			want:    true,
			wantErr: false,
		},
		{
			name:    "empty resources object",
			args:    args{dat: map[string]interface{}{resources: nil}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "empty resources requests limits object",
			args:    args{dat: map[string]interface{}{resources: map[string]interface{}{limits: nil, requests: nil}}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-empty resources limits object",
			args:    args{dat: map[string]interface{}{resources: map[string]interface{}{limits: map[string]interface{}{cpu: "10m", memory: "10Mi"}, requests: nil}}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-empty resources requests object",
			args:    args{dat: map[string]interface{}{resources: map[string]interface{}{limits: nil, requests: map[string]interface{}{cpu: "10m", memory: "12Gi"}}}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-empty  and equal resources limits and requests object",
			args:    args{dat: map[string]interface{}{resources: map[string]interface{}{limits: map[string]interface{}{cpu: "0.11", memory: "10Gi"}, requests: map[string]interface{}{cpu: "11m", memory: "9Gi"}}}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "negative: non-empty  and not equal resources limits and requests object",
			args:    args{dat: map[string]interface{}{resources: map[string]interface{}{limits: map[string]interface{}{cpu: "0.10", memory: "10Mi"}, requests: map[string]interface{}{cpu: "111m", memory: "15Gi"}}}},
			want:    false,
			wantErr: true,
		},
		{
			name:    "negative test cases - 1",
			args:    args{dat: map[string]interface{}{resources: map[string]interface{}{limits: map[string]interface{}{cpu: "0.1.0", memory: "10Mi"}, requests: map[string]interface{}{cpu: "111m", memory: "15Gi"}}}},
			want:    false,
			wantErr: true,
		},
		{
			name:    "negative test cases - 2",
			args:    args{dat: map[string]interface{}{resources: map[string]interface{}{limits: map[string]interface{}{cpu: "-0.10", memory: "10Mi"}, requests: map[string]interface{}{cpu: "111m", memory: "15Gi"}}}},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompareLimitsRequests(tt.args.dat)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareLimitsRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CompareLimitsRequests() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAutoScale(t *testing.T) {
	autoScaling := "autoscaling"
	minReplicas := "MinReplicas"
	maxReplicas := "MaxReplicas"
	enabled := "enabled"
	type args struct {
		dat map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "empty base object",
			args:    args{dat: nil},
			want:    true,
			wantErr: false,
		},
		{
			name:    "empty autoscaling object",
			args:    args{dat: map[string]interface{}{autoScaling: nil}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "negative : non-empty autoscaling empty enabled minReplicas maxReplicas object",
			args:    args{dat: map[string]interface{}{autoScaling: map[string]interface{}{}}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-empty autoscaling enabled minReplicas maxReplicas object",
			args:    args{dat: map[string]interface{}{autoScaling: map[string]interface{}{enabled: false, minReplicas: float64(10), maxReplicas: float64(11)}}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-empty autoscaling enabled minReplicas empty maxReplicas object",
			args:    args{dat: map[string]interface{}{autoScaling: map[string]interface{}{enabled: false, minReplicas: float64(11)}}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-empty autoscaling minReplicas maxReplicas object empty enabled",
			args:    args{dat: map[string]interface{}{autoScaling: map[string]interface{}{minReplicas: float64(10), maxReplicas: float64(11)}}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "negative: non-empty and greater minReplicas than maxReplicas object",
			args:    args{dat: map[string]interface{}{autoScaling: map[string]interface{}{enabled: true, minReplicas: float64(10), maxReplicas: float64(9)}}},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AutoScale(tt.args.dat)
			if (err != nil) != tt.wantErr {
				t.Errorf("AutoScale() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AutoScale() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertMemory(t *testing.T) {
	type args struct {
		memory string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "base test",
			args:    args{memory: "1Gi"},
			want:    1073741824,
			wantErr: false,
		},
		{
			name:    "negative test - scientifc notation",
			args:    args{memory: "1e2G"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "base test - scientifc notation",
			args:    args{memory: "1e2"},
			want:    100,
			wantErr: false,
		},
		{
			name:    "negative test case - Memory1",
			args:    args{memory: "1.0.1Mi"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "negative test case - Memory2",
			args:    args{memory: "-10Mi"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "negative test case - Memory2",
			args:    args{memory: "1Ki"},
			want:    1024,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MemoryToNumber(tt.args.memory)
			if (err != nil) != tt.wantErr {
				t.Errorf("memory() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("memory() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertCPU(t *testing.T) {
	type args struct {
		cpu string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "base test with unit",
			args:    args{cpu: "10m"},
			want:    10,
			wantErr: false,
		},
		{
			name:    "base test without unit",
			args:    args{cpu: "0.01"},
			want:    10,
			wantErr: false,
		},
		{
			name:    "base test - scientifc notation",
			args:    args{cpu: "1e2"},
			want:    100000,
			wantErr: false,
		},
		{
			name:    "negative test case - Cpu1",
			args:    args{cpu: "1.0.1"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "negative test case - Cpu2",
			args:    args{cpu: "-10m"},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CpuToNumber(tt.args.cpu)
			if (err != nil) != tt.wantErr {
				t.Errorf("cpu() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("memory() got = %v, want %v", got, tt.want)
			}
		})
	}
}
