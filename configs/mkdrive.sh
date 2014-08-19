#!/bin/bash
for i in 0 1 2; do
  mkdir -p /tmp/new-drive/openstack/latest
  cp core${i}.yml /tmp/new-drive/openstack/latest/user_data
  mkisofs -R -V config-2 -o core${i}.iso /tmp/new-drive
  rm -r /tmp/new-drive
done
