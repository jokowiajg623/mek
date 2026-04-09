#!/bin/bash
while true; do
    rm -f telnet.txt
    
    sudo zmap -p23 -w ip.txt -r 3000 -q 2> telnet.txt | python3 telnet.py > /tmp/telnet_output.log 2>&1 &
    
    SCAN_PID=$!
    echo "Start on Codespace - PID: $SCAN_PID (Rate: 3000)"

    while true; do
        cond1=$(grep -q "\[INFO\] zmap: completed" telnet.txt 2>/dev/null && echo 1 || echo 0)
        
        if ! kill -0 $SCAN_PID 2>/dev/null; then
            cond2="1"
        else
            cond2="0"
        fi

        if [[ "$cond1" == "1" && "$cond2" == "1" ]]; then
            echo "Scan Finished"
            break
        fi

        if ! pgrep -f zmap > /dev/null; then
            echo "Process interrupted"
            break
        fi

        sleep 10
    done
    
    pkill -f zmap
    pkill -f telnet.py
done
