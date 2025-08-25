.PHONY:mock
mock:
	@mockgen -source="E:\Book_Exp\webook\internal\service\user.go" -package=svcmocks -destination="E:\Book_Exp\webook\internal\service\mocks\user.mock.go"
	@mockgen -source="E:\Book_Exp\webook\internal\service\user.go" -package=svcmocks -destination="E:\Book_Exp\webook\internal\service\mocks\user.mock.go"
	@mockgen -source="E:\Book_Exp\webook\internal\repository\user.go" -package=repomocks -destination="E:\Book_Exp\webook\internal\repository\mocks\user.mock.go"
	@mockgen -source="E:\Book_Exp\webook\internal\repository\code.go" -package=repomocks -destination="E:\Book_Exp\webook\internal\repository\mocks\code.mock.go"
	@mockgen -source="E:\Book_Exp\webook\internal\repository\dao\user.go" -package=daomocks -destination="E:\Book_Exp\webook\internal\repository\dao\mocks\user.mock.go"
	@mockgen -source="E:\Book_Exp\webook\internal\repository\cache\user.go" -package=cachemocks -destination="E:\Book_Exp\webook\internal\repository\cache\mocks\user.mock.go"
	@mockgen -package=redismocks -destination="E:\Book_Exp\webook\internal\repository\cache\redismocks\cmdable.mock.go" github.com/redis/go-redis/v9 Cmdable
	@mockgen -source="E:\Book_Exp\webook\pkg\ginx\ratelimit\types.go" -package=limitmocks -destination="E:\Book_Exp\webook\pkg\ginx\ratelimit\mocks\ratelimit.mock.go"
	@mockgen -source="E:\Book_Exp\webook\internal\service\sms\types.go" -package=failovermocks -destination="E:\Book_Exp\webook\internal\service\sms\failover\mocks\serivce.mock.go"
	@go mod tidy



