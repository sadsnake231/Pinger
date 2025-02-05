package API


var conn *pgxpool.Pool = database.EstablishConnection()

func GetPings() gin.HandlerFunc {
	return func (c *gin.Context) {
		
	}
}

func PostPings() gin.HandlerFunc {
	return func (c *gin.Context) {
		
	}
}
