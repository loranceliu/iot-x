package types

type Device struct {
	DeviceId int64
	Sn       string
	Model    string
	Instance Instance
}
