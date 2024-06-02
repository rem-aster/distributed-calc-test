# Распределённый вычислитель арифметических выражений

## Описание проекта

Этот проект реализует распределённую систему вычисления арифметических выражений. Система состоит из двух компонентов:

- **Оркестратор**: Сервер, который принимает арифметические выражения, разбивает их на задачи с помощью постфиксной нотации и закидывает в очередь на выполнение а также создающее канал для получения результатов
- **Агент**: Демон, который получает задачи от оркестратора, выполняет их с задержкой, переданной оркестратором и возвращает результаты обратно.

Помимо этого в оркестраторе реализовано кеширование, позволяющее ускорить вычисления

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

1. Клонируйте репозиторий в нужную вам папку:

    ```sh
    git clone https://github.com/rem-aster/distributed-calc-test
    cd distributed-calc-test
    ```

2. Убедитесь, что Docker и Docker Compose установлены и работают на вашей машине.

3. Отредактируйте переменные окружения в скрипте `run.sh` или `run.ps1`, если это необходимо:

    **Для Windows (run.ps1)**:

    ```powershell
    $env:ORCH_URL = "http://orch:8080/internal/task"
    $env:COMPUTING_POWER = 3         #можно редактировать это
    $env:TIME_ADDITION_MS = 10       #это
    $env:TIME_SUBTRACTION_MS = 10    #это 
    $env:TIME_MULTIPLICATION_MS = 10 #это
    $env:TIME_DIVISION_MS = 10       #и это число

    docker-compose up --build --remove-orphans
    ```

    **Для MacOS и Linux (run.sh)**:

    ```sh
    export ORCH_URL="http://orch:8080/internal/task"
    export COMPUTING_POWER=3 #можно редактировать переменные аналогично примеру выше (run.ps1)
    export TIME_ADDITION_MS=10
    export TIME_SUBTRACTION_MS=10
    export TIME_MULTIPLICATION_MS=10
    export TIME_DIVISION_MS=10
    docker-compose up --build --remove-orphans
    ```

4. **Для Windows:**

    Откройте PowerShell и перейдите в директорию с репозиторием:

    ```powershell
    cd путь/к/папке/distributed-calc-test
    ```

    Убедитесь, что у вас есть права на выполнение PowerShell скриптов. Для этого выполните следующую команду в PowerShell ***с правами администратора***:

    ```powershell
    Set-ExecutionPolicy RemoteSigned
    ```

    Запустите скрипт `run.ps1`:

    ```powershell
    .\run.ps1
    ```

    **Для MacOS и Linux:**

    Откройте терминал и перейдите в директорию с вашим проектом:

    ```sh
    cd путь/к/папке/distributed-calc-test
    ```

    Запустите скрипт `run.sh`:

    ```sh
    ./run.sh
    ```

#### После успешного развертывания смело открываем новую вкладку консоли (***powershell*** в Windows или bash в Linux/MacOS) и курлим как в примерах ниже

### Переменные окружения

Переменные окружения, используемые в проекте, изменения которых вступят в силу только после старта/рестарта:

- `COMPUTING_POWER`: Количество горутин, используемых агентом для выполнения задач.
- `TIME_ADDITION_MS`: Время выполнения операции сложения в миллисекундах.
- `TIME_SUBTRACTION_MS`: Время выполнения операции вычитания в миллисекундах.
- `TIME_MULTIPLICATION_MS`: Время выполнения операции умножения в миллисекундах.
- `TIME_DIVISION_MS`: Время выполнения операции деления в миллисекундах.

### Тестирование и проверка работоспособности

1. **Добавление выражения для вычисления:**
    Пример для Linux/MacOS:
    ```sh
    curl --header 'Content-Type: application/json' --data '{"expression": "2 + 2 * 2"}' localhost:8080/api/v1/calculate
    ```
    Пример для Windows:
    ```powershell
    Invoke-RestMethod -Uri "http://localhost:8080/api/v1/calculate" -Method POST -Headers @{"Content-Type" = "application/json"} -Body '{"expression": "2 + 2 * 2"}'
    ```
    Вместо 2 + 2 * 2 подставляем нужное для вычисления арифметическое выражение    

    ***ВАЖНО*** В выражении могут быть только **неотрицательные** целочисленные или дробные (например 2.2) числа, символы операций (+-*/) и скобки

    В ответ будет получен id вычисления

1. **Получение всего списка выражений:**
    Пример для Linux/MacOS:
    ```sh
    curl http://localhost:8080/api/v1/expressions
    ```
    Пример для Windows:
    ```powershell
    Invoke-WebRequest -Uri "http://localhost:8080/api/v1/expressions" -Method Get
    ```
    В ответ будет получен список всех вычислений, их статусы и результаты

2. **Получение выражения по его идентификатору:**
    Пример для Linux/MacOS:
    ```sh
    curl http://localhost:8080/api/v1/expressions/1
    ```
    Пример для Windows:
    ```powershell
    Invoke-WebRequest -Uri "http://localhost:8080/api/v1/expressions/1" -Method Get
    ```
    Вместо 1 нужно написать id выражения    

    В ответ будет получен id, статус и результат указанного вычисления

### Для остановки требуется нажать Ctrl+C в терминале с запущенным проектом (где docker-compose)

## Заключение

После выполнения указанных шагов система будет запущена и готова к использованию. Вы можете добавлять арифметические выражения, получать их статус и результаты.

# Если у вас возникнут вопросы или проблемы, пожалуйста, прочитайте ещё раз внимательно данный файл или задайте мне вопрос в телеграм @aster_cmd
