package user

type Handler struct {
    store *types.UserStore
}

func NewHandler(store *types.UserStore) *Handler {
    return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
    router.HandleFunc("/register", h.RegisterUser).Methods("POST")
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
    var payload types.RegisterUserRequest
    if err := utils.parseJson(r, &payload); err != nil {
        utils.WriteError(w, http.StatusBadRequest, err)
    }
    _, err := h.store.GetUserByEmail(payload.Email)
    if err == nil {
        utils.WriteError(w, http.StatusBadRequest, err)
    }
    hashedPassword, err := auth.HashPassword(payload.Password)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, err)
    }
    err = h.store.CreateUser(types.User{
        FirstName: payload.FirstName,
        LastName: payload.LastName,
        Email: payload.Email,
        Password: hashedPassword,
    })
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, err)
    }
    utils.WriteJSON(w, http.StatusOK, "User created successfully")
}
