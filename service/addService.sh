echo '[Unit]
Description="Chatbox companion"

[Service]
User='$SUDO_USER'
WorkingDirectory='$(pwd)'
ExecStart='$(pwd)'/chatbox
Restart=always
Environment="OPENAI_KEY='$OPENAI_KEY'"

[Install]
WantedBy=multi-user.target' > /etc/systemd/system/chatbox.service
