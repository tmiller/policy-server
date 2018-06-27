#!/bin/bash

systemctl stop policy-server-monitor.timer >/dev/null || true
systemctl stop policy-server.service >/dev/null || true
systemctl disable policy-server-monitor.timer >/dev/null || true
systemctl disable policy-server.service >/dev/null || true
systemctl --system daemon-reload >/dev/null || true
