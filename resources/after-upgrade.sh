#!/bin/bash

systemctl --system daemon-reload >/dev/null || true
if ! systemctl is-enabled policy-server.service >/dev/null
then
    systemctl enable policy-server.service >/dev/null || true
    systemctl start policy-server.service >/dev/null || true
else
    systemctl restart policy-server.service >/dev/null || true
fi

if ! systemctl is-enabled policy-server-monitor.timer >/dev/null
then
    systemctl enable policy-server-monitor.timer >/dev/null || true
    systemctl start policy-server-monitor.timer >/dev/null || true
else
    systemctl restart policy-server-monitor.timer >/dev/null || true
fi
