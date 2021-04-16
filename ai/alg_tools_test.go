package ai

import (
	"board-game/core"
	"reflect"
	"testing"
)

func TestSpin90(t *testing.T) {
	type args struct {
		mat [][]*core.PlaySignal
	}
	tests := []struct {
		name    string
		args    args
		want    [][]*core.PlaySignal
		wantErr bool
	}{
		{
			name: "nil, should original",
			args: args{
				mat: nil,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "empty, should original",
			args: args{
				mat: [][]*core.PlaySignal{},
			},
			want:    [][]*core.PlaySignal{},
			wantErr: false,
		},
		{
			name: "not square, should error",
			args: args{
				mat: [][]*core.PlaySignal{
					{
						&core.PlaySignal{},
						&core.PlaySignal{},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "square, should spin",
			args: args{
				mat: [][]*core.PlaySignal{
					{
						&core.PlaySignal{Tag: "1"},
						&core.PlaySignal{Tag: "2"},
						&core.PlaySignal{Tag: "3"},
					},
					{
						&core.PlaySignal{Tag: "4"},
						&core.PlaySignal{Tag: "5"},
						&core.PlaySignal{Tag: "6"},
					},
					{
						&core.PlaySignal{Tag: "7"},
						&core.PlaySignal{Tag: "8"},
						&core.PlaySignal{Tag: "9"},
					},
				},
			},
			want: [][]*core.PlaySignal{
				{
					&core.PlaySignal{Tag: "7"},
					&core.PlaySignal{Tag: "4"},
					&core.PlaySignal{Tag: "1"},
				},
				{
					&core.PlaySignal{Tag: "8"},
					&core.PlaySignal{Tag: "5"},
					&core.PlaySignal{Tag: "2"},
				},
				{
					&core.PlaySignal{Tag: "9"},
					&core.PlaySignal{Tag: "6"},
					&core.PlaySignal{Tag: "3"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Spin90(tt.args.mat)
			if (err != nil) != tt.wantErr {
				t.Errorf("Spin90() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Spin90() got = %v, want %v", got, tt.want)
			}
		})
	}
}
