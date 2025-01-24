package main

import (
	"math/rand"
)

var packages = []string{
	"github.com/go-chi/chi/v5",
	"github.com/redis/go-redis/v9",
	"github.com/joho/godotenv",
	"github.com/sirupsen/logrus",
	"github.com/go-sql-driver/mysql",
	"gopkg.in/gomail.v2",
	"github.com/go-chi/chi/middleware",
}

// func getPackage() []string {
// 	pkgs := packages
// 	copy(pkgs, packages)
//
// 	rand.Shuffle(len(pkgs), func(i, j int) {
// 		pkgs[i], pkgs[j] = pkgs[j], pkgs[i]
// 	})
// 	for k := range pkgs {
// 		pkgs[k] = fmt.Sprintf("-%d.%d.%d", rand.Intn(10), rand.Intn(10), rand.Intn(10))
// 	}
// 	return pkgs
// }
//

func getPackage() []string {
	pkgs := make([]string, len(packages))
	copy(pkgs, packages)

	rand.Shuffle(len(pkgs), func(i, j int) {
		pkgs[i], pkgs[j] = pkgs[j], pkgs[i]
	})

	return pkgs // Mengembalikan daftar paket tanpa modifikasi
}
