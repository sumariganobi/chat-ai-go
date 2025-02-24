##Dokumentasi Buat API dengan GOLANG

Install release suitable for your Operating System
https://go.dev/dl/

open command or powershell

```bash
  go version
  #go version go1.23.4 windows/amd64
  go get -u github.com/gin-gonic/gin
  go mod init name_project
  go mod tidy

  go run main.go
  #[GIN-debug] Listening and serving HTTP on :8080
```

##Running on Centos
```bash
  wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
  sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
  sudo nano ~/.bashrc
  export PATH=$PATH:/usr/local/go/bin
  source ~/.bashrc

  go version
  #go version go1.22.0 linux/amd64

  go run main.go --HOST 10.101.34.193:3000
  #running dibackground
  go run main.go & 
  #proses tetap berjalan meskipun terminal ditutup
  nohup go run main.go > output.log 2>&1 &

  #memeriksa log
  tail -f output.log

  #memastika proses berjalan
  ps -aux | grep main
```

##Menggunakan systemd
```bash
  #buat file service di /etc/systemd/system
  sudo nano /etc/systemd/system/name_app.service
```

##Sesuaikan Config Service
```bash
  [Unit]
  Description=My Go Application
  After=network.target

  [Service]
  ExecStart=/usr/local/go/bin/go run /path/to/your/main.go
  Restart=always
  User=root
  Environment=PATH=/usr/bin:/usr/local/go/bin
  Environment=GOPATH=/root/go
  WorkingDirectory=/path/to/your/app

  [Install]
  WantedBy=multi-user.target
```

```bash
  sudo systemctl daemon-reload
  sudo systemctl start myapp.service
  sudo systemctl enable myapp.service
```

```bash
  #check status service 
  sudo systemctl status myapp.service
```