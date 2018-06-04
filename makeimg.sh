#!/bin/bash

apk add syslinux
apk add xorriso

# build our initial RAM file system
mkdir -p ramfs
cp /vagrant/scratchmachine.upx ramfs/init
cp /lib/modules/4.9.73-0-virthardened/kernel/drivers/net/ethernet/intel/e1000/e1000.ko ramfs/e1000.ko

# Make our own initramfs, with just our binary
pushd ramfs
cat <<EOF | cpio -o -H newc | gzip > initramfs.gz
init
e1000.ko
EOF
popd

# We're going to build our iso in cdroot. And we need dev to exist under that
mkdir -p cdroot/dev
# We need a kernel
mkdir -p cdroot/kernel
cp /boot/vmlinuz-virthardened cdroot/kernel
# and our initramfs
cp ramfs/initramfs.gz cdroot
# ISOLINUX is our bootloader. 
mkdir -p cdroot/isolinux
cp /usr/share/syslinux/isolinux.bin cdroot/isolinux
cp /usr/share/syslinux/ldlinux.c32 cdroot/isolinux

cat <<EOF > cdroot/isolinux/isolinux.cfg
DEFAULT linux
  SERIAL 0 115200
  SAY Now booting the kernel from ISOLINUX...
  LABEL linux
  KERNEL /kernel/vmlinuz-virthardened
  INITRD /initramfs.gz
  APPEND root=/dev/ram0 ro console=tty0 console=ttyS0,115200
EOF

# Now make the ISO
mkisofs -o /vagrant/output.iso \
   -cache-inodes -J -l \
   -b isolinux/isolinux.bin -c isolinux/boot.cat \
   -no-emul-boot -boot-load-size 4 -boot-info-table \
   cdroot/

# references
# https://wiki.alpinelinux.org/wiki/Bootloaders
# https://www.syslinux.org/wiki/index.php?title=ISOLINUX
# https://www.kernel.org/doc/html/v4.12/admin-guide/initrd.html