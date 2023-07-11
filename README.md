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
remote: Enumerating objects: 47, done.
remote: Counting objects: 100% (47/47), done.
remote: Compressing objects: 100% (30/30), done.
remote: Total 47 (delta 15), reused 37 (delta 10), pack-reused 0
Unpacking objects: 100% (47/47), 566.03 KiB | 613.00 KiB/s, done.
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
drwxrwxrwx 1 username Domain 4096 Jul 11 22:16 .
drwxrwxrwx 1 username Domain 4096 Jul 11 22:16 ..
drwxrwxrwx 1 username Domain 4096 Jul 11 22:16 .git
-rwxrwxrwx 1 username Domain  524 Jul 11 22:16 .gitignore
-rwxrwxrwx 1 username Domain 1526 Jul 11 22:16 LICENSE
-rwxrwxrwx 1 username Domain  416 Jul 11 22:16 Makefile
-rwxrwxrwx 1 username Domain 5806 Jul 11 22:16 README.md
drwxrwxrwx 1 username Domain 4096 Jul 11 22:16 app
drwxrwxrwx 1 username Domain 4096 Jul 11 22:16 csv
drwxrwxrwx 1 username Domain 4096 Jul 11 22:16 data
drwxrwxrwx 1 username Domain 4096 Jul 11 22:16 docker
-rwxrwxrwx 1 username Domain  809 Jul 11 22:16 docker-compose.yml
-rwxrwxrwx 1 username Domain  310 Jul 11 22:16 go.mod
-rwxrwxrwx 1 username Domain 1181 Jul 11 22:16 go.sum
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
```
Manticore 6.0.4 1a3a4ea82@230314
Copyright (c) 2001-2016, Andrew Aksyonoff
Copyright (c) 2008-2016, Sphinx Technologies Inc (http://sphinxsearch.com)
Copyright (c) 2017-2023, Manticore Software LTD (https://manticoresearch.com)

using config file '/etc/manticoresearch/manticore.conf'...
indexing table 'minjust_list'...
collected 5334 docs, 2.7 MB
creating lookup: 5.3 Kdocs, 100.0% done
sorted 0.4 Mhits, 100.0% done
total 5334 docs, 2708617 bytes
total 0.249 sec, 10868726 bytes/sec, 21403.46 docs/sec
total 3 reads, 0.002 sec, 990.5 kb/call avg, 0.7 msec/call avg
total 30 writes, 0.006 sec, 249.9 kb/call avg, 0.2 msec/call avg
```

После этого ОБЯЗАТЕЛЬНО перезапустите контейнер с Manticoresearch командами, изменения в БД появляются только после перезагрузки. Если не перезагрузить, база будет пустая и все проверки всегда будут ложно успешные.

```
docker-compose down --remove-orphans
```

```
docker compose up --build -d
```

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

`2023/07/11 22:22:59 file ./csv/222259_11072023_list.csv was successful writing`

### 4. Проверка книг на совпадение в списке Миньюста 

Для запуска проверки книг на соответствие списку Миньюста запустите:

`./bookschecker.exe -f %filename%`

`./books-checker.exe -f ./csv/222259_11072023_list.csv`

В качестве значения `%filename%` укажите имя файла, который вы создали с помощью парсера на предыдущем шаге.

После обработки увидите сообщение. Имя созданного файла в вашем случае будет другим.

`Обработка завершена. Создан файл: ./22-33-09_11072023_result_log.txt`

Далее можно посмотреть результаты обработки в текстовом файле, если там пусто, значит совпадений не найдено. И в списке книг нет запрещенных материалов.

При необходимости изменить значения настройки утилиты и задать свои значения, далее приведен список опции командной строки утилиты book-search:
```
Usage of books-checker.exe:
  -c, --columns string     колонки csv (author, title), которые будут соединены в строку запроса (default "all")
  -f, --file string        csv файл с наименованиями книг для проверки (default "./list.csv")
  -m, --matchMode string   режим поиска, query_string, match_phrase, match (default "query_string")
  -p, --parse              Парсинг страниц с наименованиями книг в файл csv
  -s, --showConsole        вывод результатов в консоль без сохранения в файл
```

