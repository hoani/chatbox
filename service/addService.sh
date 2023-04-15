mkdir -p $HOME/.config/systemd/user/

echo '[Unit]
Description="Chatbox companion"
After=pulseaudio.service
Wants=pulseaudio.service

[Service]
WorkingDirectory='$(pwd)'
ExecStart='$(pwd)'/chatbox
StandardError=null
Restart=always
Environment="OPENAI_KEY='$OPENAI_KEY'"

[Install]
WantedBy=default.target' > $HOME/.config/systemd/user/chatbox.service
