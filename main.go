package main

import (
	"encoding/csv"
	"os"
	"path"
	"strconv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GeoNames struct {
	GeoNameID           int    `gorm:"column:geoname_id"`
	LocaleCode          string `gorm:"column:locale_code"`
	ContinentCode       string `gorm:"column:continent_code"`
	ContinentName       string `gorm:"column:continent_name"`
	CountryIsoCode      string `gorm:"column:country_iso_code"`
	CountryName         string `gorm:"column:country_name"`
	Subdivision1IsoCode string `gorm:"column:subdivision_1_iso_code"`
	Subdivision1Name    string `gorm:"column:subdivision_1_name"`
	Subdivision2IsoCode string `gorm:"column:subdivision_2_iso_code"`
	Subdivision2Name    string `gorm:"column:subdivision_2_name"`
	CityName            string `gorm:"column:city_name"`
	MetroCode           string `gorm:"column:metro_code"`
	TimeZone            string `gorm:"column:time_zone"`
	IsInEuropeanUnion   string `gorm:"column:is_in_european_union"`
}

// importCsvToSqlite imports a CSV file into a SQLite database.
func importCsvToSqlite(dataDir string, csvFile string, geonamesdbFile string) error {
	geonames, err := loadGeonamesCsv(csvFile)
	if err != nil {
		return err
	}

	db, err := gorm.Open(sqlite.Open(path.Join(dataDir, geonamesdbFile)), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Silent),
		CreateBatchSize: 1000,
	})
	if err != nil {
		return err
	}
	defer func() {
		sql, err := db.DB()
		if err != nil {
			return
		}
		sql.Close()
	}()

	if err := db.AutoMigrate(&GeoNames{}); err != nil {
		return err
	}

	return db.Create(geonames).Error
}

func loadGeonamesCsv(filepath string) ([]GeoNames, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var geoNames []GeoNames
	for index, record := range records {
		if index == 0 {
			continue
		}
		geoNameID, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}

		geoName := GeoNames{
			GeoNameID:           geoNameID,
			LocaleCode:          record[1],
			ContinentCode:       record[2],
			ContinentName:       record[3],
			CountryIsoCode:      record[4],
			CountryName:         record[5],
			Subdivision1IsoCode: record[6],
			Subdivision1Name:    record[7],
			Subdivision2IsoCode: record[8],
			Subdivision2Name:    record[9],
			CityName:            record[10],
			MetroCode:           record[11],
			TimeZone:            record[12],
			IsInEuropeanUnion:   record[13],
		}
		geoNames = append(geoNames, geoName)
	}

	return geoNames, nil
}

func main() {
	importCsvToSqlite("./", "GeoLite2-City-Locations-en.csv", "geonames.db")
}
