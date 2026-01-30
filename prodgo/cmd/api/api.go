package api

type apiServer struct {
    addr string
    db *sql.DB
}

func newApiServer(addr string, db *sql.DB) *apiServer {
    return &apiServer{
        addr: addr,
        db: db,
    }
}

func (s *apiServer) run() error {
    router := mux.NewRouter();
    subrouter := router.PathPrefix("/api").Subrouter()
    userStore := user.NewStore(s.db)
    userHandler := user.NewHandler(userStore)
    userHandler.RegisterRoutes(subrouter)
    return http.ListenAndServe(s.addr, s)
}