#### Примеры использования:
Для изменения режима поиска укажите параметр `-m` со значениями:

- query_string — строка поиска по умолчанию. https://manual.manticoresearch.com/Searching/Full_text_matching/Basic_usage#query_string
- match_phrase — это запрос, который соответствует всей фразе. https://manual.manticoresearch.com/Searching/Full_text_matching/Basic_usage#match_phrase
- match — это простой запрос, который соответствует указанным ключевым словам в указанных полях(колонках). Все слова в запросе будут соединены через условный оператор `ИЛИ` https://manual.manticoresearch.com/Searching/Full_text_matching/Basic_usage#match


`./books-checker ./csv/222259_11072023_list.csv -m match_phrase`

В данном случае будет произведен поиск на соответствие всей фразы запроса
По умолчанию запрос создается конкатенацией (сложением двух колонок файла с книгами csv) автор книги и наименование книги, например:

`Абинякин Р. М. Офицерский корпус Добровольческой армии:`

При необходимости можно изменит это поведение и указать только колонку автора или колонку наименования книги для формирования поискового запроса, например:

`./books-checker ./csv/222259_11072023_list.csv -m match_phrase -с author` — проверка списка книг по автору
`./books-checker ./csv/222259_11072023_list.csv -m match_phrase -с title` — проверка списка книг по наименованию

Можно изменить режимы поиска:

При поиске совпадений только по автору или только по наименованию, результаты в фале скорее всего будут, но это ещё не значит, что книги находятся в списке, необходимо внимательно изучить результаты обработки, чтобы сделать окончательный вывод.

В режиме поиска `match` может быть очень много совпадения, это обусловлено самим механизмом данного метода поиска (запрос, который соответствует указанным ключевым словам в указанных полях(колонках). Все слова в запросе будут соединены через условный оператор `ИЛИ`)

`./books-checker ./csv/222259_11072023_list.csv -m query_string -с author`

`./books-checker ./csv/222259_11072023_list.csv -m query_string -с title`

`./books-checker ./csv/222259_11072023_list.csv -m match -с author`

`./books-checker ./csv/222259_11072023_list.csv -m match`


------------
`./books-checker ./csv/222259_11072023_list.csv -m match_phrase -c title`

```
...
2023/07/11 23:04:54 Строка: [Азаров И.И. Непобежденные 1973]
Запрос(match_phrase): Непобежденные
1. Стихотворение Маслова И.А. «Светослав <b>Непобежденный</b>» размещенное на сайте http:www.slavyanskaya-kultura.ru/literature/poetry/ilja-maslov-sbornik-stihov-severnoi-bolyu-i-ruskoi-pechalyu.html (решение Октябрьского районного суда г. Барнаула от 29.12.2014);
2. Размещенные в разделе «видеозаписи» на странице пользователя (электронный адрес - http://vk.com/rafburn) сайта международной сети Интернет www.vk.com следующие материалы: видеофильмы «Послание Путину из Сирии», «Обращение к президентам Шайхутдинова Ильдара», «Пресс-конференция исламской партии Хизб ут-Тахрир в УНИАН», «Обращение к властям РФ», «Первый день международного форума Хизб ут-Тахрир», «Пресс конференция, анонсирующая проведение Международного форума Хизб ут-Тахрир», «<b>Непобежденные</b>!!! Судьбой довольные!!!» (решение Альметьевского городского суда Республики Татарстан от 05.02.2015);

2023/07/11 23:04:54 Строка: [Азольский А.А. Диверсант ]
Запрос(match_phrase): Диверсант
1. Информационный материал – аудиофайл песня «Расист», обнаруженный в ходе мониторинга глобальной телекоммуникационной сети Интернет на сайте «http://vk.com/idl53056430» под псевдонимом «Женя <b>Диверсант</b>» (решение Артемовского городского суда Свердловской области от 03.07.2013);
...
```

