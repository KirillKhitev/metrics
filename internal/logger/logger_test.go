package logger

import (
	"testing"
)

func TestInitialize(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				level: "info",
			},
			wantErr: false,
		},
		{
			name: "positive test #2",
			args: args{
				level: "debug",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Initialize(tt.args.level); (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}

			if lvl := Log.Level().String(); lvl != tt.args.level {
				t.Errorf("Level log = %s, want %s", lvl, tt.args.level)
			}
		})
	}
}
