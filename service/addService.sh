mkdir -p $HOME/.config/systemd/user/

echo '[Unit]
Description="Chatbox companion"

[Service]
User='$USER'
WorkingDirectory='$(pwd)'
ExecStart='$(pwd)'/chatbox
StandardError=null
Restart=always
Environment="OPENAI_KEY='$OPENAI_KEY'"

[Install]
WantedBy=multi-user.target' > $HOME/.config/systemd/user/chatbox.service
