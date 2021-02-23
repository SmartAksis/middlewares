package parameters

var application map[string][]string = make(map[string][]string)

const SmartAkssisUserAdminKey = "MASTER_USER_SMARTAKSIS"

func GetMastersUsers() ([]string, *errorResultParameter) {
	users := application[SmartAkssisUserAdminKey]
	if users != nil {
		return users, nil
	}
	return nil, notFoutParameter()
}

func SetMastersUsers(users []string) {
	for _, user := range users {
		if contains(application[SmartAkssisUserAdminKey], user) == false{
			appended := append(application[SmartAkssisUserAdminKey], user)
			application[SmartAkssisUserAdminKey]=appended
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type errorResultParameter struct {
	Cause string
}

func (e errorResultParameter) Error() string {
	panic("implement me")
}

func notFoutParameter() *errorResultParameter {
	return &errorResultParameter{
		Cause: "not_found_parameter",
	}
}
