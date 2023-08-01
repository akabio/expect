package expect

import (
	"os"
	"path/filepath"
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

	current, is := e.value.(string)
	if !is {
		e.t.Fatalf("value of .ToBeSnaphsot mus be of type string but it is %T", e.value)
	}

	if existing == nil {
		// snapshot does not exist, create it
		err = os.WriteFile(path, []byte(current), 0o644)
		if err != nil {
			e.t.Fatalf("failed to write snapshot %v", path)
		}
	} else {
		if current == string(existing) {
			// all is well, snapshot is matched, remove a possible current version
			os.RemoveAll(path + ".current")
		} else {
			e.t.Errorf("snapshot for %v does not match current output", path)
			err = os.WriteFile(path+".current", []byte(current), 0o644)
			if err != nil {
				e.t.Fatalf("failed to write snapshot %v", path)
			}
		}
	}

	return e
}
