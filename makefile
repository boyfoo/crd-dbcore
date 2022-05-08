copyToGoPath:
	rm -rf ${GOPATH}/src/github.com/shenyisyn/dbcore
	MKDIR -p ${GOPATH}/src/github.com/shenyisyn/dbcore/pkg
	cp -r ./pkg ${GOPATH}/src/github.com/shenyisyn/dbcore
	cp -r ./main.go ${GOPATH}/src/github.com/shenyisyn/dbcore
	cp -r ./go.mod ${GOPATH}/src/github.com/shenyisyn/dbcore
	rm -rf ${GOPATH}/src/github.com/shenyisyn/dbcore/pkg/apis/dbconfig/v1/zz_generated.deepcopy.go
	rm -rf ${GOPATH}/src/github.com/shenyisyn/dbcore/pkg/client
	cd ${GOPATH}/src/github.com/shenyisyn/dbcore && ${GOPATH}/src/k8s.io/code-generator/generate-groups.sh all github.com/shenyisyn/dbcore/pkg/client github.com/shenyisyn/dbcore/pkg/apis dbconfig:v1
	cp ${GOPATH}/src/github.com/shenyisyn/dbcore/pkg/apis/dbconfig/v1/zz_generated.deepcopy.go ./pkg/apis/dbconfig/v1/zz_generated.deepcopy.go
	cp -r ${GOPATH}/src/github.com/shenyisyn/dbcore/pkg/client ./pkg/

