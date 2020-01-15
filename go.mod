module github.com/TheCacophonyProject/lepton3

go 1.12

require (
	github.com/TheCacophonyProject/go-cptv v0.0.0-20190829021603-992765fceb9b
	github.com/alexflint/go-arg v0.0.0-20180516182405-f7c0423bd11e
	github.com/alexflint/go-scalar v0.0.0-20170216020425-e80c3b7ed292 // indirect
	github.com/gferraro/go-cptv v0.0.0-20190522012101-cfce3277dd0d // indirect
	github.com/stretchr/testify v1.2.2
	golang.org/x/net v0.0.0-20180811021610-c39426892332 // indirect
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637
	periph.io/x/periph v3.6.2+incompatible
)

// We maintain a custom fork of periph.io at the moment.
replace periph.io/x/periph => github.com/TheCacophonyProject/periph v2.0.1-0.20171123021141-d06ef89e37e8+incompatible

replace github.com/TheCacophonyProject/go-cptv => github.com/gferraro/go-cptv v0.0.0-20200115004624-7672c7ff464b
