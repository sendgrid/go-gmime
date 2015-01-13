# Intro

This document describe reproducing build of the packages with recent gmime and
glib2 for CentOS 6 (EL6) and Ubuntu 12.04 Precise.

I used `sbuild` installed on server with ``Ubuntu Precise`, and docker with 
CentOS 6.4 on same machine. Both of them -- `amd64` aka `x86_64`.


## Ubuntu

### Prepare repository

```sh
apt-get install reprepro
mkdir -p /srv/www/packages/ubuntu/conf
cd /srv/www/packages/ubuntu
```

Create file `/srv/www/packages/ubuntu/conf/distributions`:

    Origin: SendGrid Labs
    Suite: precise
    Codename: precise
    Label: SendGrind labs repo
    Architectures: i386 amd64 source
    Components: main contrib non-free
    UDebComponents: main contrib non-free
    Description: labs repo
    SignWith:  maint@example.com

Create file `/srv/www/packages/ubuntu/conf/options`:

    gnupghome /srv/www/packages/ubuntu/keys

Generate keypair:

    gpg --homedir /srv/www/packages/ubuntu/keys --gen-key
    gpg --homedir /srv/www/packages/ubuntu/keys --export --armor >/srv/www/packages/ubuntu/key.asc
    chown root:sbuild /srv/www/packages/ubuntu/keys
    chmod 770 /srv/www/packages/ubuntu/keys


Configure webserver (nginx):
(this is example only, you can use any webserver, able to serve static files)

    server {
        listen 80;
        server_name build.example.com;
        location /ubuntu  {
            autoindex on;
            root /srv/www/packages;
        }
    }

### Prepare sbuild

```sh
apt-get install sbuild
mkdir /srv/packagebuilder
cd /srv/packagebuilder
sbuild-createchroot --arch=amd64 precise /srv/packagebuilder/precise-amd64 http://archive.ubuntu.com/ubuntu
sbuild-createchroot --arch=i386 precise /srv/packagebuilder/precise-i386 http://archive.ubuntu.com/ubuntu
```

Add line:

    deb http://build.example.com/ubuntu precise main non-free contrib

to both `/srv/packagebuilder/precise-i386/etc/apt/sources.list` and
`/srv/packagebuilder/precise-amd64/etc/apt/sources.list`

Add key.asc to chroots:

    cp /srv/www/packages/ubuntu/key.asc /srv/packagebuilder/precise-amd64/
    schroot-shell precise-amd64

Then inside this shell isssue command:

    apt-key add key.asc

Then exit from shell, and repeat for i386 chroot.
Last thing needed to be done, is generating sbuild own keypair.

    sbuild-update --keygen


### Build prerequisites

We need to build `libpcre3`, `highlight`, and `gtk-doc` packages.
(where urls to source packages not specified -- packages was obtained from
`debian/sid` via `apt-get source` command)

    sbuild --append-to-version=".labs1" -m "Alexander V. Nikolaev <avn@daemon.hole.ru>" -d precise pcre3_8.31-3.dsc
    sbuild --append-to-version=".labs1" -m "Alexander V. Nikolaev <avn@daemon.hole.ru>" --arch=i386 -d precise pcre3_8.31-3.dsc
    reprepro -b /srv/www/packages/ubuntu include precise pcre3_8.31-3.labs1_amd64.changes
    reprepro -b /srv/www/packages/ubuntu include precise pcre3_8.31-3.labs1_i386.changes


    sbuild --append-to-version=".labs1" -m "Alexander V. Nikolaev <avn@daemon.hole.ru>" -d precise -A highlight_3.9-1.dsc
    sbuild --append-to-version=".labs1" -m "Alexander V. Nikolaev <avn@daemon.hole.ru>" -d precise --arch=i386 highlight_3.9-1.dsc
    reprepro -b /srv/www/packages/ubuntu include precise highlight_3.9-1.labs1_amd64.changes
    reprepro -b /srv/www/packages/ubuntu include precise highlight_3.9-1.labs1_i386.changes
    sbuild --append-to-version=".labs1" -m "Alexander V. Nikolaev <avn@daemon.hole.ru>" -d precise -A gtk-doc_1.20-1.dsc
    reprepro -b /srv/www/packages/ubuntu include precise gtk-doc_1.20-1.labs1_amd64.changes
    

### Build glib and gmime itself

Obtain glib:

    wget http://archive.ubuntu.com/ubuntu/pool/main/g/glib2.0/glib2.0_2.38.1-0ubuntu1.dsc
    http://archive.ubuntu.com/ubuntu/pool/main/g/glib2.0/glib2.0_2.38.1.orig.tar.xz
    http://archive.ubuntu.com/ubuntu/pool/main/g/glib2.0/glib2.0_2.38.1-0ubuntu1.debian.tar.gz

(following next steps, prior `sbuild` run better to done in docker, becasuse
it require a lot of additional packages installed, but I done them just on
another debian machine, where all of these packages was already installed:
`dh-autoreconf`, `cdbs`, and big part of build requirements)

Unpack source package:

    dpkg-source -x glib2.0_2.38.1-0ubuntu1.dsc

Then edit `debian/control.in` file: we need to replace `python:any (>=2.6.6~)`
with `python (>= 2.6)`. After this create new changelog entry:

    dch -i

and put proper version in changelog entry: `2.38.1-0ubuntu1.labs1`
(you can have a bit different version). Then make new source package, and
bring it to machine with sbuild.

    dpkg-buildpackage -S -us -us

Next step build packages, and put them to repo:

    sbuild -m "Alexander V. Nikolaev <avn@daemon.hole.ru>" -d precise -A glib2.0_2.38.1-0ubuntu1.labs1.dsc 
    sbuild -m "Alexander V. Nikolaev <avn@daemon.hole.ru>" -d precise --arch=i386 glib2.0_2.38.1-0ubuntu1.labs1.dsc 
    reprepro -b /srv/www/packages/ubuntu include precise glib2.0_2.38.1-0ubuntu1.labs1_amd64.changes glib2.0_2.38.1-0ubuntu1.labs1_i386.changes
    reprepro -b /srv/www/packages/ubuntu include precise glib2.0_2.38.1-0ubuntu1.labs1_i386.changes

For gmime we need same repackage of source, but with different (and optional) modification. I added `libglib2.0-0 (>= 2.38.1-0ubuntu1.labs1)` to `Depends:` line of `libgmime-2.6-0` section in `debian/control`.

    sbuild -m "Alexander V. Nikolaev <avn@daemon.hole.ru>" -d precise -A gmime_2.6.19-2.labs1.dsc 
    sbuild -m "Alexander V. Nikolaev <avn@daemon.hole.ru>" -d precise --arch=i386 gmime_2.6.19-2.labs1.dsc
    reprepro -b /srv/www/packages/ubuntu include precise gmime_2.6.19-2.labs1_amd64.changes 
    reprepro -b /srv/www/packages/ubuntu include precise gmime_2.6.19-2.labs1_i386.changes


## CentOS

### Build environment

I use docker environment with `centos:6.4` for build process

```sh
yum install rpm-build gcc-c++
yum install gettext libattr-devel libselinux-devel glibc-devel systemtap-sdt-devel zlib-devel automake autoconf libtool gtk-doc python-devel libffi-devel elfutils-libelf-devel chrpath gamin-devel
```

NOTE: Later we rebuild and upgrade `gtk-doc` and `automake`

### Obtain and build glib2

```sh
wget http://www.nic.funet.fi/pub/mirrors/fedora.redhat.com/pub/fedora/linux/development/rawhide/source/SRPMS/g/glib2-2.40.0-1.fc21.src.rpm
rpmbuild --rebuild glib2-2.40.0-1.fc21.src.rpm
```

Install resulting rpms

```sh
rpm -Uvh /rpmbuild/RPMS/x86_64/glib2-2.40.0-1.el6.x86_64.rpm /rpmbuild/RPMS/x86_64/glib2-devel-2.40.0-1.el6.x86_64.rpm
```

During upgrade we can get message, about error executing `gio-querymodule-64`,
It is harmless, because we have not any `gio` modules, and `gio-querymodule-64`
work normally after install (it is just possibe upgrade issue).

### Prerequisites for GMime

Build and install fresh gtk-doc from `rawhide`:

```sh
yum install gnome-doc-utils boost-devel help2man ctags libxslt-devel
wget http://www.nic.funet.fi/pub/mirrors/fedora.redhat.com/pub/fedora/linux/development/rawhide/source/SRPMS/s/source-highlight-3.1.7-1.fc21.src.rpm
rpmbuild --rebuild source-highlight-3.1.7-1.fc21.src.rpm
rpm -Uvh /rpmbuild/RPMS/x86_64/source-highlight-3.1.7-1.el6.x86_64.rpm /rpmbuild/RPMS/x86_64/source-highlight-devel-3.1.7-1.el6.x86_64.rpm
wget http://www.nic.funet.fi/pub/mirrors/fedora.redhat.com/pub/fedora/linux/development/rawhide/source/SRPMS/i/itstool-1.2.0-4.fc20.src.rpm
rpmbuild --rebuild itstool-1.2.0-4.fc20.src.rpm
rpm -i /rpmbuild/RPMS/noarch/itstool-1.2.0-4.el6.noarch.rpь
wget http://www.nic.funet.fi/pub/mirrors/fedora.redhat.com/pub/fedora/linux/development/rawhide/source/SRPMS/y/yelp-xsl-3.12.0-1.fc21.src.rpm
rpmbuild --rebuild yelp-xsl-3.12.0-1.fc21.src.rpm
rpm -Uvh /rpmbuild/RPMS/noarch/yelp-xsl-3.12.0-1.el6.noarch.rpm /rpmbuild/RPMS/noarch/yelp-xsl-devel-3.12.0-1.el6.noarch.rpm
wget http://www.nic.funet.fi/pub/mirrors/fedora.redhat.com/pub/fedora/linux/development/rawhide/source/SRPMS/y/yelp-tools-3.12.0-1.fc21.src.rpm
rpmbuild --rebuild yelp-tools-3.12.0-1.fc21.src.rpm
rpm -Uvh /rpmbuild/RPMS/noarch/yelp-tools-3.12.0-1.el6.noarch.rpm
wget http://www.nic.funet.fi/pub/mirrors/fedora.redhat.com/pub/fedora/linux/development/rawhide/source/SRPMS/g/gtk-doc-1.20-1.fc21.src.rpm
rpmbuild --rebuild gtk-doc-1.20-1.fc21.src.rpm
rpm -Uvh /rpmbuild/RPMS/noarch/gtk-doc-1.20-1.el6.noarch.rpm
```


```sh
 yum install python-mako intltool gnome-common freetype-devel libXft-devel fontconfig-devel libX11-devel libXfixes-devel libxml2-devel mesa-libGL-devel flex bison gpgme-devel libgpg-error-devel gettext-devel
