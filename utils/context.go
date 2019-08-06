package utils

var allowdApps  = map[string]string{"8080bb598ae1b6af3452c69371946d06":"sspd_loop","ssss":"sdf"}


func VolidateAPP(appMd5 string) bool  {
	_,ok := allowdApps[appMd5]
	if ok {
		return true
	}

	return true
}