books-checker
---

Файл `exportfsm.csv` со списком запрещенных материалов сохранен в папку data по состоянию на 10.07.2023

В последующем, при обновлении файл, скаченный с сайта минюста — `exportfsm.csv`, необходимо сохранить в кодировке UTF-8. На сайте минюста этот файл в кодировке windows-1251. Manticoresearch создаст некорректные данные в базе, если не перевести файл в кодировку UTF-8
___
- Обязательное ПО для запуска **Docker**, **Git**
- Необязательные программы Postman

Если у вас уже установленны необходимые программы и создан проект, переходите к шагу 2.

### 1. Подготовка к запуску

<details><summary>Установка программ и создание проекта</summary>
<p>

1. Установить Docker для Windows https://docs.docker.com/desktop/install/windows-install/

2. Установить Git для windows https://git-scm.com/downloads При установке можно оставить все по умолчанию

3. Установить программу Postman https://www.postman.com/downloads/

После установки необходимых программ, нажмите меню пуск, найдите git консоль Git CMD, запустите Git CMD

Выберите или создайте папку с проектами, например c:\terratensor,

Для создания папки наберите в Git CMD консоли:
```
mkdir ./terratensor
```

Выберите созданную папку, наберите в Git CMD консоли:
```
cd ./terratensor
```

Для клонирования репозитория terratensor/books-checker, наберите в Git CMD консоли:

```
git clone https://github.com/terratensor/books-checker.git
```

Увидите сообщение в консоли об успешном клонировании:
```
Cloning into 'books-checker'...
fatal: helper error (-1): User cancelled dialog.
Username for 'https://github.com':
c:\terratensor>git clone https://github.com/terratensor/common_library_parser.git
Cloning into 'common_library_parser'...
remote: Enumerating objects: 20, done.
remote: Counting objects: 100% (20/20), done.
remote: Compressing objects: 100% (15/15), done.
Receiving objects: 100% (20/20), 6.48 KiB | 6.48 MiB/s, done.
Resolving deltas: 100% (4/4), done.
```

После наберите в Git CMD консоли:

```
cd ./books-checker
```

Затем наберите в консоли
```
ls -la
```

Увидите следующую структуру папки:
```
total 32
drwxrwx---+ 1 username Domain Users    0 Jun 22 13:30 .
drwxrwx---+ 1 username Domain Users    0 Jun 22 13:30 ..
drwxrwx---+ 1 username Domain Users    0 Jun 22 13:30 .git
-rwxrwx---+ 1 username Domain Users 8880 Jun 22 13:30 README.md
drwxrwx---+ 1 username Domain Users    0 Jun 22 13:30 docker
-rwxrwx---+ 1 username Domain Users  957 Jun 22 13:30 docker-compose.yml
drwxrwx---+ 1 username Domain Users    0 Jun 22 13:30 process
```

</p>
</details>

### 2. Загрузка списка Минюста в БД Manticoresearch

Запустите windows docker (через меню пуск)

После запуска необходимо проверить, что нет запущенных контейнеров с manticoresearch (меню containers), если есть остановить их или удалить, иначе будет конфликт доступа к портам.

В консоли Git CMD, где ранее создали папку и клонировали репозиторий для запуска проекта набрать:

```
docker-compose down --remove-orphans
```

```
docker compose up --build -d
```

Запустится база данных Manticoresearch

Для записи списка в БД Manticoresearch наберите:

```
docker compose exec -it manticore indexer minjust_list
```
После завершения процесса создания базы данных индекса увидите сообщение: 


### 3. Создание списка книг для проверки

Скачайте последнюю версию утилиты books-checker.exe.

Ссылка на файл находится в секции Assets страницы описания релиза

https://github.com/terratensor/books-checker/releases/latest

Сохраните books-checker.exe в папке с проектом ./books-checker

Для запуска парсера url (интернет) списка страниц с библиотеки наберите:
```
./bookschecker.exe -p
```
Парсер последовательно обработает все страницы списка алфавитного указателя и сохранит csv файл в папку `./csv`  

### 4. Проверка книг на совпадение в списке Миньюста 

Для запуска проверки книг на соответствие списку Миньюста запустите:

```
./bookschecker.exe -f %filename%
```

В качестве значения `%filename%` укажите имя файла, который вы создали с помощью парсера на предыдущем шаге.

