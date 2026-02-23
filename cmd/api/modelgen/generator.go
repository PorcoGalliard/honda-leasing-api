package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	schemas := []string{"account", "mst", "dealer", "leasing", "finance"}

	for _, schema := range schemas {
		dsn := fmt.Sprintf(
			"host=localhost user=suryana password=suryana321 dbname=leasing_db port=5434 sslmode=disable search_path=%s",
			schema,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		g := gen.NewGenerator(gen.Config{
			OutPath:      "internal/domain/query/" + schema,
			ModelPkgPath: "internal/models",
			Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		})

		g.UseDB(db)
		g.ApplyBasic(g.GenerateAllTable()...)
		g.Execute()
	}
}
