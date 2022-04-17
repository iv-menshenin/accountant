package uuid

import "testing"

func TestUUID_fromUUIDString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		want    UUID
		args    args
		wantErr bool
	}{
		{
			name:    "2fdc1211-136e-4cf4-a849-b5491f8afb2d",
			want:    UUID{47, 220, 18, 17, 19, 110, 76, 244, 168, 73, 181, 73, 31, 138, 251, 45},
			args:    args{s: "2fdc1211-136e-4cf4-a849-b5491f8afb2d"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u UUID
			if err := u.FromString(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("FromString() error = %v, wantErr %v", err, tt.wantErr)
			}
			for n := 0; n < 16; n++ {
				if u[n] != tt.want[n] {
					t.Errorf("wrong UUID data: #%d %d != %d", n, u[n], tt.want[n])
				}
			}
		})
	}
}

func TestUUID_String(t *testing.T) {
	tests := []struct {
		name string
		arg  UUID
		want string
	}{
		{
			name: "nil",
			arg:  UUID{},
			want: "00000000-0000-0000-0000-000000000000",
		},
		{
			name: "2fdc1211-136e-4cf4-a849-b5491f8afb2d",
			arg:  UUID{47, 220, 18, 17, 19, 110, 76, 244, 168, 73, 181, 73, 31, 138, 251, 45},
			want: "2fdc1211-136e-4cf4-a849-b5491f8afb2d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
