package tests

import (
	"bytes"
	"context"
	"flag"
	"log"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"text/template"

	_ "github.com/gosthome/gosthome/components"
	"github.com/gosthome/gosthome/components/api/esphomeproto"
	"github.com/gosthome/gosthome/components/api/frameshakers"
	"github.com/gosthome/gosthome/core"
	"github.com/gosthome/gosthome/core/config"
	"github.com/google/go-cmp/cmp"
)

var sampleConfig = template.Must(template.New("").Parse(`
gosthome:
    name: testABC
    friendly_name: "Testing ABC"
    mac: {{ .MAC }}

api:
    address: "127.0.0.1"
    port: {{ .Port }}
    {{ if ne .Password ""}}password: "{{ .Password }}"{{end }}
    {{ if ne .NoisePSK ""}}
    encryption:
        key: "{{.NoisePSK }}"
    {{ end }}

demo:
`))

var wantPython = template.Must(template.New("").Parse(`API version: APIVersion(major=1, minor=10)
Device info: DeviceInfo(uses_password={{if eq .Password ""}}False{{else}}True{{end}}, name='testABC', friendly_name='Testing ABC', mac_address='{{.MAC}}', compilation_time='2022', model='{{.GOOS}}/{{.GOARCH}}', manufacturer='gosthome', has_deep_sleep=False, esphome_version='{{.Version}}', project_name='', project_version='', webserver_port=0, legacy_voice_assistant_version=0, voice_assistant_feature_flags=0, legacy_bluetooth_proxy_version=0, bluetooth_proxy_feature_flags=0, suggested_area='')

Entities:
- BinarySensorInfo(object_id='demo_movement_backyard', key=1756138606, name='Demo Movement Backyard', unique_id='', disabled_by_default=False, icon='', entity_category=<EntityCategory.NONE: 0>, device_class='motion', is_status_binary_sensor=False)
- BinarySensorInfo(object_id='demo_basement_floor_wet', key=2292024046, name='Demo Basement Floor Wet', unique_id='', disabled_by_default=False, icon='', entity_category=<EntityCategory.NONE: 0>, device_class='moisture', is_status_binary_sensor=False)
- ButtonInfo(object_id='demo_regenerate_seed', key=258008683, name='Demo Regenerate Seed', unique_id='', disabled_by_default=False, icon='', entity_category=<EntityCategory.NONE: 0>, device_class='restart')

State:
- BinarySensorState(key=1756138606, state=False, missing_state=False)
- BinarySensorState(key=2292024046, state=False, missing_state=False)
`))

var debug = flag.Bool("debug", false, "debug python output")

func TestPyNode(t *testing.T) {
	if testing.Verbose() {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}
	password := "abc"
	noise, err := frameshakers.GenerateEncryptionKey()
	if err != nil {
		t.Fatal(err)
	}
	nodeMac, err := config.GenerateMAC()
	if err != nil {
		t.Fatal(err)
	}

	type tcase struct {
		name     string
		password string
		noise    *frameshakers.ConfigNoisePSK
	}

	cases := []tcase{
		{
			name:     "TestPassword",
			password: password,
			noise:    nil,
		},
		{
			name:     "TestEncryption",
			password: "",
			noise:    noise,
		},
		{
			name:     "TestEncryptionAndPassword",
			password: password,
			noise:    noise,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()
			esphomeprotoPort := getFreePort(t)
			configBytes := &bytes.Buffer{}
			err = sampleConfig.Execute(configBytes, &struct {
				Port     int
				Password string
				NoisePSK string
				MAC      string
			}{
				Port:     esphomeprotoPort,
				Password: tcase.password,
				NoisePSK: tcase.noise.String(),
				MAC:      nodeMac.String(),
			})
			if err != nil {
				t.Fatal(err)
			}
			configData := make([]byte, configBytes.Len())
			copy(configData, configBytes.Bytes())
			cfg, err := config.LoadConfig(configBytes)
			if err != nil {
				print(string(configData))
				t.Fatal(err)
			}

			n, err := core.NewNode(context.Background(), cfg)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err = n.Close(); err != nil {
					t.Error(err)
				}
			}()
			n.Start()

			testdata := filepath.Join("..", "components", "api", "esphomeproto", "testdata")

			python := filepath.Join(testdata, "setup.sh")
			// For debugging. Don't key on testing.Verbose() since the test would be
			// failing.
			if *debug {
				cmd := exec.Command(
					python, filepath.Join(testdata, "test.py"),
					"--port", strconv.Itoa(esphomeprotoPort),
					"--password", tcase.password,
					"--noise-psk", tcase.noise.String(),
					"--verbose")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err = cmd.Run(); err != nil {
					t.Error(err)
				}
				t.FailNow()
			}
			out, err := exec.Command(
				python, filepath.Join(testdata, "test.py"),
				"--port", strconv.Itoa(esphomeprotoPort),
				"--password", tcase.password,
				"--noise-psk", tcase.noise.String(),
			).CombinedOutput()
			if err != nil {
				t.Error(err)
			}
			got := string(out)
			if runtime.GOOS == "windows" {
				got = strings.ReplaceAll(got, "\r", "")
			}
			want := bytes.Buffer{}
			// _, mac := getMainAddr()
			if err := wantPython.Execute(&want, map[string]string{
				"Password":     tcase.password,
				"GOOS":         runtime.GOOS,
				"GOARCH":       runtime.GOARCH,
				"MAC":          cfg.Gosthome.MAC.String(),
				"GOHO_VERSION": core.Version(),
				"Version":      esphomeproto.ESPHOME_VERSION,
			}); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want.String(), got); diff != "" {
				t.Errorf("python client mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func getFreePort(t *testing.T) int {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	p := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		t.Fatal(err)
	}
	return p
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}
