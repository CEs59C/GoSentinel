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


## Вывод
```bash
./sentinelVAmd64
=======cpuInfo=======
CPU 0:
  Vendor: AuthenticAMD
  Model: QEMU Virtual CPU version 2.5+
  Частота: 4292 MHz
  Кэш: 512
  Ядер: 1
=======memoryInfo=======
Всего: 961 MB
Доступно: 552 MB
Использовано: 409 MB (42.54%)
Свободно: 125 MB
=======swapInfo=======
Swap Total: 1916 MB
Swap Used: 220 MB (11.52%)
Swap Free: 1696 MB
=======discInfo=======
Всего: 29 GB
Использовано: 10 GB
Свободно: 17 GB
Использовано: 37.76%
Inodes used: 7.51%
=======netInfo (Listening Ports)=======
PID: 736196 | Порт: 27017 | Адрес: 127.0.0.1
PID: 594045 | Порт: 53    | Адрес: 127.0.0.53
PID: 693    | Порт: 25413 | Адрес: 127.0.0.1
PID: 594045 | Порт: 53    | Адрес: 127.0.0.54
PID: 1      | Порт: 22    | Адрес: 0.0.0.0
PID: 502765 | Порт: 28260 | Адрес: 127.0.0.1
PID: 688    | Порт: 28262 | Адрес: 127.0.0.1
PID: 717    | Порт: 444   | Адрес: ::
PID: 689    | Порт: 443   | Адрес: ::
PID: 1      | Порт: 22    | Адрес: ::
=======hostInfo=======
Имя хоста: vps-7960
Время работы: 632h45m13s
Количество процессов: 107
ОС: linux
Платформа: ubuntu
Семейство: debian
Версия: 24.04
Виртуализация: kvm (guest)
Load Average:
  1 минута: 0.10
  5 минут: 0.05
  15 минут: 0.06
Запущенных процессов: 2
Заблокированных: 0
Создано процессов: 850416
```