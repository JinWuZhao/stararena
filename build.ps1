.\clean.ps1
New-Item -ItemType "directory" -Path "output"
go build -v -o "output\stararena.exe" .
New-Item -ItemType "directory" -Path "output\conf"
Copy-Item "conf\conf.toml" -Destination "output\conf"
New-Item -ItemType "directory" -Path "output\sc2maps\product"
Copy-Item -Path "sc2maps\product\*" -Destination "output\sc2maps\product"
Compress-Archive -Path "output" -DestinationPath ".\output.zip"