.PHONY: clean

clean:
	find -type f -name '*~' -delete

.PHONY: test
test:
	go test 
	go test $$(basename $$PWD)/scrypt
