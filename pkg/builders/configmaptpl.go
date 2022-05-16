package builders

const cmtpl = `
dbConfig:
 dsn: "[[ .Dsn ]]"
 maxOpenConn: 20
 maxLifeTime: 1800
 maxIdleConn: 5
appConfig:
 rpcPort: 8081
 httpPort: 8090
apis:
 - name: test
   sql: "select * from test"
`
