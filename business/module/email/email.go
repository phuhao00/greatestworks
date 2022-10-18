package email

type Mail interface {
	ToPB()
	LoadFrom()
	GetDBModel()
}