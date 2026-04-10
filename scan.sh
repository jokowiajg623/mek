#!/bin/bash
while true; do
    rm -f telnet.txt

    sudo zmap -p23 -w ip.txt -r 4000 -q 2> telnet.txt | python3.11 telnet.py &

    echo "Start"

    while true; do
        cond1=$(grep -q "\[INFO\] zmap: completed" telnet.txt 2>/dev/null && echo 1 || echo 0)

        # Karena tidak pakai screen, kita cek apakah proses zmap masih ada
        if ! pgrep -f zmap > /dev/null; then
            cond2="1"
        else
            cond2="0"
        fi

        if [[ "$cond1" == "1" && "$cond2" == "1" ]]; then
            pkill -f zmap
            pkill -f telnet.py
            echo "End"
            break
        fi

        if ! pgrep -f zmap > /dev/null; then
            break
        fi

        sleep 60
    done
done
