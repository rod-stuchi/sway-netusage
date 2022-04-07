waybar-netusage is an waybar companion that shows the network traffic.

How to add it to the waybar config:

```json
"custom/netusage": {
    "exec": "$HOME/.config/sway/waybar_custom_modules/sway_netusage --interface wlan0"
    //"on-click": "update-system",
},
```

You may specify the interface if autodetection does not work for you:
```json
"custom/netusage": {
  "exec": "$HOME/.config/sway/waybar_custom_modules/sway_netusage --interface wlan0"
},
```

style.css
```css
#custom-netusage {
    background-color: #4d4d4d;
    /*color: #000000;*/
}
```

Sample output:
```
19.1 MiB/s↓   21.2 MiB/s↑
```