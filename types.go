package vsz_web_backend

type Opdrachtgever struct {
	Bedrijfscode int
	Bedrijfsnaam string
	Email        string
	Wachtwoord   []byte
}

type Kruising struct {
	Kruisingscode    int
	Plaats           string
	Latitude         float64
	Longitude        float64
	Weg              string
	Bedrijfscode     int
	Laatst_Opgestart string
	Plaatsing        string
}

type Auto struct {
	Autocode      int
	DatumTijd     string
	Richting      string
	Kruisingscode int
}

type Session struct {
	User int `json:"user"`
}
