package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	revCount, err := exec.Command("git", "rev-list", "HEAD", "--count").Output()
	checkError(err)
	commitHash, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	checkError(err)
	output := fmt.Sprintf(`# Maintainer: Louis Dalibard <louis.dalibard@gmail.com>
_pkgbase="pscreen"
pkgname="$_pkgbase-git"
pkgver=r%s.%s
pkgrel=1
pkgdesc="A companion app for the pScreen desktop clock"
arch=("x86_64" "armv7h" "aarch64")
url="https://github.com/make-42/pscreen"
license=('GPL')
groups=()
depends=()
makedepends=('go>=1.18')
optdepends=()
provides=("$_pkgbase")
conflicts=("$_pkgbase")
replaces=()
backup=()
options=()
install=
changelog=
source=("git+https://github.com/make-42/pscreen.git")
noextract=()
md5sums=("SKIP") #autofill using updpkgsums

pkgver() {
	cd ${srcdir}/${_pkgbase}
	echo "r$(git rev-list --count HEAD).g$(git rev-parse --short HEAD)"
}

build() {
	cd "${srcdir}/${_pkgbase}/Software/Companion App"
	go mod tidy
	go build
}
	
package() {
	install -Dm755 "${srcdir}/${_pkgbase}/Software/Companion App/${_pkgbase}" "${pkgdir}"/usr/bin/${_pkgbase}
}`, revCount[:len(revCount)-1], commitHash[:len(commitHash)-1])
	err = os.WriteFile("PKGBUILD", []byte(output), 0644)
	checkError(err)
}
