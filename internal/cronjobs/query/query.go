package query
import "go-heartbeat/pkg/mysql"

type QueryServerId struct {
	Value int
}

func GetServerId(con *mysql.DBModel) (int, error) {
	var serverid int
	query := 'SHOW VARIABLES LIKE "server_id"'

	row := con.DBEngine.QueryRow()


}
