BIN_DIR = bin
PROTO_DIR = proto
SERVER_DIR = server
CLIENT_DIR = client

ifeq ($(OS), Windows_NT)
	SHELL := powershell.exe
	.SHELLFLAGS := -NoProfile -Command
	SHELL_VERSION = $(shell (Get-Host | Select-Object Version | Format-Table -HideTableHeaders | Out-String).Trim())
	OS = $(shell "{0} {1}" -f "windows", (Get-ComputerInfo -Property OsVersion, OsArchitecture | Format-Table -HideTableHeaders | Out-String).Trim())
	PACKAGE = $(shell (Get-Content go.mod -head 1).Split(" ")[1])
	CHECK_DIR_CMD = if (!(Test-Path $@)) { $$e = [char]27; Write-Error "$$e[31mDirectory $@ doesn't exist$${e}[0m" }
	HELP_CMD = Select-String "^[a-zA-Z_-]+:.*?\#\# .*$$" "./Makefile" | Foreach-Object { $$_data = $$_.matches -split ":.*?\#\# "; $$obj = New-Object PSCustomObject; Add-Member -InputObject $$obj -NotePropertyName ('Command') -NotePropertyValue $$_data[0]; Add-Member -InputObject $$obj -NotePropertyName ('Description') -NotePropertyValue $$_data[1]; $$obj } | Format-Table -HideTableHeaders @{Expression={ $$e = [char]27; "$$e[36m$$($$_.Command)$${e}[0m" }}, Description
	RM_F_CMD = Remove-Item -erroraction silentlycontinue -Force
	RM_RF_CMD = ${RM_F_CMD} -Recurse
	SERVER_BIN = ${SERVER_DIR}.exe
	CLIENT_BIN = ${CLIENT_DIR}.exe
else
	SHELL := bash
	SHELL_VERSION = $(shell echo $$BASH_VERSION)
	UNAME := $(shell uname -s)
	VERSION_AND_ARCH = $(shell uname -rm)
	ifeq ($(UNAME),Darwin)
		OS = macos ${VERSION_AND_ARCH}
	else ifeq ($(UNAME),Linux)
		OS = linux ${VERSION_AND_ARCH}
	else
		$(error OS not supported by this Makefile)
	endif
	PACKAGE = $(shell head -1 go.mod | awk '{print $$2}')
	CHECK_DIR_CMD = test -d $@ || (echo "\033[31mDirectory $@ doesn't exist\033[0m" && false)
	HELP_CMD = grep -E '^[a-zA-Z_-]+:.*?\#\# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?\#\# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	RM_F_CMD = rm -f
	RM_RF_CMD = ${RM_F_CMD} -r
	SERVER_BIN = ${SERVER_DIR}
	CLIENT_BIN = ${CLIENT_DIR}
endif

.DEFAULT_GOAL := help
.PHONY: students accounts mentors donators admin blog help
project := students accounts mentors donators admin blog

all: $(project) ## Generate Pbs and build

students: $@ ## Generate Pbs and build for students
accounts: $@ ## Generate Pbs and build for accounts
mentors: $@ ## Generate Pbs and build for accounts
donators: $@ ## Generate Pbs and build for donators
admin: $@ ## Generate Pbs and build for admin
blog: $@ ## Generate Pbs and build for blog

$(project):
	@${CHECK_DIR_CMD}
	protoc -I$@/${PROTO_DIR} --go_opt=module=${PACKAGE} --go_out=. --go-grpc_opt=module=${PACKAGE} --go-grpc_out=. $@/${PROTO_DIR}/*.proto
	go build -o ${BIN_DIR}/$@/${SERVER_BIN} ./$@/${SERVER_DIR}
	go build -o ${BIN_DIR}/$@/${CLIENT_BIN} ./$@/${CLIENT_DIR}

test: all ## Launch tests
	go test ./...

clean: clean_students clean_accounts clean_mentors clean_donators clean_admin clean_blog ## Clean generated files
	${RM_F_CMD} ssl/*.crt
	${RM_F_CMD} ssl/*.csr
	${RM_F_CMD} ssl/*.key
	${RM_F_CMD} ssl/*.pem
	${RM_RF_CMD} ${BIN_DIR}

clean_students: ## Clean generated files for students
	${RM_F_CMD} students/${PROTO_DIR}/*.pb.go

clean_accounts: ## Clean generated files for accounts
	${RM_F_CMD} accounts/${PROTO_DIR}/*.pb.go

clean_mentors: ## Clean generated files for mentors
	${RM_F_CMD} mentors/${PROTO_DIR}/*.pb.go

clean_donators: ## Clean generated files for donators
	${RM_F_CMD} donators/${PROTO_DIR}/*.pb.go

clean_admin: ## Clean generated files for accounts
	${RM_F_CMD} admin/${PROTO_DIR}/*.pb.go

clean_blog: ## Clean generated files for blog
	${RM_F_CMD} blog/${PROTO_DIR}/*.pb.go

rebuild: clean all ## Rebuild the whole project

bump: all ## Update packages version
	go get -u ./...

about: ## Display info related to the build
	@echo "OS: ${OS}"
	@echo "Shell: ${SHELL} ${SHELL_VERSION}"
	@echo "Protoc version: $(shell protoc --version)"
	@echo "Go version: $(shell go version)"
	@echo "Go package: ${PACKAGE}"
	@echo "Openssl version: $(shell openssl version)"

help: ## Show this help
	@${HELP_CMD}