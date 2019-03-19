package tilepix_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bcvery1/tilepix"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestReadFile(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		want     *tilepix.Map
		wantErr  bool
	}{
		{
			name:     "base64",
			filepath: "testdata/base64.tmx",
			want:     nil,
			wantErr:  false,
		},
		{
			name:     "base64-zlib",
			filepath: "testdata/base64-zlib.tmx",
			want:     nil,
			wantErr:  false,
		},
		{
			name:     "base64-gzip",
			filepath: "testdata/base64-gzip.tmx",
			want:     nil,
			wantErr:  false,
		},
		{
			name:     "csv",
			filepath: "testdata/csv.tmx",
			want:     nil,
			wantErr:  false,
		},
		{
			name:     "xml",
			filepath: "testdata/xml.tmx",
			want:     nil,
			wantErr:  false,
		},
		{
			name:     "missing file",
			filepath: "testdata/foo.tmx",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "map is infinite",
			filepath: "testdata/infinite.tmx",
			want:     nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tilepix.ReadFile(tt.filepath)

			if !tt.wantErr && err != nil {
				t.Errorf("tmx.ReadFile(): got unexpected error: %v", err)
			}
			if tt.wantErr && err == nil {
				t.Errorf("tmx.ReadFile(): expected error but not nil")
			}
		})
	}
}

func readFromFile(t *testing.T, filename string) (*tilepix.Map, error) {
	t.Log("Reading", filename)
	r, err := os.Open("testdata/poly.tmx")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	return tilepix.Read(r)
}

func TestProperties(t *testing.T) {
	m, err := tilepix.ReadFile("testdata/poly.tmx")
	if err != nil {
		t.Fatal(err)
	}
	for _, group := range m.ObjectGroups {
		for _, object := range group.Objects {
			if object.Properties[0].Name != "foo" {
				t.Error("No properties")
			}
			return
		}
	}
	t.Fatal("No property found")
}

func TestGetLayerByName(t *testing.T) {
	m, err := tilepix.ReadFile("testdata/poly.tmx")
	if err != nil {
		t.Fatal(err)
	}
	layer := m.GetLayerByName("Tile Layer 1")
	if layer.Name != "Tile Layer 1" {
		t.Error("error get layer")
	}
}

func TestGetObjectLayerByName(t *testing.T) {
	m, err := tilepix.ReadFile("testdata/poly.tmx")
	if err != nil {
		t.Fatal(err)
	}
	layer := m.GetObjectLayerByName("Object Layer 1")
	if layer.Name != "Object Layer 1" {
		t.Error("error get object layer")
	}
}
