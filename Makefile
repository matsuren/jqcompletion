.PHONY: testjsonview
testjsonview:
	go test ./jsonview --tags debug -v -count=1

.PHONY: testqueryviewonly
testqueryviewonly:
	DEBUGLOG=1 go test ./queryview --tags debug -v -count=1 -run OnlyView

.PHONY: testqueryviewquery
testqueryviewquery:
	DEBUGLOG=1 go test ./queryview --tags debug -v -count=1 -run Query
