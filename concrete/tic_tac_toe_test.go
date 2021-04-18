package concrete

import (
	"board-game/core"
	"reflect"
	"testing"
)

func TestTicTacToe_Compress(t1 *testing.T) {
	type fields struct {
		players []*core.Player
	}
	type args struct {
		mat [][]*core.PlaySignal
	}
	a, b := &core.PlaySignal{}, &core.PlaySignal{}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int32
	}{
		{
			name: "should success",
			fields: fields{
				players: []*core.Player{{Signal: a}, {Signal: b}},
			},
			args: args{
				mat: [][]*core.PlaySignal{
					{a, b, a},
					{b, a, b},
					{b, b, a},
				},
			},
			want: 0b_01_10_10_10_01_10_01_10_01,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TicTacToe{
				players: tt.fields.players,
			}
			if got := t.Zip(tt.args.mat); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Zip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTicTacToe_GenSimilar(t1 *testing.T) {
	type fields struct {
		players []*core.Player
	}
	type args struct {
		base [][]*core.PlaySignal
	}
	a, b := &core.PlaySignal{}, &core.PlaySignal{}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []int32
		wantErr bool
	}{
		{
			name: "should return 1 similar",
			fields: fields{
				players: []*core.Player{{Signal: a}, {Signal: b}},
			},
			args: args{
				base: [][]*core.PlaySignal{
					{nil, nil, nil},
					{nil, nil, nil},
					{nil, nil, nil},
				}},
			want:    []int32{0},
			wantErr: false,
		},
		{
			name: "should return 4 similar",
			fields: fields{
				players: []*core.Player{{Signal: a}, {Signal: b}},
			},
			args: args{
				base: [][]*core.PlaySignal{
					{a, nil, nil},
					{nil, nil, nil},
					{nil, nil, nil},
				}},
			want: []int32{
				0b_00_00_00_00_00_00_00_00_01,
				0b_00_00_00_00_00_00_01_00_00,
				0b_00_00_01_00_00_00_00_00_00,
				0b_01_00_00_00_00_00_00_00_00,
			},
			wantErr: false,
		},
		{
			name: "should return 8 similar",
			fields: fields{
				players: []*core.Player{{Signal: a}, {Signal: b}},
			},
			args: args{
				base: [][]*core.PlaySignal{
					{a, b, nil},
					{nil, nil, nil},
					{nil, nil, nil},
				}},
			want: []int32{
				0b_00_00_00_00_00_00_00_10_01,
				0b_00_00_00_00_00_00_01_10_00,
				0b_00_00_00_00_00_10_00_00_01,
				0b_00_00_00_10_00_00_01_00_00,
				0b_00_00_01_00_00_10_00_00_00,
				0b_00_10_01_00_00_00_00_00_00,
				0b_01_00_00_10_00_00_00_00_00,
				0b_01_10_00_00_00_00_00_00_00,

				//0b_00_00_00_00_00_00_00_10_01,
				//0b_00_00_00_10_00_00_01_00_00,
				//0b_01_10_00_00_00_00_00_00_00,
				//0b_00_00_01_00_00_10_00_00_00,
				//
				//0b_00_00_00_00_00_00_01_10_00,
				//0b_01_00_00_10_00_00_00_00_00,
				//0b_00_10_01_00_00_00_00_00_00,
				//0b_00_00_00_00_00_10_00_00_01,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TicTacToe{
				players: tt.fields.players,
			}
			got, err := t.GenSimilar(tt.args.base)
			if (err != nil) != tt.wantErr {
				t1.Errorf("GenSimilar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GenSimilar() got = %v, want %v", got, tt.want)
			}
		})
	}
}
