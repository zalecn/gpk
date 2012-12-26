


clean:
	rm -Rf target/*

prepare:
	mkdir -p target/pkg


footprint: prepare
	cp ./bin/gpk target/pkg
	rst2man doc/gpk.rst target/pkg/gpk.1

debian: footprint
	cp dist/deb/gopack target/pkg
	cd target/pkg && equivs-build gopack

tarball: footprint
	tar czvf target/gopack-1.0.0.tar.gz src/ doc/ .gpk
