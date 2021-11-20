run:
	go run main.go

normal-test:
	echo "GET http://localhost:8080/normal" | vegeta attack -duration=1s -rate=10 | vegeta report

singleflight-test:
	echo "GET http://localhost:8080/singleflight" | vegeta attack -duration=1s -rate=10 | vegeta report
