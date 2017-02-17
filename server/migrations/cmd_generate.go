package migrations

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

var tmpl = template.Must(template.New("").Parse(
	`package {{ .package }}

import (
	m "github.com/mwheatley3/ww/server/migrations"
)

func init() {
	Migrations.Add("{{ .name }}",
		func(db m.DB) error {
			return db.Exec(` + "``)" + `
		},

		func(db m.DB) error {
			return db.Exec(` + "``)" + `
		},
	)
}
`))

func gen(dir, pkg, name string) (string, error) {
	m, err := ioutil.ReadDir(dir)

	if err != nil {
		return "", err
	}

	var latest int

	for _, f := range m {
		num, err := strconv.Atoi(f.Name()[0:3])

		if err != nil {
			continue
		}

		if num > latest {
			latest = num
		}
	}

	newName := strconv.Itoa(latest + 1)
	newName = strings.Repeat("0", 3-len(newName)) + newName + "_" + name

	p := filepath.Join(dir, newName+".go")

	f, err := os.Create(p)
	defer f.Close()

	if err != nil {
		return "", err
	}

	if err := tmpl.Execute(f, map[string]interface{}{
		"name":    newName,
		"package": pkg,
	}); err != nil {
		return "", err
	}

	return p, nil
}
