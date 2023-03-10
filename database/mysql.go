package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
	vsz_web_backend "vsz-web-backend"
	"vsz-web-backend/config"
)

var db *sql.DB

func Initialize() error {
	cfg := config.Global.Mysql

	var err error
	// Setup mysql connection
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", cfg.User, cfg.Pass, cfg.Host, cfg.DB))
	if err != nil {
		return fmt.Errorf("failed to start connection: %s", err)
	}

	// Configure mysql connection
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// Test mysql connection
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to test connection: %s", err)
	}

	return nil
}

func GetOpdrachtgeverByID(id int) (*vsz_web_backend.Opdrachtgever, error) {
	res, err := db.Query("SELECT * FROM Opdrachtgever WHERE bedrijfscode = ?", id)
	if err != nil {
		return nil, err
	}

	defer res.Close()

	if res.Next() {
		var opdrachtgever vsz_web_backend.Opdrachtgever
		err := res.Scan(&opdrachtgever.Bedrijfscode, &opdrachtgever.Bedrijfsnaam, &opdrachtgever.Email, &opdrachtgever.Wachtwoord)
		if err != nil {
			return nil, err
		}
		return &opdrachtgever, nil
	} else {
		return nil, nil
	}
}

func GetOpdrachtgeverByEmail(email string) (*vsz_web_backend.Opdrachtgever, error) {
	res, err := db.Query("SELECT * FROM Opdrachtgever WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	defer res.Close()

	if res.Next() {
		var opdrachtgever vsz_web_backend.Opdrachtgever
		err := res.Scan(&opdrachtgever.Bedrijfscode, &opdrachtgever.Bedrijfsnaam, &opdrachtgever.Email, &opdrachtgever.Wachtwoord)
		if err != nil {
			return nil, err
		}
		return &opdrachtgever, nil
	} else {
		return nil, nil
	}
}

func GetOpdrachtgevers() ([]vsz_web_backend.Opdrachtgever, error) {
	res, err := db.Query("SELECT * FROM Opdrachtgever")
	if err != nil {
		return nil, err
	}

	defer res.Close()

	var opdrachtgevers []vsz_web_backend.Opdrachtgever
	for res.Next() {
		var opdrachtgever vsz_web_backend.Opdrachtgever
		err := res.Scan(&opdrachtgever.Bedrijfscode, &opdrachtgever.Bedrijfsnaam, &opdrachtgever.Email, &opdrachtgever.Wachtwoord)
		if err != nil {
			return nil, err
		}
		opdrachtgevers = append(opdrachtgevers, opdrachtgever)
	}

	return opdrachtgevers, nil
}

func GetOpdrachtgeverCount() (int, error) {
	var aantal int

	res, err := db.Query("SELECT COUNT(*) FROM Opdrachtgever")
	if err != nil {
		return aantal, err
	}

	defer res.Close()

	if res.Next() {
		err := res.Scan(&aantal)
		if err != nil {
			return aantal, err
		}
	} else {
		return aantal, errors.New("no results")
	}

	return aantal, nil
}

func GetKruisingen() ([]vsz_web_backend.Kruising, error) {
	res, err := db.Query("SELECT kruisingscode,bedrijfsnaam,plaats,weg,plaatsing,laatst_opgestart FROM Kruising JOIN Opdrachtgever ON Kruising.bedrijfscode=Opdrachtgever.bedrijfscode")
	if err != nil {
		return nil, err
	}

	defer res.Close()

	var kruisingen []vsz_web_backend.Kruising
	for res.Next() {
		var kruising vsz_web_backend.Kruising
		err := res.Scan(&kruising.Kruisingscode, &kruising.Bedrijfsnaam, &kruising.Plaats, &kruising.Weg, &kruising.Plaatsing, &kruising.Laatst_Opgestart)
		if err != nil {
			return nil, err
		}
		kruisingen = append(kruisingen, kruising)
	}

	return kruisingen, nil
}

func GetAutos() ([]vsz_web_backend.Auto, error) {
	res, err := db.Query("SELECT autocode, DATE(datumtijd), TIME(datumtijd), richting, weg FROM Auto JOIN Kruising ON Auto.kruisingscode=Kruising.kruisingscode ORDER BY datumtijd")
	if err != nil {
		return nil, err
	}

	defer res.Close()

	var autos []vsz_web_backend.Auto
	for res.Next() {
		var auto vsz_web_backend.Auto
		err := res.Scan(&auto.Autocode, &auto.Datum, &auto.Tijd, &auto.Richting, &auto.Weg)
		if err != nil {
			return nil, err
		}
		autos = append(autos, auto)
	}

	return autos, nil
}

func GetAutosWeek() ([]int, error) {
	res, err := db.Query("SELECT COUNT(*) FROM Auto GROUP BY DAY(datumtijd) ORDER BY DAY(datumtijd) DESC LIMIT 7")
	if err != nil {
		return nil, err
	}

	defer res.Close()

	var autos []int
	for res.Next() {
		var aantal int
		err := res.Scan(&aantal)
		if err != nil {
			return nil, err
		}
		autos = append(autos, aantal)
	}

	return autos, nil
}

func GetAutosMaand() (int, error) {
	var autos int

	res, err := db.Query("SELECT COUNT(*) FROM Auto WHERE MONTH(datumtijd) = MONTH(CURDATE())")
	if err != nil {
		return autos, err
	}

	defer res.Close()

	if res.Next() {
		err := res.Scan(&autos)
		if err != nil {
			return autos, err
		}
	} else {
		return autos, errors.New("no results")
	}

	return autos, nil
}

func GetAutosKruising() (map[string]float64, error) {
	res, err := db.Query("SELECT weg, COUNT(*)*100/t.s AS percentage FROM Auto JOIN Kruising ON Auto.kruisingscode=Kruising.kruisingscode CROSS JOIN (SELECT COUNT(*) AS s FROM Auto) t GROUP BY weg;")
	if err != nil {
		return nil, err
	}

	defer res.Close()

	var autos = make(map[string]float64)
	for res.Next() {
		var percentage float64
		var weg string
		err := res.Scan(&weg, &percentage)
		if err != nil {
			return nil, err
		}
		autos[weg] = percentage
	}

	return autos, nil
}

func GetAutosOpdrachtgever() (map[string]float64, error) {
	res, err := db.Query("SELECT bedrijfsnaam, COUNT(*)*100/t.s AS percentage FROM Auto JOIN Kruising ON Auto.kruisingscode=Kruising.kruisingscode JOIN Opdrachtgever ON Kruising.bedrijfscode=Opdrachtgever.bedrijfscode CROSS JOIN (SELECT COUNT(*) AS s FROM Auto) t GROUP BY bedrijfsnaam")
	if err != nil {
		return nil, err
	}

	defer res.Close()

	var autos = make(map[string]float64)
	for res.Next() {
		var percentage float64
		var bedrijfsnaam string
		err := res.Scan(&bedrijfsnaam, &percentage)
		if err != nil {
			return nil, err
		}
		autos[bedrijfsnaam] = percentage
	}

	return autos, nil
}
