module github.com/TheCacophonyProject/lepton3

go 1.12

require (
	github.com/TheCacophonyProject/go-cptv v0.0.0-20200616224711-fc633122087a
	github.com/alexflint/go-arg v0.0.0-20180516182405-f7c0423bd11e
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637
	periph.io/x/periph v3.6.2+incompatible
)

// We maintain a custom fork of periph.io at the moment.
replace periph.io/x/periph => github.com/TheCacophonyProject/periph v2.1.1-0.20200615222341-6834cd5be8c1+incompatible
