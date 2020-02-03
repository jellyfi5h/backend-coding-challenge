package main

import (
	"./server"
	"./settings"
)

var content string = `
flask=ecwcwec
Flask-erver=rvev
beautifulsoup4==4.8.2
bs4==0.0.1
certifi==2019.11.28
djangod=wecewc
Django-fbd=wcwe
chardet==3.0.4
Click==7.0

	FlaSk==1.1.1
Flask-Caching==1.8.0
Flaskgbg-Limiter==1.1.0






flasksdv==3.20.9
gevent==1.4.0
greenlet==0.4.15
httplib2==0.16.0
idna==2.8
itsdangerous==1.1.0
Jinja2==2.10.3
limits==1.5
MarkupSafe==1.1.1
pep8==1.7.1
pycodestyle==2.5.0
requests==2.22.0
six==1.14.0
soupsieve==1.9.5
urllib3==1.25.8
Werkzeug==0.16.0
	`

func main() {
	usage := settings.GetInfo() //get the imformation ncessair(host, post) for run the server
	if usage != nil {
		//print usage
		return
	}
	err := server.CreateDaemon()
	if err != nil {
		return //something goes wrong with creating the daemon process
	}
	server.WebServer() //start listen from host:port given
}
