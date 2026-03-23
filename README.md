# GoSentinel

### Todo
- [ ] сделать описание 
  - [ ] Ru
  - [ ] Eng
- [ ] модуль email
- [ ] модуль collector
  - [ ] cpu
  - [ ] memory
  - [ ] disk
  - [ ] process
  - [ ] net
- [ ] модуль report
  - [ ] formatter
  - [ ] Плановый отчет (Daily): 1 раз в сутки (например, в 9:00 утра)
  - [ ] Триггерный отчет (Alerts): Самое важное. Программа должна собирать данные каждые 5-10 минут, но отправлять письмо ТОЛЬКО ЕСЛИ:
    - [ ] RAM Used > 90%
    - [ ] Disk Used > 90%
    - [ ] Inodes Used > 90%
    - [ ] Появился новый процесс в netInfo, которого нет в «белом списке».

- [ ] версия для prometheus/grafana?

```bash
rm sentine*
```
## silicone 

```bash
#arm64 silicone
GOOS=linux GOARCH=arm64 go build -o sentinelVArd64 ./cmd/sentinel
```

```bash
path_to_file="/Users/onnikorpella/GolandProjects/goSentinel/sentinelT"
scp -P 2222 $path_to_file q@127.0.0.1:~/
```

## linux
```bash
#amd64 linux
GOOS=linux GOARCH=amd64 go build -o sentinelVAmd64 ./cmd/sentinel
```

```bash
path_to_file="/Users/onnikorpella/GolandProjects/goSentinel/sentinelVAmd64"
scp -P 22 $path_to_file root@85.155.101.203:~/
```
```bash
ssh root@85.155.101.203 "~/sentinelVAmd64" > local_result.txt
```

```bash
echo "Start"
rm sentine*
#amd64 linux
GOOS=linux GOARCH=amd64 go build -o sentinelVAmd64 ./cmd/sentinel
path_to_file="/Users/onnikorpella/GolandProjects/goSentinel/sentinelVAmd64"
scp -P 22 $path_to_file root@85.155.101.203:~/
ssh root@85.155.101.203 "~/sentinelVAmd64" > local_result.txt
echo "Done"
```


## Вывод
```txt
CPU Info:	Model=QEMU Virtual CPU version 2.5+, Vendor=AuthenticAMD, Cores=1, Usage=16.35%.
Disk:		Total=29GB, Used=10GB (37.9%) Free=17GB (62.1%), Inodes=7.5%.
Host:		vps-7960 [ubuntu 24.04], Uptime=663h35m24s, Processes=111, Running=3, Blocked=0, Created=905815, VM=kvm (guest)
Memory:		Total=961MB, Available=532MB, Used=429MB (44.67%), Free=149MB
Swap:		Total=1916MB, Used=222MB (11.59%), Free=1694MB
User:		root, Host=45.9.212.11, Started=17:27:25, Terminal=pts/0.
User:		root, Host=45.9.212.11, Started=17:27:31, Terminal=pts/1.
Process: mongod               Port: 27017 PID: 736196
Process: systemd-resolved     Port: 53    PID: 594045
Process: hysteria             Port: 25413 PID: 693
Process: systemd-resolved     Port: 53    PID: 594045
Process: systemd              Port: 22    PID: 1
Process: python               Port: 28260 PID: 502765
Process: user_auth            Port: 28262 PID: 688
Process: xray                 Port: 444   PID: 717
Process: caddy                Port: 443   PID: 689
Process: systemd              Port: 22    PID: 1

Письмо успешно отправлено!
```