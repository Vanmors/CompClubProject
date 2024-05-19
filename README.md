# CompClubProject

Требуется написать прототип системы, которая следит за работой компьютерного клуба, 
обрабатывает события и подсчитывает выручку за день и время занятости каждого стола.

## Инструкция по запуску
Создайте Docker-образ
```shell
docker build -t comp_club_project .
```
Запустите контейнер из ранее собранного образа
```shell
docker run comp_club_project
```

## Тестовые сценарии
### 1 input
```
3
09:00 19:00
10
08:48 1 client1
09:41 1 client1
09:48 1 client2
09:52 3 client1
09:54 2 client1 1
10:25 2 client2 2
10:58 1 client3
10:59 2 client3 3
11:30 1 client4
11:35 2 client4 2
11:45 3 client4
12:33 4 client1
12:43 4 client2
15:52 4 client4
19:25 1 client1
```
### 1 output
```
09:00
08:48 1 client1
08:48 13 NotOpenYet
09:41 1 client1
09:48 1 client2
09:52 3 client1
09:52 13 ICanWaitNoLonger!
09:54 2 client1 1
10:25 2 client2 2
10:58 1 client3
10:59 2 client3 3
11:30 1 client4
11:35 2 client4 2
11:35 13 PlaceIsBusy
11:45 3 client4
12:33 4 client1
12:33 12 client4 1
12:43 4 client2
15:52 4 client4
19:00 11 client3
19:25 1 client1
19:25 13 NotOpenYet
19:00
1 70 05:58
2 30 02:18
3 90 08:01
```
### 2 input (ClientUnknown)
```
2
08:00 17:00
20
08:30 1 client1
08:45 2 client1 1
09:00 2 client1 2
09:15 4 client3
```
### 2 output (ClientUnknown)
```
08:00
08:30 1 client1
08:45 2 client1 1
09:00 2 client1 2
09:15 4 client3
09:15 13 ClientUnknown
17:00 11 client1
17:00
1 20 00:15
2 160 08:00
```

### 3 input (ErrorInput)
```
2
08:00 17:00
20
08 30 1 client1
```

### 3 output (ErrorInput)
```
Input format error: 08 30 1 client1
```


