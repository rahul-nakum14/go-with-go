package main

func main(){
	db,err := db.NewSQLStorage(mysql.config{
		Host: config.Envs.DBHost,
		Port: config.Envs.DBPort,
		User: config.Envs.DBUser,
		Password: config.Envs.DBPassword,
		Name: config.Envs.DBName,
		SSLMode: config.Envs.DBSSLMode,
	})
	initStorage(db)
    server := newApiServer("localhost:8080", db)
    if err := server.run(); err != nil {
		log.Fatal("failed to start server: %v", err)
	}
}


func initStorage() (*sql.DB, error) {
	err := db.Ping(); err !=nil {
		log.Fatal("failed to connect to database: %v", err)
	}
	return db, nil
}