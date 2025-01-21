#!/bin/bash

# Minta input nama folder dari pengguna
read -p "Masukkan nama folder: " folder_name

# Cek apakah folder sudah ada
if [ -d "$folder_name" ]; then
    echo "Folder '$folder_name' sudah ada!"
    exit 1
fi

# Buat folder baru
mkdir "$folder_name"

# Pindah ke dalam folder yang baru dibuat
cd "$folder_name" || exit

# Buat file main.go dengan template dasar
cat > main.go <<EOL
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
EOL

echo "Folder '$folder_name' dan file 'main.go' berhasil dibuat!"

