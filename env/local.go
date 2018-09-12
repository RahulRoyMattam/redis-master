package env

import "os"

func setLocalVariables() {
	os.Setenv("REDIS_URL", "<redis url:port>")

	//config Vars
	os.Setenv("REDIS_MASTER_ROLE", "SLAVE")
	os.Setenv("APP_MAX_REDIS_CONN", "2")
}
