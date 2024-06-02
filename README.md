# Распределённый вычислитель арифметических выражений

## Описание проекта

Этот проект реализует распределённую систему вычисления арифметических выражений. Система состоит из двух компонентов:

- **Оркестратор**: Сервер, который принимает арифметические выражения, разбивает их на задачи и распределяет их выполнение.
- **Агент**: Демон, который получает задачи от оркестратора, выполняет их и возвращает результаты.

## Установка и запуск проекта

### Предварительные требования

- Docker
- Docker Compose

### Установка Docker и Docker Compose

#### Windows

1. Скачайте и установите [Docker Desktop для Windows](https://www.docker.com/products/docker-desktop).
2. Убедитесь, что Docker Desktop установлен и работает.
3. Docker Compose включён в состав Docker Desktop, отдельной установки не требуется.

#### MacOS

1. Скачайте и установите [Docker Desktop для Mac](https://www.docker.com/products/docker-desktop).
2. Убедитесь, что Docker Desktop установлен и работает.
3. Docker Compose включён в состав Docker Desktop, отдельной установки не требуется.

#### Linux

1. Следуйте инструкциям по установке Docker для вашей версии дистрибутива [здесь](https://docs.docker.com/engine/install/).
2. Следуйте инструкциям по установке Docker Compose [здесь](https://docs.docker.com/compose/install/).

## **ВАЖНО:** Недавно докер прижали американцы со своими санкциями и теперь у россиян возникает соответствующая ошибка при скачивании контейнеров с Dockerhub
### Шаги исправления проблемы:

#### Если используется Docker Desktop

1. Перейти в настройки (шестеренка вверху)

2. Перейти во вкладку Docker Engine

3. Добавить параметр ```"registry-mirrors": ["https://dockerhub.timeweb.cloud"]``` в конфиг

    Пример моего конфига:
    ```json
    {
      "builder": {
        "gc": {
          "defaultKeepStorage": "20GB",
          "enabled": true
        }
      },
      "experimental": false,
      "registry-mirrors": [
        "https://dockerhub.timeweb.cloud"
      ]
    }
    ```

#### Если вы линуксойд

  Следуйте инструкциям по изменению конфига на [сайте самой прокси](https://dockerhub.timeweb.cloud), изменения конфига аналогичны примеру выше

### Шаги установки и запуска

1. Клонируйте репозиторий:

    ```sh
    git clone <URL_вашего_репозитория>
    cd <название_папки_с_репозиторием>
    ```

2. Убедитесь, что Docker и Docker Compose установлены и работают на вашей машине.

3. Отредактируйте переменные окружения в скрипте `run.sh` или `run.ps1`, если это необходимо:

    **Для Windows (run.ps1)**:

    ```powershell
    # Установка переменных окружения
    $env:ORCH_URL = "http://orch:8080/internal/task"
    $env:COMPUTING_POWER = 3
    $env:TIME_ADDITION_MS = 10
    $env:TIME_SUBTRACTION_MS = 10
    $env:TIME_MULTIPLICATION_MS = 10
    $env:TIME_DIVISION_MS = 10

    # Запуск Docker Compose
    docker-compose up --build --remove-orphans
    ```

    **Для MacOS и Linux (run.sh)**:

    ```sh
    export ORCH_URL="http://orch:8080/internal/task"
    export COMPUTING_POWER=3
    export TIME_ADDITION_MS=10
    export TIME_SUBTRACTION_MS=10
    export TIME_MULTIPLICATION_MS=10
    export TIME_DIVISION_MS=10
    docker-compose up --build --remove-orphans
    ```

4. **Для Windows:**

    Откройте PowerShell и перейдите в директорию с вашим проектом:

    ```powershell
    cd путь/к/вашему/проекту
    ```

    Запустите скрипт `run.ps1`:

    ```powershell
    .\run.ps1
    ```

    **Для MacOS и Linux:**

    Откройте терминал и перейдите в директорию с вашим проектом:

    ```sh
    cd путь/к/вашему/проекту
    ```

    Запустите скрипт `run.sh`:

    ```sh
    ./run.sh
    ```

### Переменные окружения

Переменные окружения, используемые в проекте:

- `COMPUTING_POWER`: Количество горутин, используемых агентом для выполнения задач.
- `TIME_ADDITION_MS`: Время выполнения операции сложения в миллисекундах.
- `TIME_SUBTRACTION_MS`: Время выполнения операции вычитания в миллисекундах.
- `TIME_MULTIPLICATION_MS`: Время выполнения операции умножения в миллисекундах.
- `TIME_DIVISION_MS`: Время выполнения операции деления в миллисекундах.

### Тестирование и проверка работоспособности

1. **Добавление выражения для вычисления:**
    Пример:
    ```sh
    curl --location 'localhost:8080/api/v1/calculate' \
    --header 'Content-Type: application/json' \
    --data '{
        "expression": "2 + 2 * 2" #любое арифметическое выражение с целочислеными и/или дробными
    }'                            #НО НЕ ОТРИЦАТЕЛЬНЫМИ числами, скобками и знаками + - * / соответствующих операций
    ```

2. **Получение списка выражений:**
    ```sh
    curl --location 'localhost:8080/api/v1/expressions'
    ```

3. **Получение выражения по его идентификатору:**
    Пример:
    ```sh
    curl --location 'localhost:8080/api/v1/expressions/1' #вместо 1 id вашего вычисления
    ```

## Заключение

После выполнения указанных шагов система будет запущена и готова к использованию. Вы можете добавлять арифметические выражения, получать их статус и результаты.

# Если у вас возникнут вопросы или проблемы, пожалуйста, прочитайте ещё раз внимательно данный файл или задайте мне вопрос в телеграм @aster_cmd