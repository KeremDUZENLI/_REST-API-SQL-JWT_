package models

import (
	"os"
)

var TABLE string = "table"
var SECRET_KEY string = os.Getenv("SECRET_KEY")
