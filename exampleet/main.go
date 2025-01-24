package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type TemplateData struct {
	ProjectName string
	Version     string
}

func main() {
	data := TemplateData{
		ProjectName: "MyApp",
		Version:     "1.0.0",
	}

	templateDir := "template/"
	for folder, files := range projectArch("./") {
		for _, file := range files {
			// filePath := filepath.Join(folder, file)
			fmt.Println(templateDir)
			templatePath := filepath.Join(templateDir, file+".template")
			outputPath := filepath.Join(folder, file)
			content, err := os.ReadFile(templatePath)
			if err != nil {
				fmt.Printf("gagal membuat project:%v\n", err)
				content = []byte("// default content\n")
			}

			// parser template dengan subtitusi variable
			tmpl, err := template.New(file).Parse(string(content))
			if err != nil {
				fmt.Printf("gagal membuat template:%v\n", err)
			}

			// // tulis ke file baru
			// err = os.WriteFile(filePath, content, 0644)
			// if err != nil {
			// 	fmt.Printf("Gagal menulis file %s: %v\n", filePath, err)
			// } else {
			// 	fmt.Printf("Berhasil membuat file %s\n", filePath)
			// }

			file, err := os.Create(outputPath)
			if err != nil {
				fmt.Printf("Gagal mem-parsing template %s: %v\n", templatePath, err)
			}
			defer file.Close()

			err = tmpl.Execute(file, data)
			if err != nil {
				fmt.Printf("Gagal menulis ke file %s: %v\n", outputPath, err)
				continue
			}
			fmt.Printf("Berhasil membuat file %s\n", outputPath)
		}
	}

}

func projectArch(rootDir string) map[string][]string {
	structures := map[string][]string{
		fmt.Sprintf("%s/exampleet", rootDir):                {"testx.go"},
		fmt.Sprintf("%s/exampleet/commons/helper", rootDir): {"helper.go"},
	}
	return structures
}
