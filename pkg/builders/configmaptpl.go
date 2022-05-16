package builders

const cmtpl = `
dbConfig:
 dsn: "[[ .Dsn ]]"
 maxOpenConn: [[ .MaxOpenConn ]]
 maxLifeTime: [[ .MaxLifeTime ]]
 maxIdleConn: [[ .MaxIdleConn ]]
appConfig:
 rpcPort: 8081
 httpPort: 8090
apis:
 - name: test
   sql: "select * from test"
`
