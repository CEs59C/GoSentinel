# GoSentinel

### Todo
- [ ] сделать описание 
  - [ ] Ru
  - [ ] Eng
- [x] модуль email
- [x] модуль collector
  - [x] cpu
  - [x] memory
  - [x] disk
  - [x] process
  - [x] net
- [x] модуль report
  - [x] formatter
  - [ ] Плановый отчет (Daily): 1 раз в сутки (например, в 9:00 утра)
  - [ ] Триггерный отчет (Alerts): Самое важное. Программа должна собирать данные каждые 5-10 минут, но отправлять письмо ТОЛЬКО ЕСЛИ:
    - [ ] RAM Used > 90%
    - [ ] Disk Used > 90%
    - [ ] Inodes Used > 90%
    - [ ] Появился новый процесс в netInfo, которого нет в «белом списке».
- [ ] наладить работу с .env
- [x] сделать Make файл
- [ ] версия для prometheus?
  - [ ] тогда надо делать вэб-сервис с постоянной трансяцией данных на http://ip:9100/metrics

## Вывод/Письмо
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

## Схема работы с .env
```mermaid
flowchart LR
    subgraph COMPILE[ПРОЦЕСС КОМПИЛЯЦИИ]
        direction TB
        ENV[".env (исходный)"]
        
        subgraph  STEP1[Шифруем ТОЛЬКО пароли]
            POST_IN1["POST_IN=email<br/>(не тронуто)"]
            POST_TO1["POST_TO=recipient<br/>(не тронуто)"]
            PASSWORD1["PASSWORD=real123<br/>(зашифровано)"]
        end
        
        ENV --> STEP1
        
        subgraph BUILD[build/ структура]
            SENTINEL["sentinel (бинарник)"]
            ENV_BUILD[".env с зашифрованным паролем<br/>POST_IN=email<br/>POST_TO=recipient<br/>PASSWORD=ENC[AES256_GCM,data:abc123...]"]
        end

        STEP1 --> BUILD
    end

    subgraph RUNTIME[РАБОТА ПРИЛОЖЕНИЯ НА СЕРВЕРЕ]
        direction TB
        READ[1. Бинарник читает .env файл]
        DETECT[2. Обнаруживает зашифрованный пароль]
        DECRYPT[3. Пытается расшифровать с встроенным ключом]
        MEMORY[4. Получает реальный пароль в памяти]
        
        READ --> DETECT --> DECRYPT --> MEMORY
    end

    subgraph  TRAP[ЛОВУШКА ДЛЯ ЗЛОУМЫШЛЕННИКА]
        direction RL
        THEFT1["Украл бинарник → нет .env, нет пароля"]
        THEFT2["Украл .env → пароль зашифрован, не может расшифровать"]
        THEFT3["Сделал strings на бинарник → видит ложные пароли (ловушка)"]
        THEFT4["Запускает с подменой .env → использует подставные данные"]
    end

    COMPILE --> RUNTIME
    COMPILE -.-> TRAP

    style COMPILE fill:#e1f5fe,stroke:#01579b
    style RUNTIME fill:#fff3e0,stroke:#e65100
    style TRAP fill:#ffebee,stroke:#b71c1c
    style STEP1 fill:#f5f5f5,stroke:#9e9e9e
    style BUILD fill:#e8f5e9,stroke:#2e7d32
```

```bash
make encrypt
#echo $ENCRYPTION_KEY 
```

```bash
make run-remote-amd64
```
```bash
make run-local-arm64
```
```bash
make run-mac
```
