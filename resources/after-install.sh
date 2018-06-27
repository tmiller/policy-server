#!/bin/bash

systemctl --system daemon-reload >/dev/null || true
systemctl enable policy-server.service >/dev/null || true
systemctl enable policy-server-monitor.timer >/dev/null || true
systemctl start policy-server.service >/dev/null || true
systemctl start policy-server-monitor.timer >/dev/null || true
