package expect

import (
	"os"
	"path/filepath"

	_ "image/gif"
	_ "image/jpeg"

	"github.com/ghodss/yaml"
	"golang.org/x/exp/slices"
)

func (e Val) ToBeSnapshot(path string) Val {
	e.t.Helper()

	folder := filepath.Dir(path)
	if folder != "." {
		err := os.MkdirAll(folder, 0o755)
		if err != nil {
			e.t.Fatalf("failed to create target folder %v", folder)
		}
	}

	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		e.t.Fatalf("failed to read snaphsot %v: %v", path, err)
	}

	current, err := asBytes(e.value)
	if err != nil {
		e.t.Error(err)
	}

	if existing == nil {
		// snapshot does not exist, create it
		err = os.WriteFile(path, current, 0o644)
		if err != nil {
			e.t.Fatalf("failed to write snapshot %v", path)
		}
	} else {
		if slices.Equal(current, existing) {
			// all is well, snapshot is matched, remove a possible current version
			os.RemoveAll(path + ".current")
		} else {
			e.t.Errorf("snapshot for %v does not match current output", path)
			err = os.WriteFile(path+".current", current, 0o644)
			if err != nil {
				e.t.Fatalf("failed to write snapshot %v", path)
			}
		}
	}

	return e
}

func asBytes(in any) ([]byte, error) {
	switch t := in.(type) {
	case []byte:
		return t, nil
	case string:
		return []byte(t), nil
	default:
		return yaml.Marshal(in)
	}
}
