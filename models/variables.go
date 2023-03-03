package models

import (
	"os"
)

const EMPTY = ""

const ADMIN string = "ADMIN"
const USER string = "USER"

const TABLE string = "table"

var SECRET_KEY string = os.Getenv("SECRET_KEY")
