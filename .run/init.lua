if run.argv[1] == nil then
  run.argv[1] = "build"
end
if run.argv[1] == "build" then
  run.env["RUNCTR"] = "./etc/Dockerfile.dev"
end
run.env["NAME"] = "defers"
