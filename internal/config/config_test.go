package configuration

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type field struct {
	helper     bool
	maxDepth   int64
	jsonOutput bool
	fileExt    string
	logLevel   zerolog.Level
}

func TestNew(t *testing.T) {
	cfg := New()
	assert.NotNilf(t, cfg, "Config is Nil = %v", cfg)
}

func Test_Helper(t *testing.T) {
	tests := []struct {
		name    string
		fields  field
		want    bool
		wantErr bool
	}{
		{
			name: "success result",
			fields: field{
				helper:     false,
				maxDepth:   0,
				jsonOutput: false,
				fileExt:    ".go",
				logLevel:   0,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success result - want error ",
			fields: field{
				helper:     true,
				maxDepth:   0,
				jsonOutput: false,
				fileExt:    ".go",
				logLevel:   0,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "error result",
			fields: field{
				helper:     true,
				maxDepth:   0,
				jsonOutput: false,
				fileExt:    ".go",
				logLevel:   0,
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				helper:     tt.fields.helper,
				maxDepth:   tt.fields.maxDepth,
				jsonOutput: false,
				fileExt:    DefaultFileExt,
				logLevel:   tt.fields.logLevel,
			}
			if got := c.Helper(); got != tt.want && !tt.wantErr {
				t.Errorf("Helper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_JsonLog(t *testing.T) {
	tests := []struct {
		name    string
		fields  field
		want    bool
		wantErr bool
	}{
		{
			name: "success result",
			fields: field{
				helper:     false,
				maxDepth:   0,
				jsonOutput: false,
				fileExt:    ".go",
				logLevel:   0,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success result - want error",
			fields: field{
				helper:     false,
				maxDepth:   0,
				jsonOutput: false,
				fileExt:    ".go",
				logLevel:   0,
			},
			want:    true,
			wantErr: true,
		},
		{
			name: "error result",
			fields: field{
				helper:     true,
				maxDepth:   0,
				jsonOutput: true,
				fileExt:    ".go",
				logLevel:   0,
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				helper:     tt.fields.helper,
				maxDepth:   tt.fields.maxDepth,
				jsonOutput: false,
				fileExt:    DefaultFileExt,
				logLevel:   tt.fields.logLevel,
			}
			if got := c.JSONnLog(); got != tt.want && !tt.wantErr {
				t.Errorf("JSONnLog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FileExt(t *testing.T) {
	tests := []struct {
		name    string
		fields  field
		want    string
		wantErr bool
	}{
		{
			name: "success result",
			fields: field{
				helper:     false,
				maxDepth:   0,
				jsonOutput: false,
				fileExt:    ".go",
				logLevel:   0,
			},
			want:    ".go",
			wantErr: false,
		},
		{
			name: "error result",
			fields: field{
				helper:     true,
				maxDepth:   0,
				jsonOutput: false,
				fileExt:    ".gow",
				logLevel:   0,
			},
			want:    ".go",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				helper:     tt.fields.helper,
				maxDepth:   tt.fields.maxDepth,
				jsonOutput: false,
				fileExt:    DefaultFileExt,
				logLevel:   tt.fields.logLevel,
			}
			if got := c.FileExt(); got != tt.want && !tt.wantErr {
				t.Errorf("JSONnLog() = %v, want %v", got, tt.want)
			}
		})
	}
}
