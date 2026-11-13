package handler

type MockValidatorManager struct {
	RegisterValidatorsError error
}

func (m *MockValidatorManager) RegisterValidators() error {
	return m.RegisterValidatorsError
}
