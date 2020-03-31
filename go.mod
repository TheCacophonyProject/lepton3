module github.com/TheCacophonyProject/lepton3

go 1.12

require (
	github.com/TheCacophonyProject/go-cptv v0.0.0-20200116020937-858bd8b71512
	github.com/alexflint/go-arg v0.0.0-20180516182405-f7c0423bd11e
	github.com/alexflint/go-scalar v1.0.0 // indirect
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637
	periph.io/x/periph v3.6.2+incompatible
)

// We maintain a custom fork of periph.io at the moment.
replace periph.io/x/periph => github.com/TheCacophonyProject/periph v1.0.1-0.20200331204442-4717ddfb6980
