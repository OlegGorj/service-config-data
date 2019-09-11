package kernel


type Kernel struct {
	Name string
	JSON []byte
}

//// TODO: Finish this
//func GetKernels(environment *conf.Environment) (int, error){
//	//listOfUsers := environment.Users
//	//for i := range listOfUsers{
//	//	if listOfUsers[i].Email == email {
//	//		return i,nil
//	//	}
//	//}
//	return -1, fmt.Errorf("user not found")
//}