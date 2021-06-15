package filter

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	resource2 "github.com/cloudskiff/driftctl/test/resource"

	"github.com/cloudskiff/driftctl/pkg/resource"
)

func TestDriftIgnore_IsResourceIgnored(t *testing.T) {
	tests := []struct {
		name      string
		resources []resource.Resource
		want      []bool
	}{
		{
			name: "drift_ignore_no_file",
			resources: []resource.Resource{
				&resource2.FakeResource{
					Type: "type1",
					Id:   "id1",
				},
			},

			want: []bool{
				false,
			},
		},
		{
			name: "drift_ignore_empty",
			resources: []resource.Resource{
				&resource2.FakeResource{
					Type: "type1",
					Id:   "id1",
				},
			},
			want: []bool{
				false,
			},
		},
		{
			name: "drift_ignore_invalid_lines",
			resources: []resource.Resource{
				&resource2.FakeResource{
					Type: "type1",
					Id:   "id1",
				},
				&resource2.FakeResource{
					Type: "ignored_resource",
					Id:   "id2",
				},
			},
			want: []bool{
				false,
				true,
			},
		},
		{
			name: "drift_ignore_valid",
			resources: []resource.Resource{
				&resource2.FakeResource{
					Type: "type1",
					Id:   "id1",
				},
				&resource2.FakeResource{
					Type: "wildcard_resource",
					Id:   "id1",
				},
				&resource2.FakeResource{
					Type: "wildcard_resource",
					Id:   "id2",
				},
				&resource2.FakeResource{
					Type: "wildcard_resource",
					Id:   "id3",
				},
				&resource2.FakeResource{
					Type: "ignored_resource",
					Id:   "id2",
				},
				&resource2.FakeResource{
					Type: "resource_type",
					Id:   "id.with.dots",
				},
				&resource2.FakeResource{
					Type: "resource_type",
					Id:   "idwith\\",
				},
				&resource2.FakeResource{
					Type: "resource_type",
					Id:   "idwith\\backslashes",
				},
			},
			want: []bool{
				false,
				true,
				true,
				true,
				true,
				true,
				true,
				true,
			},
		},
		{
			name: "drift_ignore_wildcard",
			resources: []resource.Resource{
				&resource2.FakeResource{
					Type: "type1",
					Id:   "id1",
				},
				&resource2.FakeResource{
					Type: "type2",
					Id:   "id1",
				},
				&resource2.FakeResource{
					Type: "type2",
					Id:   "id11",
				},
				&resource2.FakeResource{
					Type: "type2",
					Id:   "id2",
				},
				&resource2.FakeResource{
					Type: "type3",
					Id:   "id100",
				},
				&resource2.FakeResource{
					Type: "type3",
					Id:   "id101",
				},
				&resource2.FakeResource{
					Type: "type4",
					Id:   "id\\WithBac*slash***\\*\\",
				},
			},
			want: []bool{
				false,
				true,
				true,
				false,
				true,
				false,
				true,
			},
		},
		{
			name: "drift_ignore_all_exclude",
			resources: []resource.Resource{
				&resource2.FakeResource{
					Type: "type1",
					Id:   "id1",
				},
				&resource2.FakeResource{
					Type: "type2",
					Id:   "id1",
				},
				&resource2.FakeResource{
					Type: "type2",
					Id:   "id11",
				},
				&resource2.FakeResource{
					Type: "type2",
					Id:   "id2",
				},
				&resource2.FakeResource{
					Type: "type3",
					Id:   "id100",
				},
				&resource2.FakeResource{
					Type: "type3",
					Id:   "id101",
				},
				&resource2.FakeResource{
					Type: "iam_user",
					Id:   "id\\WithBac*slash***\\*\\",
				},
			},
			want: []bool{
				true,
				true,
				true,
				true,
				true,
				true,
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cwd, _ := os.Getwd()
			defer func() { _ = os.Chdir(cwd) }()
			if err := os.Chdir(path.Join("testdata", tt.name)); err != nil {
				t.Fatal(err)
			}
			r := NewDriftIgnore()
			got := make([]bool, 0, len(tt.want))
			for _, res := range tt.resources {
				got = append(got, r.IsResourceIgnored(res))
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDriftIgnore_IsFieldIgnored(t *testing.T) {

	type Args struct {
		Res  resource.Resource
		Path []string
		Want bool
	}

	tests := []struct {
		name string
		args []Args
	}{
		{
			name: "drift_ignore_no_file",
			args: []Args{

				{
					Res:  &resource2.FakeResource{Type: "type1", Id: "id1"},
					Path: []string{"Id"},
					Want: false,
				},
				{
					Res:  &resource2.FakeResource{Type: "type2", Id: "id2"},
					Path: []string{"Id"},
					Want: false,
				},
			},
		},
		{
			name: "drift_ignore_empty",
			args: []Args{
				{
					Res:  &resource2.FakeResource{Type: "type1", Id: "id1"},
					Path: []string{"Id"},
					Want: false,
				},
				{
					Res:  &resource2.FakeResource{Type: "type2", Id: "id2"},
					Path: []string{"Id"},
					Want: false,
				},
			},
		},
		{
			name: "drift_ignore_fields",
			args: []Args{
				{
					Res:  &resource2.FakeResource{Type: "res_type", Id: "full_drift_ignored"},
					Path: []string{"json"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "res_type", Id: "full_drift_ignored"},
					Path: []string{"foobar"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "res_type", Id: "partial_drift_ignored"},
					Path: []string{"json"},
					Want: false,
				},
				{
					Res:  &resource2.FakeResource{Type: "res_type", Id: "partial_drift_ignored"},
					Path: []string{"foobar"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "resource_type", Id: "id.with.dots"},
					Path: []string{"json"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "resource_type", Id: "id.with.dots"},
					Path: []string{"json"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "resource_type", Id: "idwith\\"},
					Path: []string{"json"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "resource_type", Id: "idwith\\backslashes"},
					Path: []string{"json"},
					Want: false,
				},
				{
					Res:  &resource2.FakeResource{Type: "resource_type", Id: "idwith\\backslashes"},
					Path: []string{"foobar"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "res_type", Id: "wildcard_drift_ignored"},
					Path: []string{"struct", "baz"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "res_type", Id: "wildcard_drift_ignored"},
					Path: []string{"struct", "bar"},
					Want: false,
				},
				{
					Res:  &resource2.FakeResource{Type: "res_type", Id: "endofpath_drift_ignored"},
					Path: []string{"struct", "baz"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "res_type", Id: "endofpath_drift_ignored"},
					Path: []string{"struct", "bar"},
					Want: true,
				},
			},
		},
		{
			name: "drift_ignore_all_exclude_field",
			args: []Args{
				{
					Res:  &resource2.FakeResource{Type: "res_type", Id: "full_drift_ignored"},
					Path: []string{"json"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "res_type", Id: "full_drift_ignored"},
					Path: []string{"foobar"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "res_type", Id: "partial_drift_ignored"},
					Path: []string{"json"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "res_type", Id: "partial_drift_ignored"},
					Path: []string{"foobar"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "resource_type", Id: "id.with.dots"},
					Path: []string{"json"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "resource_type", Id: "id.with.dots"},
					Path: []string{"json"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "resource_type", Id: "idwith\\"},
					Path: []string{"json"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "resource_type", Id: "idwith\\backslashes"},
					Path: []string{"json"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "resource_type", Id: "idwith\\backslashes"},
					Path: []string{"foobar"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "res_type", Id: "wildcard_drift_ignored"},
					Path: []string{"struct", "baz"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "res_type", Id: "wildcard_drift_ignored"},
					Path: []string{"struct", "bar"},
					Want: false,
				},
				{
					Res:  &resource2.FakeResource{Type: "res_type", Id: "endofpath_drift_ignored"},
					Path: []string{"struct", "baz"},
					Want: true,
				},
				{
					Res:  &resource2.FakeResource{Type: "res_type", Id: "endofpath_drift_ignored"},
					Path: []string{"struct", "bar"},
					Want: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cwd, _ := os.Getwd()
			defer func() { _ = os.Chdir(cwd) }()
			if err := os.Chdir(path.Join("testdata", tt.name)); err != nil {
				t.Fatal(err)
			}
			r := NewDriftIgnore()
			for _, arg := range tt.args {
				got := r.IsFieldIgnored(arg.Res, arg.Path)
				if arg.Want != got {
					t.Errorf("%s.%s.%s expected %v got %v", arg.Res.TerraformType(), arg.Res.TerraformId(), strings.Join(arg.Path, "."), arg.Want, got)
				}
			}
		})
	}
}
