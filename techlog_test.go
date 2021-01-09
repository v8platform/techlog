package techlog

import (
	"github.com/k0kubun/pp"
	"testing"
)

func TestWatch(t *testing.T) {
	type args struct {
		dir  string
		opts []Options
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"small chunk",
			args{
				dir: "./log",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs, _ := watch(tt.args.dir, tt.args.opts...)

			for log := range logs {
				pp.Println(log)
			}

		})
	}
}

func TestStreamEvents(t *testing.T) {
	type args struct {
		file      string
		maxEvents int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"simple",
			args{
				file:      "./log/18100509.log",
				maxEvents: 10,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream, err := StreamRead(tt.args.file, tt.args.maxEvents)
			if (err != nil) != tt.wantErr {
				t.Errorf("StreamRead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			count := 0
			for _ = range stream {
				count++
			}

			pp.Println("events", count)

		})
	}
}

func TestRead(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"simple",
			args{
				file: "./log/18100509.log",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			pp.Println("events count", len(got))

		})
	}
}
