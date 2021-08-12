package manager

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/setting"
)

func TestLoader_Load(t *testing.T) {
	type fields struct {
		Cfg                           *setting.Cfg
		allowUnsignedPluginsCondition unsignedPluginV2ConditionFunc
		log                           log.Logger
	}
	type args struct {
		pluginJSONPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *plugins.PluginV2
		wantErr bool
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Loader{
				Cfg:                           tt.fields.Cfg,
				allowUnsignedPluginsCondition: tt.fields.allowUnsignedPluginsCondition,
				log:                           tt.fields.log,
			}
			got, err := l.Load(tt.args.pluginJSONPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoader_LoadAll(t *testing.T) {
	corePluginDir, err := filepath.Abs("./../../../public")
	if err != nil {
		t.Errorf("could not construct absolute path of core plugins dir")
		return
	}
	currentPath, err := filepath.Abs(".")
	if err != nil {
		t.Errorf("could not construct absolute path of current dir")
		return
	}
	tests := []struct {
		name            string
		cfg             *setting.Cfg
		log             log.Logger
		pluginJSONPaths []string
		want            []*plugins.PluginV2
		wantErr         bool
	}{
		{
			name: "Load a Core plugin",
			cfg: &setting.Cfg{
				StaticRootPath: corePluginDir,
			},
			log:             &FakeLogger{},
			pluginJSONPaths: []string{filepath.Join(corePluginDir, "app/plugins/datasource/cloudwatch/plugin.json")},
			want: []*plugins.PluginV2{
				{
					JSONData: plugins.JSONData{
						ID:   "cloudwatch",
						Type: "datasource",
						Name: "CloudWatch",
						Info: plugins.PluginInfo{
							Author: plugins.PluginInfoLink{
								Name: "Grafana Labs",
								Url:  "https://grafana.com",
							},
							Description: "Data source for Amazon AWS monitoring service",
							Logos: plugins.PluginLogos{
								Small: "img/amazon-web-services.png",
								Large: "img/amazon-web-services.png",
							},
						},
						Includes: []*plugins.PluginInclude{
							{Name: "EC2", Path: "dashboards/ec2.json", Type: "dashboard"},
							{Name: "EBS", Path: "dashboards/EBS.json", Type: "dashboard"},
							{Name: "Lambda", Path: "dashboards/Lambda.json", Type: "dashboard"},
							{Name: "Logs", Path: "dashboards/Logs.json", Type: "dashboard"},
							{Name: "RDS", Path: "dashboards/RDS.json", Type: "dashboard"},
						},
						Category:     "cloud",
						Signature:    "internal",
						Annotations:  true,
						Metrics:      true,
						Alerting:     true,
						Logs:         true,
						QueryOptions: map[string]bool{"minInterval": true},
					},
					PluginDir: filepath.Join(corePluginDir, "app/plugins/datasource/cloudwatch"),
					Class:     "core",
				},
			},
			wantErr: false,
		}, {
			name: "Load a Bundled plugin",
			cfg: &setting.Cfg{
				BundledPluginsPath: filepath.Join(currentPath, "testdata/unsigned-datasource"),
			},
			log:             &FakeLogger{},
			pluginJSONPaths: []string{"./testdata/unsigned-datasource/plugin/plugin.json"},
			want: []*plugins.PluginV2{
				{
					JSONData: plugins.JSONData{
						ID:   "test",
						Type: "datasource",
						Name: "Test",
						Info: plugins.PluginInfo{
							Author: plugins.PluginInfoLink{
								Name: "Grafana Labs",
								Url:  "https://grafana.com",
							},
							Description: "Test",
						},
						Backend:   true,
						Signature: "unsigned",
						State:     "alpha",
					},
					PluginDir: filepath.Join(currentPath, "testdata/unsigned-datasource/plugin/"),
					Class:     "bundled",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Loader{
				Cfg: tt.cfg,
				log: tt.log,
			}
			got, err := l.LoadAll(tt.pluginJSONPaths)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Fatalf("Result mismatch (-want +got):\n%s", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestLoader_readPluginJSON(t *testing.T) {
	tests := []struct {
		name           string
		pluginJSONPath string
		expected       plugins.JSONData
		failed         bool
	}{
		{
			name:           "Valid plugin",
			pluginJSONPath: "./testdata/test-app/plugin.json",
			expected: plugins.JSONData{
				ID:   "test-app",
				Type: "app",
				Name: "Test App",
				Info: plugins.PluginInfo{
					Author: plugins.PluginInfoLink{
						Name: "Test Inc.",
						Url:  "http://test.com",
					},
					Description: "Official Grafana Test App & Dashboard bundle",
					Version:     "1.0.0",
					Links: []plugins.PluginInfoLink{
						{Name: "Project site", Url: "http://project.com"},
						{Name: "License & Terms", Url: "http://license.com"},
					},
					Logos: plugins.PluginLogos{
						Small: "img/logo_small.png",
						Large: "img/logo_large.png",
					},
					Screenshots: []plugins.PluginScreenshots{
						{Path: "img/screenshot1.png", Name: "img1"},
						{Path: "img/screenshot2.png", Name: "img2"},
					},
					Updated: "2015-02-10",
				},
				Dependencies: plugins.PluginDependencies{
					GrafanaVersion: "3.x.x",
					Plugins: []plugins.PluginDependencyItem{
						{Type: "datasource", Id: "graphite", Name: "Graphite", Version: "1.0.0"},
						{Type: "panel", Id: "graph", Name: "Graph", Version: "1.0.0"},
					},
				},
				Includes: []*plugins.PluginInclude{
					{Name: "Nginx Connections", Path: "dashboards/connections.json", Type: "dashboard"},
					{Name: "Nginx Memory", Path: "dashboards/memory.json", Type: "dashboard"},
					{Name: "Nginx Panel", Type: "panel"},
					{Name: "Nginx Datasource", Type: "datasource"},
				},
				Backend: false,
			},
		},
		{
			name:           "Invalid plugin JSON",
			pluginJSONPath: "./testdata/invalid-plugin-json/plugin.json",
			failed:         true,
		},
		{
			name:           "Non-existing JSON file",
			pluginJSONPath: "./nonExistingFile.json",
			failed:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Loader{
				log: &FakeLogger{},
			}
			got, err := l.readPluginJSON(tt.pluginJSONPath)
			if (err != nil) && !tt.failed {
				t.Errorf("readPluginJSON() error = %v, failed %v", err, tt.failed)
				return
			}
			if !cmp.Equal(got, tt.expected) {
				t.Errorf("Unexpected pluginJSONData: %v", cmp.Diff(got, tt.expected))
			}
		})
	}
}

func Test_validatePluginJSON(t *testing.T) {
	type args struct {
		data plugins.JSONData
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Valid case",
			args: args{
				data: plugins.JSONData{
					ID:   "grafana-plugin-id",
					Type: plugins.DataSource,
				},
			},
		},
		{
			name: "Invalid plugin ID",
			args: args{
				data: plugins.JSONData{
					Type: plugins.Panel,
				},
			},
			err: InvalidPluginJSON,
		},
		{
			name: "Invalid plugin type",
			args: args{
				data: plugins.JSONData{
					ID:   "grafana-plugin-id",
					Type: "test",
				},
			},
			err: InvalidPluginJSON,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePluginJSON(tt.args.data); !errors.Is(err, tt.err) {
				t.Errorf("validatePluginJSON() = %v, want %v", err, tt.err)
			}
		})
	}
}
