waybar-netusage is an waybar companion that shows the network traffic.

How to add it to the waybar config:

```json
"custom/netusage": {
    "exec": "$HOME/..path_to../waybar-netusage -interface wlo1",
    "return-type": "json"
},
```

You may specify the interface if autodetection does not work for you:
```json
"custom/netusage": {
    "exec": "$HOME/..path_to../waybar-netusage -down -interface wlo1",
    "return-type": "json"
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
Sample output with `-down` argument (upload status will be displayed in tooltip):
```
5.1 MiB/s↓

```

Sample output with `-up` argument (download output will be displayed in tooltip):
```
3.0 MiB/s↑
```
