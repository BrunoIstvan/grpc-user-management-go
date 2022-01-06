# Geração das mensagens e serviços.

Para isso, instale a última versão do grpc para o Go

    go get -u google.golang.org/protobuf/cmd/protoc-gen-go
    go install google.golang.org/protobuf/cmd/protoc-gen-go

    go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

Em seguida rode o comando a seguir à partir da raíz do projeto

    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative usermgmt/usermgmt.proto

Instalando as demais dependências

    go mod tidy

    