`./books-checker.exe -f ./csv/222259_11072023_list.csv -c title`

```
...
2023/07/11 23:09:12 Строка: [Воинов А.И. Отважные 1961]
Запрос(query_string): Отважные 
1. Размещенные в сети Интернет, в том числе на веб сервисе сайта «vk.com» аудиозаписи «Apraxia - Russia, WakeUp!» длительностью 05 минут 12 секунд, «Apraxia –Русь, проснись!» длительностью от 04 минуты 44 секунды до 05 минут 12 секунд, Apraxia – 14-Русь, проснись!» длительностью 05 минут 10 секунд, «APRAXIA – «Русь, проснись!» длительностью 05 минут 10 секунд, «Apraxia/Молат – Русь, Проснись!» длительностью 05 минут 10 секунд, «apraxia – Русь проснись!!!» длительностью 05 минут 12 секунд, Apraxia & Молот – Русь, Проснись!» длительностью 05 минут 12 секунд, «Apraxia – Русь, проснись! (BL)» длительностью 05 минут 09 секунд, «Молот – Русь,проснись!» длительностью 04 минуты 44 секунды, «Apraxia Molot – Вставай, Священная Русь!» длительностью 05 минут 11 секунд, «Apraxia Tribute – Коло Прави / Русь, проснись!» длительностью 04 минуты 44 секунды, «коло прави (апраксия) – Русь, проснись!» длительностью 04 минуты 38 секунд, начинающихся словами: «Вставай, просыпайся, священная Русь. Из праха восстань...» и заканчивающихся словами: «...<b>отважных</b> борцов гремит - Проснись, наша Русь, проснись!» (решение Благовещенского городского суда Амурской области от 24.09.2018);
2. Видеозапись под названием «Роль евреев в работорговле - шокирующая правда. Часть 1», продолжительностью 11 минут 31 секунда, на которой демонстрируется выступление седого мужчины, речь которого начинается со слов: «Здравствуйте. Я Дэйвид Дюк. За исключением войны в истории человечества не было ничего, что принесло бы больше страданий, смертей и насилия, чем работорговля...» и заканчивается словами: «...Еврейский судья Стевин Гринберг вы нес приговор с требованием запретить трансляцию программы, пока он ее не одобрит», затем следует заставка с надписью: «Поддержи своим неравнодушием <b>отважную</b> работу Дэйвида Дюка» (решение Сыктывкарского городского суда Республики Коми от 22.06.2020);

2023/07/11 23:09:12 Строка: [Войтиков С.С. Армия и власть. 2016]
Запрос(query_string): Армия и власть. 
1. Информационный материал, поступивший с адреса электронной почты 22sud2010@gmail.ru от неизвестного лица, представляющий собой электронное письмо начинающееся словами «Соотечественники!...» <b>и</b> заканчивающееся «...для противодействия провокациям еврейской <b>власти</b> против вас лично <b>и</b> ваших народов, как моральное <b>и</b> юридическое основание для ликвидации еврейских оккупационных сил на нашей территории», состоящее из 12 разделов («Воззвание к русскому народу, к офицерам <b>армии и</b> флота, к казачеству, к русской молодежи <b>и</b> православному духовенству», «Уничтожение <b>армии</b>. Ликвидация военной разведки России», «Товарищи офицеры!», «Разъяснение о характере тайных операций против русского народа», «Война с русской молодежью», «Псевдонационалисты», «Война с молодежью Кавказа», «Организация войн <b>и</b> финансовый кризис», «Экономический террор», «Диверсионно-террористическая деятельность <b>и</b> психологический террор», «О разделении русского народа», «Задачи бойцов народного сопротивления»)» (решение Ленинского районного суда г. Саранска Республики Мордовия от 29.07.2010).
...
```
