package idcard_cn


func MakeIDCard(uniqueIDCode string, check bool) (IDCard, error) {
	idcard := &IDParser{}
	return idcard, idcard.Init(uniqueIDCode, check)
}