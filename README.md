# Geração das mensagens e serviços.

## Versão 1

Para isso, instale a última versão do grpc para o Go

    go get -u google.golang.org/protobuf/cmd/protoc-gen-go
    go install google.golang.org/protobuf/cmd/protoc-gen-go

    go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

Em seguida rode o comando a seguir à partir da raíz do projeto para gerar as classes de apoio do grpc

    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative usermgmt/usermgmt.proto

Instalando as demais dependências:

    go mod tidy

Para executar, em um teerminal execute:

    go run usermgmt_server/usermgmt_server.go

Você deverá receber a seguinte resposta:

    Server listening at [::]:50051

Em um novo terminal, na raíz do projeto, digite:

    go run usermgmt_client/usermgmt_client.go

Você obterá a seguinte resposta:

    User Details: [Name: User 1, Age: 38, Id: 408]
    User Details: [Name: User 2, Age: 8, Id: 387]
    User Details: [Name: User 3, Age: 29, Id: 831]
    User Details: [Name: User 4, Age: 66, Id: 429]

## Versão 2

Nessa segunda versão iremos persistir os dados.

Em um terminal, rode o comando a seguir à partir da raíz do projeto para atualizar as classes de apoio do grpc

    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative usermgmt/usermgmt.proto

