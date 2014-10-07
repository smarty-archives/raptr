package manifest

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReaderWriter(t *testing.T) {
	Convey("When a Debian is read and then written", t, func() {
		Convey("The hash should be the same.", nil)
	})

	handle, err := os.Open("debian/control")
	if err != nil {
		panic(err)
	}
	defer handle.Close()

	// payload, _ := ioutil.ReadAll(handle)
	// handle.Seek(0, 0)
	// computed := md5.Sum(payload)
	// fmt.Println("MD5:", hex.EncodeToString(computed[:]))

	// reader := manifest.NewReader(handle)

	// buffer := bytes.NewBuffer([]byte{})
	// writer := manifest.NewWriter(buffer)

	// for {
	// 	if item, err := reader.Read(); err == io.EOF {
	// 		break
	// 	} else if err != nil {
	// 		panic(err)
	// 	} else {
	// 		fmt.Println(item)
	// 		writer.Write(item)
	// 	}
	// }

	// fmt.Println("--------------------------------")
	// written := string(buffer.Bytes())
	// fmt.Println(written)
	// computed1 := md5.Sum(buffer.Bytes())
	// fmt.Println("MD5:", hex.EncodeToString(computed1[:]))

}

const debianFile = `Format: 3.0 (quilt)
Source: nginx
Binary: nginx, nginx-doc, nginx-common, nginx-full, nginx-full-dbg, nginx-light, nginx-light-dbg, nginx-extras, nginx-extras-dbg, nginx-naxsi, nginx-naxsi-dbg, nginx-naxsi-ui
Architecture: any all
Version: 1.7.4-1
Maintainer: Kartik Mistry <kartik@debian.org>
Uploaders: Jose Parrella <bureado@debian.org>, Fabio Tranchitella <kobold@debian.org>, Michael Lustfield <michael@lustfield.net>, Dmitry E. Oboukhov <unera@debian.org>, Cyril Lavier <cyril.lavier@davromaniak.eu>, Christos Trochalakis <yatiohi@ideopolis.gr>
Homepage: http://nginx.net
Standards-Version: 3.9.5
Vcs-Browser: http://anonscm.debian.org/gitweb/?p=collab-maint/nginx.git;a=summary
Vcs-Git: git://anonscm.debian.org/collab-maint/nginx.git
Build-Depends: autotools-dev, debhelper (>= 9), dh-systemd (>= 1.5), dpkg-dev (>= 1.15.5), libexpat-dev, libgd2-dev | libgd2-noxpm-dev, libgeoip-dev, liblua5.1-dev, libmhash-dev, libpam0g-dev, libpcre3-dev, libperl-dev, libssl-dev, libxslt1-dev, po-debconf, zlib1g-dev
Package-List: 
 nginx deb httpd optional
 nginx-common deb httpd optional
 nginx-doc deb doc optional
 nginx-extras deb httpd extra
 nginx-extras-dbg deb debug extra
 nginx-full deb httpd optional
 nginx-full-dbg deb debug extra
 nginx-light deb httpd extra
 nginx-light-dbg deb debug extra
 nginx-naxsi deb httpd extra
 nginx-naxsi-dbg deb debug extra
 nginx-naxsi-ui deb httpd extra
Checksums-Sha1: 
 94f4ac8ddb4a05349e75c43b84f24dbacdbac6e9 817174 nginx_1.7.4.orig.tar.gz
 728f8d1c6441dd7b095c1db7be338b778af53216 1568360 nginx_1.7.4-1.debian.tar.gz
Checksums-Sha256: 
 935c5a5f35d8691d73d3477db2f936b2bbd3ee73668702af3f61b810587fbf68 817174 nginx_1.7.4.orig.tar.gz
 c836bfe0ed55ef1bcf5f51d594eef27e2a1a809fe6ba856452c6c4cb8d6a4e8b 1568360 nginx_1.7.4-1.debian.tar.gz
Files: 
 bfc256cf72123601af56501b0a6a41f5 817174 nginx_1.7.4.orig.tar.gz
 b14cd918daa408d9c4ea54734dfef26a 1568360 nginx_1.7.4-1.debian.tar.gz
`
