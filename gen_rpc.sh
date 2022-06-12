echo 'export GOPATH=$HOME/Go' >> $HOME/.bashrc source $HOME/.bashrc

# rm -rf protos/rpc/*

# protoc --go_out=. --go_opt=paths=source_relative ./protos/*
protoc --go_out=rpc --go_opt=paths=source_relative ./proto/inventory/inventory.proto

# --go-grpc_out=protos  --go-grpc_opt=paths=source_relative ./proto/inventory/*

echo "Copying into protos.."
cp -R rpc/proto/inventory/* rpc

echo "Cleaning protos directory.."
rm -rf rpc/proto

echo "Commtting code so that cahnges will reflect.."
git commit -am 'Updating submodules'

echo "DONE"