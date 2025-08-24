.PHONY:mock
mock:
	@mockgen -source="E:\Book_Exp\webook\internal\service\user.go" -package=svcmocks -destination="E:\Book_Exp\webook\internal\service\mocks\user.mock.go"
	@mockgen -source="E:\Book_Exp\webook\internal\service\user.go" -package=svcmocks -destination="E:\Book_Exp\webook\internal\service\mocks\user.mock.go"
	@mockgen -source="E:\Book_Exp\webook\internal\respository\user.go" -package=repomocks -destination="E:\Book_Exp\webook\internal\respository\mocks\user.mock.go"
	@mockgen -source="E:\Book_Exp\webook\internal\respository\code.go" -package=repomocks -destination="E:\Book_Exp\webook\internal\respository\mocks\code.mock.go"

	@go mod tidy



