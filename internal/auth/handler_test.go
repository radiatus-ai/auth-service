package auth

type mockService struct {
	registerFunc    func(email, password string) error
	loginFunc       func(email, password string) (string, error)
	verifyTokenFunc func(token string) (string, error)
}

func (m *mockService) Register(email, password string) error {
	return m.registerFunc(email, password)
}

func (m *mockService) Login(email, password string) (string, error) {
	return m.loginFunc(email, password)
}

func (m *mockService) VerifyToken(token string) (string, error) {
	return m.verifyTokenFunc(token)
}
