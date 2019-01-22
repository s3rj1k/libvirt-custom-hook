# libvirt-mhst-hook

Docs: 
  - https://www.libvirt.org/hooks.html

Install:
  - copy `qemu` to `/etc/libvirt/hooks/qemu`
  - copy `qemu-hook.json` to `/etc/libvirt/hooks/qemu-hook.json`
  - restart libvirt daemon `systemctl restart libvirtd`

Todo:
  - add tests for config validation
