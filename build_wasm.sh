project_list="
dodger
easyui
fifteen
flappy
gocycle
infection
splatcard
vvv
"

for project in $project_list; do
  echo "Building: $project"
  
  pushd $project
  env GOOS=js GOARCH=wasm go build -o ../web/$project.wasm
  popd
done


