Start-Process powershell -ArgumentList "cd server; go run .\cmd\server"
Start-Process powershell -ArgumentList "cd ui; npm run dev"