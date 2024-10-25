build:
    go build -o waybar-netusage

run:
    go build -o waybar-netusage
    ./waybar-netusage -down -interface wlo1
