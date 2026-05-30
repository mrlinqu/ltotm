package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("aaaa")
		return
	}

	root := os.Args[1]

	m := minify.New()
	m.AddFunc(".css", css.Minify)
	m.AddFunc(".js", js.Minify)
	m.AddFunc("html", html.Minify)

	fileSystem := os.DirFS(root)

	//includes := map[string]*bytes.Buffer{}
	css := strings.Builder{}
	js := strings.Builder{}

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if path == "." || path == "index.html" || path == "bundle.html" {
			return nil
		}

		ext := filepath.Ext(path)

		if ext != ".js" && ext != ".css" {
			return nil
		}

		fmt.Printf("minify %s... ", path)

		fullPath := filepath.Join(root, path)
		buf := &bytes.Buffer{}

		f, err := os.Open(fullPath)
		if err != nil {
			log.Fatal(err)
		}

		if err := m.Minify(ext, buf, f); err != nil {
			log.Fatal(err)
		}

		if err := f.Close(); err != nil {
			log.Fatal(err)
		}

		switch ext {
		case ".js":
			js.WriteString(`<script type="text/javascript">`)
			js.Write(buf.Bytes())
			js.WriteString("</script>")
		case ".css":
			css.WriteString("<style>")
			css.Write(buf.Bytes())
			css.WriteString("</style>")
		}

		fmt.Printf("ok\n")

		return nil
	})

	fmt.Printf("minify index.html... ")

	fileNameIndex := filepath.Join(root, "index.html")
	in, err := os.ReadFile(fileNameIndex)
	if err != nil {
		log.Fatal(err)
	}

	str := strings.Replace(string(in), "{{css}}", css.String(), 1)
	str = strings.Replace(str, "{{js}}", js.String(), 1)

	fileNameBundle := filepath.Join(root, "bundle.html")
	err = os.WriteFile(fileNameBundle, []byte(str), 0666)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ok\n")

	// buf := &bytes.Buffer{}

	// f, err := os.Open(fullPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := m.Minify(ext, buf, f); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := f.Close(); err != nil {
	// 	log.Fatal(err)
	// }

	//fmt.Println(includes)
}