yum install libXrender-devel libX11-devel libpng-devel libxml2-devel pixman-devel freetype-devel fontconfig-devel librsvg2-devel mesa-libGL-devel mesa-libEGL-devel
wget http://www.nic.funet.fi/pub/mirrors/fedora.redhat.com/pub/fedora/linux/development/20/source/SRPMS/p/pixman-0.30.0-3.fc20.src.rpm
rpmbuild --rebuild pixman-0.30.0-3.fc20.src.rpm
rpm -Uvh /rpmbuild/RPMS/x86_64/pixman-*
wget http://www.nic.funet.fi/pub/mirrors/fedora.redhat.com/pub/fedora/linux/development/rawhide/source/SRPMS/c/cairo-1.13.1-0.1.git337ab1f.fc21.src.rpm
rpmbuild --rebuild cairo-1.13.1-0.1.git337ab1f.fc21.src.rpm
rpm -Uvh /rpmbuild/RPMS/x86_64/cairo-1.13.1-0.1.git337ab1f.el6.x86_64.rpm /rpmbuild/RPMS/x86_64/cairo-gobject-1.13.1-0.1.git337ab1f.el6.x86_64.rpm /rpmbuild/RPMS/x86_64/cairo-devel-1.13.1-0.1.git337ab1f.el6.x86_64.rpm /rpmbuild/RPMS/x86_64/cairo-gobject-devel-1.13.1-0.1.git337ab1f.el6.x86_64.rpm
wget http://www.nic.funet.fi/pub/mirrors/fedora.redhat.com/pub/fedora/linux/development/rawhide/source/SRPMS/g/gobject-introspection-1.40.0-1.fc21.src.rpm
rpmbuild --rebuild gobject-introspection-1.40.0-1.fc21.src.rpm
```

Also `gonject-introspection` require some manual patches:

 * remove `perl-macros` from `Build-Requirement:`
 * comment out following two lines:
   ```
   #%dir %{_datadir}/gtk-doc/html/gi
   #%{_datadir}/gtk-doc/html/gi/*
   ```

Then finish `gobject-introspection` build:

```sh
rpmbuild -ba /rpmbuild/SPECS/gobject-introspection.spec
rpm -Uvh /rpmbuild/RPMS/x86_64/gobject-introspection-1.40.0-1.el6.x86_64.rpm /rpmbuild/RPMS/x86_64/gobject-introspection-devel-1.40.0-1.el6.x86_64.rpm
```

I think dependency on `Vala` is avoidable, but cheaper to build it from `rawhide`, than patch specs:

```sh
yum install emacs emacs-el
wget http://www.nic.funet.fi/pub/mirrors/fedora.redhat.com/pub/fedora/linux/development/rawhide/source/SRPMS/v/vala-0.24.0-1.fc21.src.rpm
rpmbuild --rebuild vala-0.24.0-1.fc21.src.rpm
rpm -Uvh /rpmbuild/RPMS/x86_64/vala-0.24.0-1.el6.x86_64.rpm /rpmbuild/RPMS/x86_64/vala-devel-0.24.0-1.el6.x86_64.rpm /rpmbuild/RPMS/x86_64/vala-tools-0.24.0-1.el6.x86_64.rpm
```

And final stanzas: CentOS 6 have enought broken `autoconf`, which unable to
process `configure.ac` from gmime.
(we take `m4` from `Fedora Core 20`, because requirements of `m4` make unwanted
dependency loop for us)

```sh
yum install perl-macros gcc-gfortran
wget http://www.nic.funet.fi/pub/mirrors/fedora.redhat.com/pub/fedora/linux/development/20/source/SRPMS/m/m4-1.4.16-10.fc20.src.rpm
rpmbuild --rebuild m4-1.4.16-10.fc20.src.rpm
rpm -Uvh /rpmbuild/RPMS/x86_64/m4-1.4.16-10.el6.x86_64.rpm
wget http://www.nic.funet.fi/pub/mirrors/fedora.redhat.com/pub/fedora/linux/development/rawhide/source/SRPMS/a/autoconf-2.69-14.fc21.src.rpm
rpmbuild --rebuild autoconf-2.69-14.fc21.src.rpm
```

We can get build failure, due failing testsuite.
Edit `/rpmbuild/SPECS/autoconf.spec`
And finish build with `rpmbuild -ba /rpmbuild/SPECS/autoconf.spec`
(this failure not affect on building of gmime, but I am suggest not use this
`autoconf` for anything, except building gmime)

```sh
rpm -Uvh /rpmbuild/RPMS/noarch/autoconf-2.69-14.el6.noarch.rpm
```

Then we can build `gmime` itself:

```sh
wget http://www.nic.funet.fi/pub/mirrors/fedora.redhat.com/pub/fedora/linux/development/rawhide/source/SRPMS/g/gmime-2.6.20-1.fc21.src.rpm
rpmbuild --rebuild gmime-2.6.20-1.fc21.src.rpm
```
