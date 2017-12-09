response=$(curl -ksf https://localhost:843 | head -1)
check='<?xml version="1.0"?>'
if [ "$response" = "$check" ]; then
        echo "Working - "$(date)
else
        echo "It's down! - "$(date)
        sudo service policy-server-go restart
fi
